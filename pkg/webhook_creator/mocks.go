package webhook_creator

import (
	context2 "context"
	"errors"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/quickube/piper/pkg/git_provider"
	"golang.org/x/net/context"
	"net/http"
)

type MockGitProviderClient struct {
	ListFilesFunc           func(ctx context.Context, repo string, branch string, path string) ([]string, error)
	GetFileFunc             func(ctx context.Context, repo string, branch string, path string) (*git_provider.CommitFile, error)
	GetFilesFunc            func(ctx context.Context, repo string, branch string, paths []string) ([]*git_provider.CommitFile, error)
	SetWebhookFunc          func(ctx context.Context, repo *string) (*git_provider.HookWithStatus, error)
	UnsetWebhookFunc        func(ctx context.Context, hook *git_provider.HookWithStatus) error
	HandlePayloadFunc       func(request *http.Request, secret []byte) (*git_provider.WebhookPayload, error)
	SetStatusFunc           func(ctx context.Context, repo *string, commit *string, linkURL *string, status *string, message *string) error
	PingHookFunc            func(ctx context.Context, hook *git_provider.HookWithStatus) error
	GetCorrelatingEventFunc func(ctx context.Context, workflowEvent *v1alpha1.WorkflowPhase) (string, error)
}

func (m *MockGitProviderClient) ListFiles(ctx context2.Context, repo string, branch string, path string) ([]string, error) {
	if m.ListFilesFunc != nil {
		return m.ListFilesFunc(ctx, repo, branch, path)
	}
	return nil, errors.New("unimplemented")
}

func (m *MockGitProviderClient) GetFile(ctx context2.Context, repo string, branch string, path string) (*git_provider.CommitFile, error) {
	if m.GetFileFunc != nil {
		return m.GetFileFunc(ctx, repo, branch, path)
	}
	return nil, errors.New("unimplemented")
}

func (m *MockGitProviderClient) GetFiles(ctx context2.Context, repo string, branch string, paths []string) ([]*git_provider.CommitFile, error) {
	if m.GetFilesFunc != nil {
		return m.GetFilesFunc(ctx, repo, branch, paths)
	}
	return nil, errors.New("unimplemented")
}

func (m *MockGitProviderClient) SetWebhook(ctx context2.Context, repo *string) (*git_provider.HookWithStatus, error) {
	if m.SetWebhookFunc != nil {
		return m.SetWebhookFunc(ctx, repo)
	}
	return nil, errors.New("unimplemented")
}

func (m *MockGitProviderClient) UnsetWebhook(ctx context2.Context, hook *git_provider.HookWithStatus) error {
	if m.UnsetWebhookFunc != nil {
		return m.UnsetWebhookFunc(ctx, hook)
	}
	return errors.New("unimplemented")
}

func (m *MockGitProviderClient) HandlePayload(ctx context2.Context, request *http.Request, secret []byte) (*git_provider.WebhookPayload, error) {
	if m.HandlePayloadFunc != nil {
		return m.HandlePayloadFunc(request, secret)
	}
	return nil, errors.New("unimplemented")
}

func (m *MockGitProviderClient) SetStatus(ctx context2.Context, repo *string, commit *string, linkURL *string, status *string, message *string) error {
	if m.SetStatusFunc != nil {
		return m.SetStatusFunc(ctx, repo, commit, linkURL, status, message)
	}
	return errors.New("unimplemented")
}
func (m *MockGitProviderClient) GetCorrelatingEvent(ctx context.Context, workflowEvent *v1alpha1.WorkflowPhase) (string, error) {
	if m.GetCorrelatingEventFunc != nil {
		return m.GetCorrelatingEventFunc(ctx, workflowEvent)
	}
	return "", errors.New("unimplemented")
}

func (m *MockGitProviderClient) PingHook(ctx context2.Context, hook *git_provider.HookWithStatus) error {
	if m.PingHookFunc != nil {
		return m.PingHookFunc(ctx, hook)
	}
	return errors.New("unimplemented")
}
