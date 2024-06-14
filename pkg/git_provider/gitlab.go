package git_provider

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/quickube/piper/pkg/conf"
	"github.com/quickube/piper/pkg/utils"

	"github.com/xanzy/go-gitlab"
)

type GitlabClientImpl struct {
	client *gitlab.Client
	cfg    *conf.GlobalConfig
}

func NewGitlabClient(cfg *conf.GlobalConfig) (Client, error) {
	ctx := context.Background()

	client, err := gitlab.NewClient(cfg.GitProviderConfig.Token)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	err = ValidateGitlabPermissions(ctx, client, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to validate permissions: %v", err)
	}

	group, resp, err := client.Groups.GetGroup(cfg.OrgName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get organization data %s", resp.Status)
	}

	cfg.OrgID = int64(group.ID)

	log.Printf("Org ID is: %d\n", cfg.OrgID)

	return &GitlabClientImpl{
		client: client,
		cfg:    cfg,
	}, err
}

func (c *GitlabClientImpl) ListFiles(ctx *context.Context, repo string, branch string, path string) ([]string, error) {
	var files []string
	opt := &gitlab.ListTreeOptions{
		Ref: &branch,
		Path: &path,}
	
	dirFiles, resp, err := c.client.Repositories.ListTree(repo, opt)


	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("gitlab provider returned %d: failed to get contents of %s/%s%s", resp.StatusCode, repo, branch, path)
	}
	if files == nil {
		return nil, nil
	}
	for _, file := range dirFiles {
		files = append(files, file.Name)
	}
	return files, nil
}

func (c *GitlabClientImpl) GetFile(ctx *context.Context, repo string, branch string, path string) (*CommitFile, error) {
	var commitFile CommitFile

	opt := &gitlab.GetFileOptions{Ref: &branch,}
	fileContent, resp, err := c.client.RepositoryFiles.GetFile(repo, path,opt)
	if err != nil {
		return &commitFile, err
	}
	if resp.StatusCode == 404 {
		log.Printf("File %s not found in repo %s branch %s", path, repo, branch)
		return nil, nil
	}
	if resp.StatusCode != 200 {
		return &commitFile, err
	}
	if fileContent == nil {
		return &commitFile, nil
	}
	filePath := fileContent.FilePath
	commitFile.Path = &filePath
	fileContentString := fileContent.Content
	commitFile.Content = &fileContentString

	return &commitFile, nil
}

func (c *GitlabClientImpl) GetFiles(ctx *context.Context, repo string, branch string, paths []string) ([]*CommitFile, error) {
	var commitFiles []*CommitFile
	for _, path := range paths {
		file, err := c.GetFile(ctx, repo, branch, path)
		if err != nil {
			return nil, err
		}
		if file == nil {
			log.Printf("file %s not found in repo %s branch %s", path, repo, branch)
			continue
		}
		commitFiles = append(commitFiles, file)
	}
	return commitFiles, nil
}

