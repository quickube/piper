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
	var options []gitlab.ClientOptionFunc
	ctx := context.Background()

	if cfg.GitProviderConfig.Url != "" {
		options = append(options, gitlab.WithBaseURL(cfg.GitProviderConfig.Url))
	}
	client, err := gitlab.NewClient(cfg.GitProviderConfig.Token, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate user: %v", err)
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

	cfg.GitProviderConfig.OrgID = int64(group.ID)

	log.Printf("Group ID is: %d\n", cfg.OrgID)

	return &GitlabClientImpl{
		client: client,
		cfg:    cfg,
	}, err
}

func (c *GitlabClientImpl) ListFiles(ctx *context.Context, repo string, branch string, path string) ([]string, error) {
	log.Printf("Listing files for repo: %s", repo)
	var files []string
	opt := &gitlab.ListTreeOptions{
		Ref:  &branch,
		Path: &path}

	projectId, err := GetProjectId(ctx, c, &repo)
	if err != nil {
		return nil, err
	}
	dirFiles, resp, err := c.client.Repositories.ListTree(*projectId, opt, gitlab.WithContext(*ctx))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gitlab provider returned %d: failed to get contents of %s/%s%s", resp.StatusCode, repo, branch, path)
	}
	for _, file := range dirFiles {
		files = append(files, file.Name)
	}
	return files, nil
}

func (c *GitlabClientImpl) GetFile(ctx *context.Context, repo string, branch string, path string) (*CommitFile, error) {
	log.Printf("Getting file: %s", path)
	commitFile := &CommitFile{}
	projectId, err := GetProjectId(ctx, c, &repo)
	if err != nil {
		return nil, err
	}
	fmt.Println("got project id: ", *projectId)
	fileContent, resp, err := c.client.RepositoryFiles.GetFile(*projectId, path, &gitlab.GetFileOptions{Ref: &branch}, gitlab.WithContext(*ctx))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, err
	}

	decodedText, err := DecodeBase64ToStringPtr(fileContent.Content)
	if err != nil {
		return nil, err
	}

	commitFile.Path = &fileContent.FilePath
	commitFile.Content = decodedText

	return commitFile, nil
}

func (c *GitlabClientImpl) GetFiles(ctx *context.Context, repo string, branch string, paths []string) ([]*CommitFile, error) {
	log.Println("Getting multiple files")
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
	log.Println("commit file", commitFiles)
	return commitFiles, nil
}

