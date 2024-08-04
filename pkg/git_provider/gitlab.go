package git_provider

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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

	client, err := gitlab.NewClient(cfg.GitProviderConfig.Token,  gitlab.WithBaseURL(cfg.GitProviderConfig.Url))
	if err != nil {
		return nil, err
	}

	err = ValidateGitlabPermissions(ctx, client, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to validate permissions: %v", err)
	}

	group, resp, err := client.Groups.GetGroup(cfg.GitProviderConfig.OrgName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get organization data %s", resp.Status)
	}

	c := &GitlabClientImpl{
		client: client,
		cfg:    cfg,
	}

	c.cfg.GitProviderConfig.OrgID = int64(group.ID)
	log.Printf("Org ID is: %d\n", cfg.OrgID)
	
	return c, err
}

func (c *GitlabClientImpl) ListFiles(ctx *context.Context, repo string, branch string, path string) ([]string, error) {
	var files []string
	opt := &gitlab.ListTreeOptions{
		Ref: &branch,
		Path: &path,}
	
	projectName := fmt.Sprintf("%s/%s", c.cfg.GitProviderConfig.OrgName ,repo)
	dirFiles, resp, err := c.client.Repositories.ListTree(projectName, opt)


	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("gitlab provider returned %d: failed to get contents of %s/%s%s", resp.StatusCode, repo, branch, path)
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
	if resp.StatusCode != 200 {
		return &commitFile, err
	}
	commitFile.Path = &fileContent.FilePath
	commitFile.Content = &fileContent.Content

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
		return nil, fmt.Errorf("trying to set project scope. project: %s", *repo)
	}
	var gitlabHook gitlab.Hook

	if repo == nil {
		respHook, ok := IsGroupWebhookEnabled(c)

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
			gitlabHook, resp, err := c.client.Groups.EditGroupHook(c.cfg.GitProviderConfig.OrgName,respHook.ID,&editedGroupHookOpt,)
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
		respHook, ok := IsProjectWebhookEnabled(c, *repo)
		if !ok {
			addProjectHookOpts := gitlab.AddProjectHookOptions{
				URL: gitlab.Ptr(c.cfg.GitProviderConfig.WebhookURL),
				Token: gitlab.Ptr(c.cfg.GitProviderConfig.WebhookSecret),
				MergeRequestsEvents: gitlab.Ptr(true),
				PushEvents: gitlab.Ptr(true),
				ReleasesEvents: gitlab.Ptr(true),
				TagPushEvents: gitlab.Ptr(true),
			}
			
			gitlabHook, resp, err := c.client.Projects.AddProjectHook(*repo, &addProjectHookOpts)
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
			gitlabHook, resp, err := c.client.Projects.EditProjectHook(*repo, respHook.ID, &editProjectHookOpts)
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
	var webhookPayload *WebhookPayload
	payload, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %v", err)
	}
	event, err := gitlab.ParseWebhook(gitlab.WebhookEventType(request), payload)
	if err != nil {
		return nil, err
	}


	switch e := event.(type) {
	case gitlab.PushEvent:
		webhookPayload = &WebhookPayload{
			Event:     "push",
			Action:    e.EventName,
			Repo:      e.Project.Name,
			Branch:    strings.TrimPrefix(e.Ref, "refs/heads/"),
			Commit:    e.CheckoutSHA,
			User:      e.UserName,
			UserEmail: e.UserEmail,
			OwnerID:   int64(e.UserID),
		}
	case gitlab.MergeEvent:
		webhookPayload = &WebhookPayload{
			Event:            "merge_request",
			Action:           e.ObjectAttributes.Action,
			Repo:             e.Repository.Name,
			Branch:           e.ObjectAttributes.SourceBranch,
			Commit:           e.ObjectAttributes.LastCommit.ID,
			User:             e.User.Name,
			UserEmail:        e.User.Email,
			PullRequestTitle: e.ObjectAttributes.Title,
			PullRequestURL:   e.ObjectAttributes.URL,
			DestBranch:       e.ObjectAttributes.TargetBranch,
			Labels:           ExtractLabelsId(e.Labels),
			OwnerID:          int64(e.User.ID),
		}
	case gitlab.ReleaseEvent:
		webhookPayload = &WebhookPayload{
			Event:     "release",
			Action:    e.Action, // "create" | "update" | "delete"
			Repo:      e.Project.Name,
			Branch:    e.Tag,
			Commit:    e.Commit.ID,
			User:      e.Commit.Author.Name,
			UserEmail: e.Commit.Author.Email,
		}
	}

	if (webhookPayload.OwnerID == 0) || (webhookPayload.OwnerID != c.cfg.OrgID) {
		return nil, fmt.Errorf("webhook send from non organizational member")
	}
	return webhookPayload, nil
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
		return fmt.Errorf("trying o ping repo scope webhook while configured for org level webhook. repo: %s", *hook.RepoName)
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
		_ , resp, err := c.client.Projects.GetProjectHook(*hook.RepoName, int(hook.HookID), nil)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("unable to find repo webhook for repo:%s hookID: %d", *hook.RepoName, hook.HookID)
		}
	}

	return nil
}
