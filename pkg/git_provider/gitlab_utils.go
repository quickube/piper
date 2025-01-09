package git_provider

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/quickube/piper/pkg/conf"
	"github.com/quickube/piper/pkg/utils"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/net/context"
)

func ValidateGitlabPermissions(ctx context.Context, client *gitlab.Client, cfg *conf.GlobalConfig) error {

	repoAdminScopes := []string{"api"}
	repoGranularScopes := []string{"write_repository", "read_api"}

	token, _, err := client.PersonalAccessTokens.GetSinglePersonalAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get scopes: %v", err)
	}
	scopes := token.Scopes

	if len(scopes) == 0 {
		return fmt.Errorf("permissions error: no scopes found for the gitlab client")
	}

	if utils.ListContains(repoAdminScopes, scopes) {
		return nil
	}
	if utils.ListContains(repoGranularScopes, scopes) {
		return nil
	}

	return fmt.Errorf("permissions error: %v is not a valid scope for the project level permissions", scopes)
}

func IsGroupWebhookEnabled(ctx context.Context, c *GitlabClientImpl) (*gitlab.GroupHook, bool) {
	emptyHook := gitlab.GroupHook{}
	hooks, resp, err := c.client.Groups.ListGroupHooks(c.cfg.GitProviderConfig.OrgName, nil, gitlab.WithContext(ctx))

	if err != nil {
		return &emptyHook, false
	}
	if resp.StatusCode != 200 {
		return &emptyHook, false
	}
	if len(hooks) != 0 {
		for _, hook := range hooks {
			if hook.URL == c.cfg.GitProviderConfig.WebhookURL {
				return hook, true
			}
		}
	}
	return &emptyHook, false
}

func IsProjectWebhookEnabled(ctx context.Context, c *GitlabClientImpl, projectId int) (*gitlab.ProjectHook, bool) {
	emptyHook := gitlab.ProjectHook{}

	hooks, resp, err := c.client.Projects.ListProjectHooks(projectId, nil, gitlab.WithContext(ctx))
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

func ExtractLabelsId(labels []*gitlab.EventLabel) []string {
	var returnLabelsList []string
	for _, label := range labels {
		returnLabelsList = append(returnLabelsList, fmt.Sprint(label.ID))
	}
	return returnLabelsList
}

func GetProjectId(ctx context.Context, c *GitlabClientImpl, repo *string) (*int, error) {
	projectFullName := fmt.Sprintf("%s/%s", c.cfg.GitProviderConfig.OrgName, *repo)
	IProject, _, err := c.client.Projects.GetProject(projectFullName, nil, gitlab.WithContext(ctx))
	if err != nil {
		log.Printf("Failed to get project (%s): %v", *repo, err)
		return nil, err
	}
	return &IProject.ID, nil
}

func ValidatePayload(r *http.Request, secret []byte) ([]byte, error) {
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

	isEqual := hmac.Equal([]byte(gitlabSignature), []byte(expectedMAC))
	if !isEqual {
		return nil, fmt.Errorf("secret not correct")
	}
	return payload, nil
}

func FixRepoNames(c *GitlabClientImpl) error {
	var formattedRepos []string
	for _, repo := range strings.Split(c.cfg.GitProviderConfig.RepoList, ",") {
		userRepo := fmt.Sprintf("%s/%s", c.cfg.GitProviderConfig.OrgName, repo)
		formattedRepos = append(formattedRepos, userRepo)
	}
	c.cfg.GitProviderConfig.RepoList = strings.Join(formattedRepos, ",")
	return nil
}

func DecodeBase64ToStringPtr(encoded string) (*string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	result := string(decoded)
	return &result, nil
}