func (c *GitlabClientImpl) SetWebhook(ctx *context.Context, repo *string) (*HookWithStatus, error) {
	var gitlabHookId *int
	if *repo == "" {
		log.Println("starting with group level hooks")
		respHook, ok := IsGroupWebhookEnabled(ctx, c)
		if !ok {
			groupHookOptions := gitlab.AddGroupHookOptions{
				URL:                 &c.cfg.GitProviderConfig.WebhookURL,
				Token:               &c.cfg.GitProviderConfig.WebhookSecret,
				MergeRequestsEvents: gitlab.Ptr(true),
				PushEvents:          gitlab.Ptr(true),
				ReleasesEvents:      gitlab.Ptr(true),
			}

			gitlabHook, resp, err := c.client.Groups.AddGroupHook(c.cfg.GitProviderConfig.OrgName, &groupHookOptions, gitlab.WithContext(*ctx))
			if resp != nil {
				if resp.StatusCode == http.StatusForbidden {
					return nil, fmt.Errorf("for org level webhook, group token must be Owner level")
				} else if resp.StatusCode != http.StatusCreated {
					return nil, fmt.Errorf("failed to create group level webhhok, API returned %d", resp.StatusCode)
				}
			}
			if err != nil {
				return nil, err
			}
			gitlabHookId = &gitlabHook.ID
			log.Printf("added webhook: %d for %s name: %s\n", gitlabHook.ID, c.cfg.GitProviderConfig.OrgName, gitlabHook.URL)
		} else {
			editedGroupHookOpt := gitlab.EditGroupHookOptions{
				URL:                 gitlab.Ptr(c.cfg.GitProviderConfig.WebhookURL),
				Token:               gitlab.Ptr(c.cfg.GitProviderConfig.WebhookSecret),
				MergeRequestsEvents: gitlab.Ptr(true),
				PushEvents:          gitlab.Ptr(true),
				ReleasesEvents:      gitlab.Ptr(true),
			}
			gitlabHook, resp, err := c.client.Groups.EditGroupHook(c.cfg.GitProviderConfig.OrgName, respHook.ID, &editedGroupHookOpt, gitlab.WithContext(*ctx))
			if resp != nil {
				if resp.StatusCode == http.StatusForbidden {
					return nil, fmt.Errorf("for org level webhook, group token must be Owner level")
				} else if resp.StatusCode != http.StatusOK {
					fmt.Println(resp.Request.URL, err)
					return nil, fmt.Errorf(
						"failed to update group level webhook for %s, API returned %d",
						c.cfg.GitProviderConfig.OrgName,
						resp.StatusCode,
					)
				}
			}
			if err != nil {
				return nil, err
			}
			gitlabHookId = &gitlabHook.ID
			log.Printf("edited webhook for %s: %s\n", c.cfg.GitProviderConfig.OrgName, gitlabHook.URL)
		}
	} else {
		projectId, err := GetProjectId(ctx, c, repo)
		if err != nil {
			return nil, err
		}
		log.Printf("project id is: %d\n", *projectId)
		respHook, ok := IsProjectWebhookEnabled(ctx, c, *projectId)

		if !ok {
			addProjectHookOpts := gitlab.AddProjectHookOptions{
				URL:                 &c.cfg.GitProviderConfig.WebhookURL,
				Token:               &c.cfg.GitProviderConfig.WebhookSecret,
				MergeRequestsEvents: gitlab.Ptr(true),
				PushEvents:          gitlab.Ptr(true),
				ReleasesEvents:      gitlab.Ptr(true),
			}
			gitlabHook, resp, err := c.client.Projects.AddProjectHook(*projectId, &addProjectHookOpts, gitlab.WithContext(*ctx))
			if resp != nil {
				if resp.StatusCode == http.StatusForbidden {
					return nil, fmt.Errorf("for projects specific webhook, group token must be Maintainer level or above")
				} else if resp.StatusCode != http.StatusCreated {
					return nil, fmt.Errorf("failed to create repo level webhhok for %s, API returned %d", *repo, resp.StatusCode)
				}
			}
			if err != nil {
				log.Println("url", *addProjectHookOpts.URL, *projectId, *addProjectHookOpts.Token)
				return nil, fmt.Errorf("failed to add project hook ,%d", err)
			}
			gitlabHookId = &gitlabHook.ID
			log.Printf("created webhook: %d for %s: %s\n", gitlabHook.ID, *repo, gitlabHook.URL)
		} else {
			editProjectHookOpts := gitlab.EditProjectHookOptions{
				URL:                 gitlab.Ptr(c.cfg.GitProviderConfig.WebhookURL),
				Token:               gitlab.Ptr(c.cfg.GitProviderConfig.WebhookSecret),
				MergeRequestsEvents: gitlab.Ptr(true),
				PushEvents:          gitlab.Ptr(true),
				ReleasesEvents:      gitlab.Ptr(true),
			}
			gitlabHook, resp, err := c.client.Projects.EditProjectHook(*projectId, respHook.ID, &editProjectHookOpts, gitlab.WithContext(*ctx))
			if resp != nil {
				if resp.StatusCode == http.StatusForbidden {
					return nil, fmt.Errorf("for projects specific webhook, group token must be Maintainer level or above")
				} else if resp.StatusCode != http.StatusOK {
					return nil, fmt.Errorf("failed to update repo level webhhok for %s, API returned %d", *repo, resp.StatusCode)
				}
			}
			if err != nil {
				return nil, err
			}
			gitlabHookId = &gitlabHook.ID
			log.Printf("edited webhook: %d for %s: %s\n", *gitlabHookId, *repo, gitlabHook.URL)
		}

	}

	hookID := int64(*gitlabHookId)
	return &HookWithStatus{HookID: hookID, HealthStatus: true, RepoName: repo}, nil
}

