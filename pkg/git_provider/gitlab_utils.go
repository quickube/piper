package git_provider

import (
	"fmt"
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