func (c *GitlabClientImpl) SetWebhook(ctx *context.Context, repo *string) (*HookWithStatus, error) {
	if c.cfg.OrgLevelWebhook && repo != nil {
		return nil, fmt.Errorf("trying to set repo scope. repo: %s", *repo)
	}
	var gitlabHook gitlab.Hook

	if repo == nil {
		respHook, ok := isGroupWebhookEnabled(c)

		if !ok {
			groupHookOptions := gitlab.AddGroupHookOptions{
				URL: gitlab.Ptr(c.cfg.GitProviderConfig.WebhookURL),
				Token: gitlab.Ptr(c.cfg.GitProviderConfig.WebhookSecret),
				MergeRequestsEvents: gitlab.Ptr(true),
				PushEvents: gitlab.Ptr(true),
				ReleasesEvents: gitlab.Ptr(true),
				TagPushEvents: gitlab.Ptr(true),				
			}

			gitlabHook, resp, err := c.client.Groups.AddGroupHook(c.cfg.GitProviderConfig.OrgName, &groupHookOptions)
			if err != nil {
				return nil, err
			}
			if resp.StatusCode != 201 {
				return nil, fmt.Errorf("failed to create group level webhhok, API returned %d", resp.StatusCode)
			}
			log.Printf("added webhook for %s name: %s\n", c.cfg.GitProviderConfig.OrgName, gitlabHook.URL)
		} else {
			editedGroupHookOpt := gitlab.EditGroupHookOptions{
				URL: gitlab.Ptr(c.cfg.GitProviderConfig.WebhookURL),
				Token: gitlab.Ptr(c.cfg.GitProviderConfig.WebhookSecret),
				MergeRequestsEvents: gitlab.Ptr(true),
				PushEvents: gitlab.Ptr(true),
				ReleasesEvents: gitlab.Ptr(true),
				TagPushEvents: gitlab.Ptr(true),				
			}
			gitlabHook, resp, err := c.client.Groups.EditGroupHook(
				c.cfg.GitProviderConfig.OrgName,
				respHook.ID,
				&editedGroupHookOpt,
			)
			if err != nil {
				return nil, err
			}
			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf(
					"failed to update group level webhook for %s, API returned %d",
					c.cfg.GitProviderConfig.OrgName,
					resp.StatusCode,
				)
			}
			log.Printf("edited webhook for %s: %s\n", c.cfg.GitProviderConfig.OrgName, gitlabHook.URL)
		}
	} else {
		respHook, ok := isProjectWebhookEnabled(*ctx, c, *repo)
		if !ok {
			addProjectHookOpts := gitlab.AddProjectHookOptions{
				URL: gitlab.Ptr(c.cfg.GitProviderConfig.WebhookURL),
				Token: gitlab.Ptr(c.cfg.GitProviderConfig.WebhookSecret),
				MergeRequestsEvents: gitlab.Ptr(true),
				PushEvents: gitlab.Ptr(true),
				ReleasesEvents: gitlab.Ptr(true),
				TagPushEvents: gitlab.Ptr(true),
			}

			gitlabHook, resp, err := c.client.Projects.AddProjectHook(repo, &addProjectHookOpts)
			if err != nil {
				return nil, err
			}
			if resp.StatusCode != 201 {
				return nil, fmt.Errorf("failed to create repo level webhhok for %s, API returned %d", *repo, resp.StatusCode)
			}
			log.Printf("created webhook for %s: %s\n", *repo, gitlabHook.URL)
		} else {
			
			editProjectHookOpts := gitlab.EditProjectHookOptions{
				URL: gitlab.Ptr(c.cfg.GitProviderConfig.WebhookURL),
				Token: gitlab.Ptr(c.cfg.GitProviderConfig.WebhookSecret),
				MergeRequestsEvents: gitlab.Ptr(true),
				PushEvents: gitlab.Ptr(true),
				ReleasesEvents: gitlab.Ptr(true),
				TagPushEvents: gitlab.Ptr(true),
			}
			gitlabHook, resp, err := c.client.Projects.EditProjectHook(repo, respHook.ID, &editProjectHookOpts)
			if err != nil {
				return nil, err
			}
			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("failed to update repo level webhhok for %s, API returned %d", *repo, resp.StatusCode)
			}
			log.Printf("edited webhook for %s: %s\n", *repo, gitlabHook.URL)
		}

	}

	hookID := int64(gitlabHook.ID)
	return &HookWithStatus{HookID: hookID, HealthStatus: true, RepoName: repo}, nil
}

func (c *GitlabClientImpl) UnsetWebhook(ctx *context.Context, hook *HookWithStatus) error {

	if hook.RepoName == nil {
		resp, err := c.client.Groups.DeleteGroupHook( c.cfg.GitProviderConfig.OrgName, int(hook.HookID))
		if err != nil {
			return err
		}

		if resp.StatusCode != 204 {
			return fmt.Errorf("failed to delete group level webhhok, API call returned %d", resp.StatusCode)
		}
		log.Printf("removed group webhook, hookID :%d\n", hook.HookID)
	} else {
		resp, err := c.client.Projects.DeleteProjectHook(*hook.RepoName, int(hook.HookID))

		if err != nil {
			return fmt.Errorf("failed to delete project level webhhok for %s, API call returned %d. %s", *hook.RepoName, resp.StatusCode, err)
		}

		if resp.StatusCode != 204 {
			return fmt.Errorf("failed to delete project level webhhok for %s, API call returned %d", *hook.RepoName, resp.StatusCode)
		}
		log.Printf("removed project webhook, project:%s hookID :%d\n", *hook.RepoName, hook.HookID) // INFO
	}

	return nil
}