func (c *GitlabClientImpl) UnsetWebhook(ctx *context.Context, hook *HookWithStatus) error {
	log.Println("unsetting webhook")
	if hook.RepoName == nil {
		resp, err := c.client.Groups.DeleteGroupHook(c.cfg.GitProviderConfig.OrgName, int(hook.HookID), gitlab.WithContext(*ctx))
		if resp != nil {
			if resp.StatusCode != http.StatusNoContent {
				return fmt.Errorf("failed to delete group level webhhok, API call returned %d", resp.StatusCode)
			}
		}
		if err != nil {
			return err
		}
		log.Printf("removed group webhook, hookID :%d\n", hook.HookID)
	} else {
		resp, err := c.client.Projects.DeleteProjectHook(*hook.RepoName, int(hook.HookID), gitlab.WithContext(*ctx))
		if resp != nil {
			if resp.StatusCode != http.StatusNoContent {
				return fmt.Errorf("failed to delete project level webhhok for %s, API call returned %d", *hook.RepoName, resp.StatusCode)
			}
		}
		if err != nil {
			return fmt.Errorf("failed to delete project level webhhok for %s, API call returned %s", *hook.RepoName, err)
		}
		log.Printf("removed project webhook, project:%s hookID :%d\n", *hook.RepoName, hook.HookID) // INFO
	}

	return nil
}

func (c *GitlabClientImpl) HandlePayload(ctx *context.Context, request *http.Request, secret []byte) (*WebhookPayload, error) {
	log.Printf("starting with payload")
	var webhookPayload WebhookPayload
	payload, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading  request body: %v", err)
	}
	event, err := gitlab.ParseWebhook(gitlab.WebhookEventType(request), payload)
	if err != nil {
		return nil, err
	}
	switch e := event.(type) {
	case *gitlab.PushEvent:
		webhookPayload = WebhookPayload{
			Event:     "push",
			Repo:      e.Project.Name,
			Branch:    strings.TrimPrefix(e.Ref, "refs/heads/"),
			Commit:    e.CheckoutSHA,
			User:      e.UserName,
			UserEmail: e.UserEmail,
			OwnerID:   int64(e.UserID),
		}
	case *gitlab.MergeEvent:
		webhookPayload = WebhookPayload{
			Event:            "merge_request",
			Action:           e.ObjectAttributes.Action, //open, close, reopen, update, approved, unapproved, approval, unapproval, merge
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
	case *gitlab.ReleaseEvent:
		webhookPayload = WebhookPayload{
			Event:     "release",
			Action:    e.Action, // "create" | "update" | "delete"
			Repo:      e.Project.Name,
			Branch:    e.Tag,
			Commit:    e.Commit.ID,
			User:      e.Commit.Author.Name,
			UserEmail: e.Commit.Author.Email,
		}
	}
	log.Printf("sending payload: %s, %s", webhookPayload.Repo, webhookPayload.User)
	return &webhookPayload, nil
}

func (c *GitlabClientImpl) SetStatus(ctx *context.Context, repo *string, commit *string, linkURL *string, status *string, message *string) error {
	if !utils.ValidateHTTPFormat(*linkURL) {
		log.Println("invalid link URL", *linkURL)
		return fmt.Errorf("invalid linkURL")
	}
	projectId, err := GetProjectId(ctx, c, repo)
	if err != nil {
		return err
	}

	currCommit, resp, err := c.client.Commits.GetCommitStatuses(*projectId, *commit, nil, gitlab.WithContext(*ctx))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get commit status on repo:%s, commit:%s, API call returned %d", *repo, *commit, resp.StatusCode)
	}

	if len(currCommit) != 0 {
		if currCommit[0].Status == *status {
			// https://forum.gitlab.com/t/cannot-transition-status-via-run-from-running-reason-s-status-cannot-transition-via-run/42588/6
			log.Printf("cannot change commit description without status also, status stays: %s", *status)
			return nil
		}
	}

	repoStatus := gitlab.SetCommitStatusOptions{
		State:       gitlab.BuildStateValue(*status), // pending, success, error, or failure.
		TargetURL:   linkURL,
		Description: gitlab.Ptr(fmt.Sprintf("Workflow %s %s", *status, *message)),
		Context:     gitlab.Ptr("Piper/ArgoWorkflows"),
	}
	_, resp, err = c.client.Commits.SetCommitStatus(*projectId, *commit, &repoStatus, gitlab.WithContext(*ctx))
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
	//TODO implement me
	panic("implement me")
}
