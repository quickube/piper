package git_provider

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/quickube/piper/pkg/conf"
	"github.com/quickube/piper/pkg/utils"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/net/context"
)

func ValidateGitlabPermissions(ctx context.Context, client *gitlab.Client, cfg *conf.GlobalConfig) error {

	orgScopes := []string{"admin:org_hook"}
	repoAdminScopes := []string{"admin:repo_hook"}
	repoGranularScopes := []string{"write:repo_hook", "read:repo_hook"}

	scopes, err := GetGitlabScopes(ctx, client)

	if err != nil {
		return fmt.Errorf("failed to get scopes: %v", err)
	}
	if len(scopes) == 0 {
		return fmt.Errorf("permissions error: no scopes found for the gitlab client")
	}

	if cfg.GitProviderConfig.OrgLevelWebhook {
		if utils.ListContains(orgScopes, scopes) {
			return nil
		}
		return fmt.Errorf("permissions error: %v is not a valid scope for the org level permissions", scopes)
	}

	if utils.ListContains(repoAdminScopes, scopes) {
		return nil
	}
	if utils.ListContains(repoGranularScopes, scopes) {
		return nil
	}

	return fmt.Errorf("permissions error: %v is not a valid scope for the repo level permissions", scopes)
}

func GetGitlabScopes(ctx context.Context, client *gitlab.Client) ([]string, error) {
	req, err := retryablehttp.NewRequest("GET", "https://gitlab.com/api/v4/user", nil)
    if err != nil {
        log.Fatalf("Failed to create request: %v", err)
    }

    resp, err := client.Do(req, nil)
    if err != nil {
        log.Fatalf("Failed to perform request: %v", err)
		return nil, err

    }
    defer resp.Body.Close()

    // Check for successful response
    if resp.StatusCode != http.StatusOK {
        log.Fatalf("Failed to get user: %v", resp.Status)
		return nil, err
    }

	scopes := resp.Header.Get("X-OAuth-Scopes")
	fmt.Println("Github Token Scopes are:", scopes)

	scopes = strings.ReplaceAll(scopes, " ", "")
	return strings.Split(scopes, ","), nil
}

func isGroupWebhookEnabled(c *GitlabClientImpl) (*gitlab.GroupHook, bool) {
	emptyHook := gitlab.GroupHook{}
	hooks, resp, err := c.client.Groups.ListGroupHooks(c.cfg.GitProviderConfig.OrgName, nil)
	if err != nil {
		return &emptyHook, false
	}
	if resp.StatusCode != 200 {
		return &emptyHook, false
	}
	if len(hooks) == 0 {
		return &emptyHook, false
	}
	for _, hook := range hooks {
		if hook.AlertStatus == "triggered" && hook.URL == c.cfg.GitProviderConfig.WebhookURL {
			return hook, true
		}
	}
	return &emptyHook, false
}

func isProjectWebhookEnabled(c *GitlabClientImpl, repo string) (*gitlab.ProjectHook, bool) {
	emptyHook := gitlab.ProjectHook{}
	hooks, resp, err := c.client.Projects.ListProjectHooks(repo, nil)
	if err != nil {
		return &emptyHook, false
	}
	if resp.StatusCode != 200 {
		return &emptyHook, false
	}
	if len(hooks) == 0 {
		return &emptyHook, false
	}

	for _, hook := range hooks {
		if hook.URL == c.cfg.GitProviderConfig.WebhookURL {
			return hook, true
		}
	}

	return &emptyHook, false
}

func extractLabelsId(labels []*gitlab.EventLabel) []string {
	var returnLabelsList []string
	for _, label := range labels {
		returnLabelsList = append(returnLabelsList, fmt.Sprint(label.ID))
	}
	return returnLabelsList
}

func validatePayload(r *http.Request, secret []byte) ([]byte, error){
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %v", err)
	}

	// Get GitLab signature from headers
	gitlabSignature := r.Header.Get("X-Gitlab-Token")
	if gitlabSignature == "" {
		return nil, fmt.Errorf("no GitLab signature found in headers")
	}

	h := hmac.New(sha256.New, secret)
	_, err = h.Write(payload)
	if err != nil {
		return nil, fmt.Errorf("error computing HMAC: %v", err)
	}
	expectedMAC := hex.EncodeToString(h.Sum(nil))

	isEquall := hmac.Equal([]byte(gitlabSignature), []byte(expectedMAC))
	if !isEquall {
		return nil, fmt.Errorf("secret not correct")
	}
	return payload, nil
}