func (c *GitlabClientImpl) HandlePayload(ctx *context.Context, request *http.Request, secret []byte) (*WebhookPayload, error) {
// 	var webhookPayload *WebhookPayload

// 	payload, err := github.ValidatePayload(request, secret)
// 	if err != nil {
// 		return nil, err
// 	}

// 	event, err := github.ParseWebHook(github.WebHookType(request), payload)
// 	if err != nil {
// 		return nil, err
// 	}

// 	switch e := event.(type) {
// 	case *github.PingEvent:
// 		webhookPayload = &WebhookPayload{
// 			Event:   "ping",
// 			Repo:    e.GetRepo().GetFullName(),
// 			HookID:  e.GetHookID(),
// 			OwnerID: e.GetSender().GetID(),
// 		}
// 	case *github.PushEvent:
// 		webhookPayload = &WebhookPayload{
// 			Event:     "push",
// 			Action:    e.GetAction(),
// 			Repo:      e.GetRepo().GetName(),
// 			Branch:    strings.TrimPrefix(e.GetRef(), "refs/heads/"),
// 			Commit:    e.GetHeadCommit().GetID(),
// 			User:      e.GetSender().GetLogin(),
// 			UserEmail: e.GetHeadCommit().GetAuthor().GetEmail(),
// 			OwnerID:   e.GetSender().GetID(),
// 		}
// 	case *github.PullRequestEvent:
// 		webhookPayload = &WebhookPayload{
// 			Event:            "pull_request",
// 			Action:           e.GetAction(),
// 			Repo:             e.GetRepo().GetName(),
// 			Branch:           e.GetPullRequest().GetHead().GetRef(),
// 			Commit:           e.GetPullRequest().GetHead().GetSHA(),
// 			User:             e.GetPullRequest().GetUser().GetLogin(),
// 			UserEmail:        e.GetSender().GetEmail(), // e.GetPullRequest().GetUser().GetEmail() Not working. GitHub missing email for PR events in payload.
// 			PullRequestTitle: e.GetPullRequest().GetTitle(),
// 			PullRequestURL:   e.GetPullRequest().GetHTMLURL(),
// 			DestBranch:       e.GetPullRequest().GetBase().GetRef(),
// 			Labels:           c.extractLabelNames(e.GetPullRequest().Labels),
// 			OwnerID:          e.GetSender().GetID(),
// 		}
// 	case *github.CreateEvent:
// 		webhookPayload = &WebhookPayload{
// 			Event:     "create",
// 			Action:    e.GetRefType(), // Possible values are: "repository", "branch", "tag".
// 			Repo:      e.GetRepo().GetName(),
// 			Branch:    e.GetRef(),
// 			Commit:    e.GetRef(),
// 			User:      e.GetSender().GetLogin(),
// 			UserEmail: e.GetSender().GetEmail(),
// 			OwnerID:   e.GetSender().GetID(),
// 		}
// 	case *github.ReleaseEvent:
// 		commitSHA, _err := c.refToSHA(ctx, e.GetRelease().GetName(), e.GetRepo().GetName())
// 		if _err != nil {
// 			return webhookPayload, _err
// 		}
// 		webhookPayload = &WebhookPayload{
// 			Event:     "release",
// 			Action:    e.GetAction(), // "created", "edited", "deleted", or "prereleased".
// 			Repo:      e.GetRepo().GetName(),
// 			Branch:    e.GetRelease().GetTagName(),
// 			Commit:    *commitSHA,
// 			User:      e.GetSender().GetLogin(),
// 			UserEmail: e.GetSender().GetEmail(),
// 			OwnerID:   e.GetSender().GetID(),
// 		}
// 	}

// 	if c.cfg.EnforceOrgBelonging && (webhookPayload.OwnerID == 0 || webhookPayload.OwnerID != c.cfg.OrgID) {
// 		return nil, fmt.Errorf("webhook send from non organizational member")
// 	}
// 	return webhookPayload, nil
panic("implement me")

}

func (c *GitlabClientImpl) SetStatus(ctx *context.Context, repo *string, commit *string, linkURL *string, status *string, message *string) error {
	if !utils.ValidateHTTPFormat(*linkURL) {
		return fmt.Errorf("invalid linkURL")
	}

	repoStatus := &gitlab.SetCommitStatusOptions{
		State:       gitlab.BuildStateValue(*status), // pending, success, error, or failure.
		Ref: commit,
		TargetURL:   linkURL,
		Description: gitlab.Ptr(fmt.Sprintf("Workflow %s %s", *status, *message)),
		Context:     gitlab.Ptr("Piper/ArgoWorkflows"),
	}
	
	_, resp, err := c.client.Commits.SetCommitStatus(*repo, *commit, repoStatus)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to set status on repo:%s, commit:%s, API call returned %d", *repo, *commit, resp.StatusCode)
	}

	log.Printf("successfully set status on repo:%s commit: %s to status: %s\n", *repo, *commit, *status)
	return nil
}

func (c *GitlabClientImpl) PingHook(ctx *context.Context, hook *HookWithStatus) error {
	if c.cfg.OrgLevelWebhook && hook.RepoName != nil {
		return fmt.Errorf("trying to ping repo scope webhook while configured for org level webhook. repo: %s", *hook.RepoName)
	}
	if hook.RepoName == nil {
		_,resp, err := c.client.Groups.GetGroupHook(c.cfg.OrgName,int(hook.HookID), nil)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("unable to find organization webhook for hookID: %d", hook.HookID)
		}
	} else {
		_,resp, err := c.client.Projects.GetProjectHook(hook.RepoName, int(hook.HookID), nil)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("unable to find repo webhook for repo:%s hookID: %d", *hook.RepoName, hook.HookID)
		}
	}

	return nil
}
