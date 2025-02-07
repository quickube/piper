package webhook_handler

import (
	"context"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/quickube/piper/pkg/clients"
	"github.com/quickube/piper/pkg/common"
	"github.com/quickube/piper/pkg/git_provider"
	"github.com/quickube/piper/pkg/utils"
	assertion "github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

var fileContentMap = map[string]*string{
	"main.yaml":       utils.SPtr("main.yaml"),
	"exit.yaml":       utils.SPtr("exit.yaml"),
	"parameters.yaml": utils.SPtr("parameters.yaml"),
}

var commitFileMap = map[string]*git_provider.CommitFile{
	"repo1/branch1/.workflows/main.yaml": &git_provider.CommitFile{
		Path:    utils.SPtr(".workflows/main.yaml"),
		Content: fileContentMap["main.yaml"],
	},
	"repo1/branch1/.workflows/exit.yaml": &git_provider.CommitFile{
		Path:    utils.SPtr(".workflows/exit.yaml"),
		Content: fileContentMap["exit.yaml"],
	},
	"repo1/branch2/.workflows/main.yaml": &git_provider.CommitFile{
		Path:    utils.SPtr(".workflows/main.yaml"),
		Content: fileContentMap["main.yaml"],
	},
	"repo1/branch2/.workflows/parameters.yaml": &git_provider.CommitFile{
		Path:    utils.SPtr(".workflows/parameters.yaml"),
		Content: fileContentMap["parameters.yaml"],
	},
}

// mockGitProvider is a mock implementation of the git_provider.Client interface.
type mockGitProvider struct{}

func (m *mockGitProvider) GetFile(ctx context.Context, repo string, branch string, path string) (*git_provider.CommitFile, error) {

	fullPath := fmt.Sprintf("%s/%s/%s", repo, branch, path)
	if fileInfo, ok := commitFileMap[fullPath]; ok {
		return fileInfo, nil
	}
	return &git_provider.CommitFile{}, nil
}

func (m *mockGitProvider) GetFiles(ctx context.Context, repo string, branch string, paths []string) ([]*git_provider.CommitFile, error) {
	var commitFiles []*git_provider.CommitFile

	for _, path := range paths {
		f, err := m.GetFile(ctx, repo, branch, path)
		if err != nil {
			return nil, err
		}
		commitFiles = append(commitFiles, f)
	}

	return commitFiles, nil

}

func (m *mockGitProvider) ListFiles(ctx context.Context, repo string, branch string, path string) ([]string, error) {
	var files []string

	fullPath := fmt.Sprintf("%s/%s/%s/", repo, branch, path)

	for key := range commitFileMap {
		if strings.Contains(key, fullPath) {
			trimmed := strings.Replace(key, fullPath, "", -1)
			files = append(files, trimmed)
		}
	}

	return files, nil
}

func (m *mockGitProvider) SetWebhook(ctx context.Context, repo *string) (*git_provider.HookWithStatus, error) {
	return nil, nil
}

func (m *mockGitProvider) UnsetWebhook(ctx context.Context, hook *git_provider.HookWithStatus) error {
	return nil
}

func (m *mockGitProvider) HandlePayload(ctx context.Context, request *http.Request, secret []byte) (*git_provider.WebhookPayload, error) {
	return nil, nil
}

func (m *mockGitProvider) SetStatus(ctx context.Context, repo *string, commit *string, linkURL *string, status *string, message *string) error {
	return nil
}
func (m *mockGitProvider) GetCorrelatingEvent(ctx context.Context, workflowEvent *v1alpha1.WorkflowPhase) (string, error) {
	return "", nil
}
func (m *mockGitProvider) PingHook(ctx context.Context, hook *git_provider.HookWithStatus) error {
	return nil
}

func TestPrepareBatchForMatchingTriggers(t *testing.T) {
	assert := assertion.New(t)
	ctx := context.Background()
	tests := []struct {
		name                  string
		triggers              *[]Trigger
		payload               *git_provider.WebhookPayload
		expectedWorkflowBatch []*common.WorkflowsBatch
	}{
		{name: "Event without action",
			triggers: &[]Trigger{{
				Events:    &[]string{"event1", "event2.action2"},
				Branches:  &[]string{"branch1", "branch2"},
				Templates: &[]string{""},
				OnStart:   &[]string{"main.yaml"},
				OnExit:    &[]string{"exit.yaml"},
				Config:    "default",
			}},
			payload: &git_provider.WebhookPayload{
				Event:            "event1",
				Action:           "",
				Repo:             "repo1",
				Branch:           "branch1",
				Commit:           "commitHSA",
				User:             "piper",
				UserEmail:        "piper@quickube.com",
				PullRequestURL:   "",
				PullRequestTitle: "",
				DestBranch:       "",
			},
			expectedWorkflowBatch: []*common.WorkflowsBatch{
				&common.WorkflowsBatch{
					OnStart: []*git_provider.CommitFile{
						{
							Path:    utils.SPtr(".workflows/main.yaml"),
							Content: fileContentMap["main.yaml"],
						},
					},
					OnExit: []*git_provider.CommitFile{
						{
							Path:    utils.SPtr(".workflows/exit.yaml"),
							Content: fileContentMap["exit.yaml"],
						},
					},
					Templates: []*git_provider.CommitFile{
						&git_provider.CommitFile{
							Path:    nil,
							Content: nil,
						},
					},
					Parameters: &git_provider.CommitFile{
						Path:    nil,
						Content: nil,
					},
					Config:  utils.SPtr("default"),
					Payload: &git_provider.WebhookPayload{},
				},
			},
		},
		{name: "Event and action",
			triggers: &[]Trigger{{
				Events:    &[]string{"event1", "event2.action2"},
				Branches:  &[]string{"branch1", "branch2"},
				Templates: &[]string{""},
				OnStart:   &[]string{"main.yaml"},
				OnExit:    &[]string{"exit.yaml"},
				Config:    "default",
			}},
			payload: &git_provider.WebhookPayload{
				Event:            "event2",
				Action:           "action2",
				Repo:             "repo1",
				Branch:           "branch1",
				Commit:           "commitHSA",
				User:             "piper",
				UserEmail:        "piper@quickube.com",
				PullRequestURL:   "",
				PullRequestTitle: "",
				DestBranch:       "",
			},
			expectedWorkflowBatch: []*common.WorkflowsBatch{
				&common.WorkflowsBatch{
					OnStart: []*git_provider.CommitFile{
						{
							Path:    utils.SPtr(".workflows/main.yaml"),
							Content: fileContentMap["main.yaml"],
						},
					},
					OnExit: []*git_provider.CommitFile{
						{
							Path:    utils.SPtr(".workflows/exit.yaml"),
							Content: fileContentMap["exit.yaml"],
						},
					},
					Templates: []*git_provider.CommitFile{
						&git_provider.CommitFile{
							Path:    nil,
							Content: nil,
						},
					},
					Parameters: &git_provider.CommitFile{
						Path:    nil,
						Content: nil,
					},
					Config:  utils.SPtr("default"),
					Payload: &git_provider.WebhookPayload{},
				},
			},
		},
		{name: "none existing branch",
			triggers: &[]Trigger{{
				Events:    &[]string{"event1", "event2.action2"},
				Branches:  &[]string{"branch1", "branch2"},
				Templates: &[]string{""},
				OnStart:   &[]string{"main.yaml"},
				OnExit:    &[]string{"exit.yaml"},
				Config:    "default",
			}},
			payload: &git_provider.WebhookPayload{
				Event:            "event1",
				Action:           "branch2",
				Repo:             "repo1",
				Branch:           "branch1",
				Commit:           "commitHSA",
				User:             "piper",
				UserEmail:        "piper@quickube.com",
				PullRequestURL:   "",
				PullRequestTitle: "",
				DestBranch:       "",
			},
			expectedWorkflowBatch: nil,
		},
		{name: "none existing repo",
			triggers: &[]Trigger{{
				Events:    &[]string{"event1", "event2.action2"},
				Branches:  &[]string{"branch1", "branch2"},
				Templates: &[]string{""},
				OnStart:   &[]string{"main.yaml"},
				OnExit:    &[]string{"exit.yaml"},
				Config:    "default",
			}},
			payload: &git_provider.WebhookPayload{
				Event:            "event1",
				Action:           "branch1",
				Repo:             "non-existing",
				Branch:           "branch1",
				Commit:           "commitHSA",
				User:             "piper",
				UserEmail:        "piper@quickube.com",
				PullRequestURL:   "",
				PullRequestTitle: "",
				DestBranch:       "",
			},
			expectedWorkflowBatch: nil,
		},
		{name: "missing OnStart",
			triggers: &[]Trigger{{
				Events:    &[]string{"event1", "event2.action2"},
				Branches:  &[]string{"branch1", "branch2"},
				Templates: &[]string{""},
				OnStart:   &[]string{""},
				OnExit:    &[]string{"exit.yaml"},
				Config:    "default",
			}},
			payload: &git_provider.WebhookPayload{
				Event:            "event1",
				Action:           "branch1",
				Repo:             "repo1",
				Branch:           "branch1",
				Commit:           "commitHSA",
				User:             "piper",
				UserEmail:        "piper@quickube.com",
				PullRequestURL:   "",
				PullRequestTitle: "",
				DestBranch:       "",
			},
			expectedWorkflowBatch: nil,
		},
		{name: "missing OnExit",
			triggers: &[]Trigger{{
				Events:    &[]string{"event1", "event2.action2"},
				Branches:  &[]string{"branch1", "branch2"},
				Templates: &[]string{""},
				OnStart:   &[]string{"main.yaml"},
				OnExit:    &[]string{""},
				Config:    "default",
			}},
			payload: &git_provider.WebhookPayload{
				Event:            "event1",
				Action:           "",
				Repo:             "repo1",
				Branch:           "branch1",
				Commit:           "commitHSA",
				User:             "piper",
				UserEmail:        "piper@quickube.com",
				PullRequestURL:   "",
				PullRequestTitle: "",
				DestBranch:       "",
			},
			expectedWorkflowBatch: []*common.WorkflowsBatch{
				&common.WorkflowsBatch{
					OnStart: []*git_provider.CommitFile{
						{
							Path:    utils.SPtr(".workflows/main.yaml"),
							Content: fileContentMap["main.yaml"],
						},
					},
					OnExit: []*git_provider.CommitFile{
						&git_provider.CommitFile{
							Path:    nil,
							Content: nil,
						},
					},
					Templates: []*git_provider.CommitFile{
						&git_provider.CommitFile{
							Path:    nil,
							Content: nil,
						},
					},
					Parameters: &git_provider.CommitFile{
						Path:    nil,
						Content: nil,
					},
					Config:  utils.SPtr("default"),
					Payload: &git_provider.WebhookPayload{},
				},
			},
		},
		{name: "Multiple OnStart",
			triggers: &[]Trigger{{
				Events:    &[]string{"event1", "event2.action2"},
				Branches:  &[]string{"branch1", "branch2"},
				Templates: &[]string{""},
				OnStart:   &[]string{"main.yaml", "main.yaml"},
				OnExit:    &[]string{""},
				Config:    "default",
			}},
			payload: &git_provider.WebhookPayload{
				Event:            "event1",
				Action:           "",
				Repo:             "repo1",
				Branch:           "branch1",
				Commit:           "commitHSA",
				User:             "piper",
				UserEmail:        "piper@quickube.com",
				PullRequestURL:   "",
				PullRequestTitle: "",
				DestBranch:       "",
			},
			expectedWorkflowBatch: []*common.WorkflowsBatch{
				&common.WorkflowsBatch{
					OnStart: []*git_provider.CommitFile{
						{
							Path:    utils.SPtr(".workflows/main.yaml"),
							Content: fileContentMap["main.yaml"],
						},
						{
							Path:    utils.SPtr(".workflows/main.yaml"),
							Content: fileContentMap["main.yaml"],
						},
					},
					OnExit: []*git_provider.CommitFile{
						&git_provider.CommitFile{
							Path:    nil,
							Content: nil,
						},
					},
					Templates: []*git_provider.CommitFile{
						&git_provider.CommitFile{
							Path:    nil,
							Content: nil,
						},
					},
					Parameters: &git_provider.CommitFile{
						Path:    nil,
						Content: nil,
					},
					Config:  utils.SPtr("default"),
					Payload: &git_provider.WebhookPayload{},
				},
			},
		},
		{name: "Multiple OnExit",
			triggers: &[]Trigger{{
				Events:    &[]string{"event1", "event2.action2"},
				Branches:  &[]string{"branch1", "branch2"},
				Templates: &[]string{""},
				OnStart:   &[]string{"main.yaml"},
				OnExit:    &[]string{"exit.yaml", "exit.yaml"},
				Config:    "default",
			}},
			payload: &git_provider.WebhookPayload{
				Event:            "event1",
				Action:           "",
				Repo:             "repo1",
				Branch:           "branch1",
				Commit:           "commitHSA",
				User:             "piper",
				UserEmail:        "piper@quickube.com",
				PullRequestURL:   "",
				PullRequestTitle: "",
				DestBranch:       "",
			},
			expectedWorkflowBatch: []*common.WorkflowsBatch{
				&common.WorkflowsBatch{
					OnStart: []*git_provider.CommitFile{
						{
							Path:    utils.SPtr(".workflows/main.yaml"),
							Content: fileContentMap["main.yaml"],
						},
					},
					OnExit: []*git_provider.CommitFile{
						&git_provider.CommitFile{
							Path:    utils.SPtr(".workflows/exit.yaml"),
							Content: fileContentMap["exit.yaml"],
						},
						&git_provider.CommitFile{
							Path:    utils.SPtr(".workflows/exit.yaml"),
							Content: fileContentMap["exit.yaml"],
						},
					},
					Templates: []*git_provider.CommitFile{
						&git_provider.CommitFile{
							Path:    nil,
							Content: nil,
						},
					},
					Parameters: &git_provider.CommitFile{
						Path:    nil,
						Content: nil,
					},
					Config:  utils.SPtr("default"),
					Payload: &git_provider.WebhookPayload{},
				},
			},
		},
		{name: "Branch with parameters",
			triggers: &[]Trigger{{
				Events:    &[]string{"event1", "event2.action2"},
				Branches:  &[]string{"branch1", "branch2"},
				Templates: &[]string{""},
				OnStart:   &[]string{"main.yaml"},
				OnExit:    &[]string{""},
				Config:    "default",
			}},
			payload: &git_provider.WebhookPayload{
				Event:            "event1",
				Action:           "",
				Repo:             "repo1",
				Branch:           "branch2",
				Commit:           "commitHSA",
				User:             "piper",
				UserEmail:        "piper@quickube.com",
				PullRequestURL:   "",
				PullRequestTitle: "",
				DestBranch:       "",
			},
			expectedWorkflowBatch: []*common.WorkflowsBatch{
				&common.WorkflowsBatch{
					OnStart: []*git_provider.CommitFile{
						{
							Path:    utils.SPtr(".workflows/main.yaml"),
							Content: fileContentMap["main.yaml"],
						},
					},
					OnExit: []*git_provider.CommitFile{
						&git_provider.CommitFile{
							Path:    nil,
							Content: nil,
						},
					},
					Templates: []*git_provider.CommitFile{
						&git_provider.CommitFile{
							Path:    nil,
							Content: nil,
						},
					},
					Parameters: &git_provider.CommitFile{
						Path:    utils.SPtr(".workflows/parameters.yaml"),
						Content: fileContentMap["parameters.yaml"],
					},
					Config:  utils.SPtr("default"),
					Payload: &git_provider.WebhookPayload{},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wh := &WebhookHandlerImpl{
				Triggers: test.triggers,
				Payload:  test.payload,
				clients: &clients.Clients{
					GitProvider: &mockGitProvider{},
				},
			}
			WorkflowsBatches, err := wh.PrepareBatchForMatchingTriggers(ctx)
			if test.expectedWorkflowBatch == nil {
				assert.Nil(WorkflowsBatches)
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				for iwf, wf := range WorkflowsBatches {
					for i := range wf.OnStart {
						assert.Equal(*test.expectedWorkflowBatch[iwf].OnStart[i].Path, *WorkflowsBatches[iwf].OnStart[i].Path)
						assert.Equal(*test.expectedWorkflowBatch[iwf].OnStart[i].Content, *WorkflowsBatches[iwf].OnStart[i].Content)
					}
					for j := range wf.OnExit {
						if test.expectedWorkflowBatch[iwf].OnExit[j].Path == nil || test.expectedWorkflowBatch[iwf].OnExit[j].Content == nil {
							assert.Nil(WorkflowsBatches[iwf].Templates[j].Path)
							assert.Nil(WorkflowsBatches[iwf].Templates[j].Content)
						} else {
							assert.Equal(*test.expectedWorkflowBatch[iwf].OnExit[j].Path, *WorkflowsBatches[iwf].OnExit[j].Path)
							assert.Equal(*test.expectedWorkflowBatch[iwf].OnExit[j].Content, *WorkflowsBatches[iwf].OnExit[j].Content)
						}
					}

					for k := range wf.Templates {
						if test.expectedWorkflowBatch[iwf].Templates[k].Path == nil || test.expectedWorkflowBatch[iwf].Templates[k].Content == nil {
							assert.Nil(WorkflowsBatches[iwf].Templates[k].Path)
							assert.Nil(WorkflowsBatches[iwf].Templates[k].Content)
						} else {
							assert.Equal(*test.expectedWorkflowBatch[iwf].Templates[k].Path, *WorkflowsBatches[iwf].Templates[k].Path)
							assert.Equal(*test.expectedWorkflowBatch[iwf].Templates[k].Content, *WorkflowsBatches[iwf].Templates[k].Content)
						}

					}

					if test.expectedWorkflowBatch[iwf].Parameters.Path == nil || test.expectedWorkflowBatch[iwf].Parameters.Content == nil {
						assert.Nil(WorkflowsBatches[iwf].Parameters.Path)
						assert.Nil(WorkflowsBatches[iwf].Parameters.Content)
					} else {
						assert.Equal(*test.expectedWorkflowBatch[iwf].Parameters.Path, *WorkflowsBatches[iwf].Parameters.Path)
					}
					assert.Equal(*test.expectedWorkflowBatch[iwf].Config, *WorkflowsBatches[iwf].Config)

				}
			}

		})
	}

}
