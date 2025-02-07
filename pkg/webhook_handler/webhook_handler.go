package webhook_handler

import (
	"context"
	"fmt"
	"github.com/quickube/piper/pkg/clients"
	"github.com/quickube/piper/pkg/common"
	"github.com/quickube/piper/pkg/conf"
	"github.com/quickube/piper/pkg/git_provider"
	"github.com/quickube/piper/pkg/utils"
	"gopkg.in/yaml.v3"
	"log"
)

type WebhookHandlerImpl struct {
	cfg      *conf.GlobalConfig
	clients  *clients.Clients
	Triggers *[]Trigger
	Payload  *git_provider.WebhookPayload
}

func NewWebhookHandler(cfg *conf.GlobalConfig, clients *clients.Clients, payload *git_provider.WebhookPayload) (*WebhookHandlerImpl, error) {
	var err error

	return &WebhookHandlerImpl{
		cfg:      cfg,
		clients:  clients,
		Triggers: &[]Trigger{},
		Payload:  payload,
	}, err
}

func (wh *WebhookHandlerImpl) RegisterTriggers(ctx context.Context) error {
	if !IsFileExists(ctx, wh, "", ".workflows") {
		return fmt.Errorf(".workflows folder does not exist in %s/%s", wh.Payload.Repo, wh.Payload.Branch)
	}

	if !IsFileExists(ctx, wh, ".workflows", "triggers.yaml") {
		return fmt.Errorf(".workflows/triggers.yaml file does not exist in %s/%s", wh.Payload.Repo, wh.Payload.Branch)
	}

	triggers, err := wh.clients.GitProvider.GetFile(ctx, wh.Payload.Repo, wh.Payload.Branch, ".workflows/triggers.yaml")
	if err != nil {
		return fmt.Errorf("failed to get triggers content: %v", err)
	}

	log.Printf("triggers content is: \n %s \n", *triggers.Content) // DEBUG

	err = yaml.Unmarshal([]byte(*triggers.Content), wh.Triggers)
	if err != nil {
		return fmt.Errorf("failed to unmarshal triggers content: %v", err)
	}
	return nil
}

func (wh *WebhookHandlerImpl) PrepareBatchForMatchingTriggers(ctx context.Context) ([]*common.WorkflowsBatch, error) {
	triggered := false
	var workflowBatches []*common.WorkflowsBatch
	for _, trigger := range *wh.Triggers {
		if trigger.Branches == nil {
			return nil, fmt.Errorf("trigger from repo %s branch %s missing branch field", wh.Payload.Repo, wh.Payload.Branch)
		}
		if trigger.Events == nil {
			return nil, fmt.Errorf("trigger from repo %s branch %s missing event field", wh.Payload.Repo, wh.Payload.Branch)
		}

		eventToCheck := wh.Payload.Event
		if wh.Payload.Action != "" {
			eventToCheck += "." + wh.Payload.Action
		}
		if utils.IsElementMatch(wh.Payload.Branch, *trigger.Branches) && utils.IsElementMatch(eventToCheck, *trigger.Events) {
			log.Printf(
				"Triggering event %s for repo %s branch %s are triggered.",
				wh.Payload.Event,
				wh.Payload.Repo,
				wh.Payload.Branch,
			)
			triggered = true
			onStartFiles, err := wh.clients.GitProvider.GetFiles(
				ctx,
				wh.Payload.Repo,
				wh.Payload.Branch,
				utils.AddPrefixToList(*trigger.OnStart, ".workflows/"),
			)
			if len(onStartFiles) == 0 {
				return nil, fmt.Errorf("one or more of onStart: %s files found in repo: %s branch %s", *trigger.OnStart, wh.Payload.Repo, wh.Payload.Branch)
			}
			if err != nil {
				return nil, err
			}

			onExitFiles := make([]*git_provider.CommitFile, 0)
			if trigger.OnExit != nil {
				onExitFiles, err = wh.clients.GitProvider.GetFiles(
					ctx,
					wh.Payload.Repo,
					wh.Payload.Branch,
					utils.AddPrefixToList(*trigger.OnExit, ".workflows/"),
				)
				if len(onExitFiles) == 0 {
					log.Printf("one or more of onExist: %s files not found in repo: %s branch %s", *trigger.OnExit, wh.Payload.Repo, wh.Payload.Branch)
				}
				if err != nil {
					return nil, err
				}
			}

			templatesFiles := make([]*git_provider.CommitFile, 0)
			if trigger.Templates != nil {
				templatesFiles, err = wh.clients.GitProvider.GetFiles(
					ctx,
					wh.Payload.Repo,
					wh.Payload.Branch,
					utils.AddPrefixToList(*trigger.Templates, ".workflows/"),
				)
				if len(templatesFiles) == 0 {
					log.Printf("one or more of templates: %s files not found in repo: %s branch %s", *trigger.Templates, wh.Payload.Repo, wh.Payload.Branch)
				}
				if err != nil {
					return nil, err
				}
			}

			parameters := &git_provider.CommitFile{
				Path:    nil,
				Content: nil,
			}
			if IsFileExists(ctx, wh, ".workflows", "parameters.yaml") {
				parameters, err = wh.clients.GitProvider.GetFile(
					ctx,
					wh.Payload.Repo,
					wh.Payload.Branch,
					".workflows/parameters.yaml",
				)
				if err != nil {
					return nil, err
				}
			} else {
				log.Printf("parameters.yaml not found in repo: %s branch %s", wh.Payload.Repo, wh.Payload.Branch)
			}

			workflowBatches = append(workflowBatches, &common.WorkflowsBatch{
				OnStart:    onStartFiles,
				OnExit:     onExitFiles,
				Templates:  templatesFiles,
				Parameters: parameters,
				Config:     &trigger.Config,
				Payload:    wh.Payload,
			})
		}
	}
	if !triggered {
		return nil, fmt.Errorf("no matching trigger found for event: %s action: %s in branch :%s", wh.Payload.Event, wh.Payload.Action, wh.Payload.Branch)
	}
	return workflowBatches, nil
}

func IsFileExists(ctx context.Context, wh *WebhookHandlerImpl, path string, file string) bool {
	files, err := wh.clients.GitProvider.ListFiles(ctx, wh.Payload.Repo, wh.Payload.Branch, path)
	if err != nil {
		log.Printf("Error listing files in repo: %s branch: %s. %v", wh.Payload.Repo, wh.Payload.Branch, err)
		return false
	}
	if len(files) == 0 {
		log.Printf("Empty list of files in repo: %s branch: %s", wh.Payload.Repo, wh.Payload.Branch)
		return false
	}

	if utils.IsElementExists(files, file) {
		return true
	}

	return false
}

func HandleWebhook(ctx context.Context, wh *WebhookHandlerImpl) ([]*common.WorkflowsBatch, error) {
	err := wh.RegisterTriggers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to register triggers, error: %v", err)
	} else {
		log.Printf("successfully registered triggers for repo: %s branch: %s", wh.Payload.Repo, wh.Payload.Branch)
	}

	workflowsBatches, err := wh.PrepareBatchForMatchingTriggers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare matching triggers, error: %v", err)
	}

	if len(workflowsBatches) == 0 {
		log.Printf("no workflows to execute")
		return nil, fmt.Errorf("no workflows to execute for repo: %s branch: %s",
			wh.Payload.Repo,
			wh.Payload.Branch,
		)
	}
	return workflowsBatches, nil
}
