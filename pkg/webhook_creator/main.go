package webhook_creator

import (
	"fmt"
	"github.com/emicklei/go-restful/v3/log"
	"github.com/quickube/piper/pkg/clients"
	"github.com/quickube/piper/pkg/conf"
	"github.com/quickube/piper/pkg/git_provider"
	"github.com/quickube/piper/pkg/utils"
	"golang.org/x/net/context"
	"strconv"
	"strings"
	"sync"
	"time"
)

type WebhookCreatorImpl struct {
	clients *clients.Clients
	cfg     *conf.GlobalConfig
	hooks   map[int64]*git_provider.HookWithStatus
	mu      sync.Mutex
}

func NewWebhookCreator(cfg *conf.GlobalConfig, clients *clients.Clients) *WebhookCreatorImpl {
	wr := &WebhookCreatorImpl{
		clients: clients,
		cfg:     cfg,
		hooks:   make(map[int64]*git_provider.HookWithStatus, 0),
	}

	return wr
}

func (wc *WebhookCreatorImpl) GetHooks() *map[int64]*git_provider.HookWithStatus {
	return &wc.hooks
}

func (wc *WebhookCreatorImpl) Start(ctx context.Context) {

	err := wc.initWebhooks(ctx)
	if err != nil {
		log.Print(err)
		panic("failed in initializing webhooks")
	}
}

func (wc *WebhookCreatorImpl) setWebhook(hookID int64, healthStatus bool, repoName string) {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	wc.hooks[hookID] = &git_provider.HookWithStatus{HookID: hookID, HealthStatus: healthStatus, RepoName: &repoName}
}

func (wc *WebhookCreatorImpl) getWebhook(hookID int64) *git_provider.HookWithStatus {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	hook, ok := wc.hooks[hookID]
	if !ok {
		return nil
	}
	return hook
}

func (wc *WebhookCreatorImpl) deleteWebhook(hookID int64) {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	delete(wc.hooks, hookID)
}

func (wc *WebhookCreatorImpl) SetWebhookHealth(hookID int64, status bool) error {

	hook, ok := wc.hooks[hookID]
	if !ok {
		return fmt.Errorf("unable to find hookID: %d in internal hooks map %v", hookID, wc.hooks)
	}
	wc.setWebhook(hookID, status, *hook.RepoName)
	log.Printf("set health status to %s for hook id: %d", strconv.FormatBool(status), hookID)
	return nil
}

func (wc *WebhookCreatorImpl) setAllHooksHealth(status bool) {
	for hookID, hook := range wc.hooks {
		wc.setWebhook(hookID, status, *hook.RepoName)
	}
	log.Printf("set all hooks health status for to %s", strconv.FormatBool(status))
}

func (wc *WebhookCreatorImpl) initWebhooks(ctx context.Context) error {

	if wc.cfg.GitProviderConfig.OrgLevelWebhook && len(wc.cfg.GitProviderConfig.RepoList) != 0 {
		return fmt.Errorf("org level webhook wanted but provided repositories list")
	} else if !wc.cfg.GitProviderConfig.OrgLevelWebhook && len(wc.cfg.GitProviderConfig.RepoList) == 0 {
		return fmt.Errorf("either org level webhook or repos list must be provided")
	}
	for _, repo := range strings.Split(wc.cfg.GitProviderConfig.RepoList, ",") {
		if wc.cfg.GitProviderConfig.Provider == "bitbucket" {
			repo = utils.SanitizeString(repo)
		}
		hook, err := wc.clients.GitProvider.SetWebhook(ctx, &repo)
		if err != nil {
			return err
		}
		wc.setWebhook(hook.HookID, hook.HealthStatus, *hook.RepoName)
	}

	return nil
}

func (wc *WebhookCreatorImpl) Stop(ctx context.Context) {
	if wc.cfg.GitProviderConfig.WebhookAutoCleanup {
		err := wc.deleteWebhooks(ctx)
		if err != nil {
			log.Printf("Failed to delete webhooks, error: %v", err)
		}
	}
}

func (wc *WebhookCreatorImpl) deleteWebhooks(ctx context.Context) error {
	for hookID, hook := range wc.hooks {
		err := wc.clients.GitProvider.UnsetWebhook(ctx, hook)
		if err != nil {
			return err
		}
		wc.deleteWebhook(hookID)
	}

	return nil
}

func (wc *WebhookCreatorImpl) checkHooksHealth(timeoutSeconds time.Duration) bool {
	startTime := time.Now()

	for {
		allHealthy := true
		for _, hook := range wc.hooks {
			if !hook.HealthStatus {
				allHealthy = false
				break
			}
		}

		if allHealthy {
			return true
		}

		if time.Since(startTime) >= timeoutSeconds {
			break
		}

		time.Sleep(1 * time.Second) // Adjust the sleep duration as per your requirement
	}

	return false
}

func (wc *WebhookCreatorImpl) recoverHook(ctx context.Context, hookID int64) error {

	log.Printf("started recover of hook %d", hookID)
	hook := wc.getWebhook(hookID)
	if hook == nil {
		return fmt.Errorf("failed to recover hook, hookID %d not found", hookID)
	}
	newHook, err := wc.clients.GitProvider.SetWebhook(ctx, hook.RepoName)
	if err != nil {
		return err
	}
	wc.deleteWebhook(hookID)
	wc.setWebhook(newHook.HookID, newHook.HealthStatus, *newHook.RepoName)
	log.Printf("successful recover of hook %d", hookID)
	return nil

}

func (wc *WebhookCreatorImpl) pingHooks(ctx context.Context) error {
	for _, hook := range wc.hooks {
		err := wc.clients.GitProvider.PingHook(ctx, hook)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wc *WebhookCreatorImpl) RunDiagnosis(ctx context.Context) error {
	log.Printf("Starting webhook diagnostics")
	wc.setAllHooksHealth(false)
	err := wc.pingHooks(ctx)
	if err != nil {
		return err
	}
	if !wc.checkHooksHealth(5 * time.Second) {
		for hookID, hook := range wc.hooks {
			if !hook.HealthStatus {
				return fmt.Errorf("hook %d is not healthy", hookID)
			}
		}
	}

	log.Print("Successful webhook diagnosis")
	return nil
}
