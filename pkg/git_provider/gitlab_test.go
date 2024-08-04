package git_provider

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/quickube/piper/pkg/conf"
	"github.com/quickube/piper/pkg/utils"
	assertion "github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/net/context"
)

func TestGitlabListFiles(t *testing.T) {
	// Prepare
	mux, client := setupGitlab(t)

	repoContent := &gitlab.TreeNode{
		Type: "file",
		Name: "exit.yaml",
		Path: ".workflows/exit.yaml",
	}

	repoContent2 := &gitlab.TreeNode{
		Type: "file",
		Name: "main.yaml",
		Path: ".workflows/main.yaml",
	}

	treeNodes := []gitlab.TreeNode{*repoContent, *repoContent2}
	expectedRef := "branch1"
	project := "project1"

	c := GitlabClientImpl{
		client: client,
		cfg: &conf.GlobalConfig{
			GitProviderConfig: conf.GitProviderConfig{
				OrgLevelWebhook: true,
				OrgName:         "group1",
				RepoList:        project,
			},
		},
	}
	projectUrl := fmt.Sprintf("/api/v4/projects/%s/%s/repository/tree",c.cfg.GitProviderConfig.OrgName, project)
	mux.HandleFunc(projectUrl, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		ref := r.URL.Query().Get("ref")

		// Check if the ref value matches the expected value
		if ref != expectedRef {
			http.Error(w, "Invalid ref value", http.StatusBadRequest)
			return
		}
		mockHTTPResponse(t, w,  treeNodes)
	})

	
	ctx := context.Background()

	// Execute
	actualContent, err := c.ListFiles(&ctx, project, expectedRef, ".workflows")

	var expectedFilesNames []string
	for _, file := range treeNodes{
		expectedFilesNames = append(expectedFilesNames, file.Name)
	}

	// Assert
	assert := assertion.New(t)
	assert.NotNil(t, err)
	assert.Equal(expectedFilesNames, actualContent)
}

func TestGitlabGetFile(t *testing.T) {
	// Prepare
	mux, client := setupGitlab(t)


	expectedFile := gitlab.File{
			Content: "file",
			FileName: "file.yaml",
			CommitID: "1",
			FilePath: ".workflows/file.yaml",
		}

	c := GitlabClientImpl{
		client: client,
		cfg: &conf.GlobalConfig{
			GitProviderConfig: conf.GitProviderConfig{
				OrgLevelWebhook: true,
				OrgName:         "group1",
				RepoList:        "project1",
			},
		},
	}

	mux.HandleFunc("/api/v4/projects/project1/repository/files/.workflows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		ref := r.URL.Query().Get("ref")
		// Check if the ref value matches the expected value
		if ref != "branch1" {
			http.Error(w, "Invalid ref value", http.StatusBadRequest)
			return
		}
		
		mockHTTPResponse(t, w,  expectedFile)
	})

	
	ctx := context.Background()

	// Execute
	actualFile, err := c.GetFile(&ctx, "project1", "branch1", ".workflows")

	// Assert
	assert := assertion.New(t)
	assert.NotNil(t, err)

	assert.Equal(*actualFile.Path, expectedFile.FilePath)
	assert.Equal(*actualFile.Content, expectedFile.Content)
}

func TestGitlabPingHook(t *testing.T) {
	// Prepare
	ctx := context.Background()
	assert := assertion.New(t)
	mux, client := setupGitlab(t)

	hookUrl := "https://url"

	// Test-repo2 existing webhook
	mux.HandleFunc("/api/v4/projects/test-repo1/hooks/234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		mockHTTPResponse(t, w, nil)
	})

	mux.HandleFunc("/api/v4/groups/test/hooks/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		mockHTTPResponse(t, w, nil)

	})

	c := GitlabClientImpl{
		client: client,
		cfg: &conf.GlobalConfig{
			GitProviderConfig: conf.GitProviderConfig{},
		},
	}

	// Define test cases
	tests := []struct {
		name        string
		repo        *string
		hook        *HookWithStatus
		config      *conf.GitProviderConfig
	}{
		{
			name: "Ping repo webhook",
			hook: &HookWithStatus{
				HookID:       234,
				HealthStatus: true,
				RepoName:     utils.SPtr("test-repo1"),
			},
			config: &conf.GitProviderConfig{
				OrgLevelWebhook: false,
				OrgName:         "test",
				WebhookURL:      hookUrl,
			},
		},
		{
			name: "Ping org webhook",
			hook: &HookWithStatus{
				HookID:       123,
				HealthStatus: true,
				RepoName:     nil,
			},
			config: &conf.GitProviderConfig{
				OrgLevelWebhook: true,
				OrgName:         "test",
				WebhookURL:      hookUrl,
			},
		},
	}
	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c.cfg.GitProviderConfig = *test.config
			// Call the function being tested
			err := c.PingHook(&ctx, test.hook)

			assert.Nil(err)


		})
	}
}

func TestGitlabSetStatus(t *testing.T) {
	// Prepare
	ctx := context.Background()
	assert := assertion.New(t)
	mux, client := setupGitlab(t)

	mux.HandleFunc("/api/v4/projects/test-repo1/statuses/test-commit", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		w.WriteHeader(http.StatusCreated)
		jsonBytes := []byte(`{"status": "ok"}`)
		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

	c := GitlabClientImpl{
		client: client,
		cfg: &conf.GlobalConfig{
			GitProviderConfig: conf.GitProviderConfig{
				OrgLevelWebhook: false,
				OrgName:         "test",
				RepoList:        "test-repo1",
			},
		},
	}

	// Define test cases
	tests := []struct {
		name        string
		repo        *string
		commit      *string
		linkURL     *string
		status      *string
		message     *string
		wantedError error
	}{
		{
			name:        "Notify success",
			repo:        utils.SPtr("test-repo1"),
			commit:      utils.SPtr("test-commit"),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("success"),
			message:     utils.SPtr(""),
			wantedError: nil,
		},
		{
			name:        "Notify pending",
			repo:        utils.SPtr("test-repo1"),
			commit:      utils.SPtr("test-commit"),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("pending"),
			message:     utils.SPtr(""),
			wantedError: nil,
		},
		{
			name:        "Notify error",
			repo:        utils.SPtr("test-repo1"),
			commit:      utils.SPtr("test-commit"),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("error"),
			message:     utils.SPtr("some message"),
			wantedError: nil,
		},
		{
			name:        "Notify failure",
			repo:        utils.SPtr("test-repo1"),
			commit:      utils.SPtr("test-commit"),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("failure"),
			message:     utils.SPtr(""),
			wantedError: nil,
		},
		{
			name:        "Non managed repo",
			repo:        utils.SPtr("non-existing-repo"),
			commit:      utils.SPtr("test-commit"),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("error"),
			message:     utils.SPtr(""),
			wantedError: errors.New("some error"),
		},
		{
			name:        "Non existing commit",
			repo:        utils.SPtr("test-repo1"),
			commit:      utils.SPtr("not-exists"),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("error"),
			message:     utils.SPtr(""),
			wantedError: errors.New("some error"),
		},
		{
			name:        "Wrong URL",
			repo:        utils.SPtr("test-repo1"),
			commit:      utils.SPtr("test-commit"),
			linkURL:     utils.SPtr("argo"),
			status:      utils.SPtr("error"),
			message:     utils.SPtr(""),
			wantedError: errors.New("some error"),
		},
	}
	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Call the function being tested
			err := c.SetStatus(&ctx, test.repo, test.commit, test.linkURL, test.status, test.message)

			// Use assert to check the equality of the error
			if test.wantedError != nil {
				assert.Error(err)
				assert.NotNil(err)
			} else {
				assert.NoError(err)
				assert.Nil(err)
			}
		})
	}
}

func TestGitlabSetWebhook(t *testing.T) {
	// Prepare
	ctx := context.Background()
	assert := assertion.New(t)
	mux, client := setupGitlab(t)

	hookUrl := "https://url"

	// new group webhook
	mux.HandleFunc("/api/v4/groups/groupA/hooks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			mockHTTPResponse(t,w,[]*gitlab.GroupHook{})
		} else if r.Method == "POST"{
			w.WriteHeader(http.StatusCreated)
			mockHTTPResponse(t,w,gitlab.GroupHook{ID:123,URL: hookUrl})
		}
	})
	// existing group Webhook
	mux.HandleFunc("/api/v4/groups/groupB/hooks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			mockHTTPResponse(t,w,[]*gitlab.GroupHook{{ID:123,URL: hookUrl}})
		}
	})
	mux.HandleFunc("/api/v4/groups/groupB/hooks/123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT"{
			w.WriteHeader(http.StatusOK)
			mockHTTPResponse(t,w,gitlab.GroupHook{ID:123,URL: hookUrl})
		}
	})

	// new project Webhook
	mux.HandleFunc("/api/v4/projects/test/test-repo1/hooks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			mockHTTPResponse(t,w,[]*gitlab.ProjectHook{{}})
		}		
	})
	mux.HandleFunc("/api/v4/projects/test-repo1/hooks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
			mockHTTPResponse(t,w, gitlab.ProjectHook{ID:123,URL: hookUrl})
		}		
	})
	// existing project webhook
	mux.HandleFunc("/api/v4/projects/test/test-repo2/hooks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			mockHTTPResponse(t,w,[]*gitlab.ProjectHook{{ID:123,URL: hookUrl}})
		}
	})
	mux.HandleFunc("/api/v4/projects/test-repo2/hooks/123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			w.WriteHeader(http.StatusOK)
			mockHTTPResponse(t,w, gitlab.ProjectHook{ID:123,URL: hookUrl})
		}		
	})



	c := GitlabClientImpl{
		client: client,
		cfg: &conf.GlobalConfig{
			GitProviderConfig: conf.GitProviderConfig{},
		},
	}

	// Define test cases
	tests := []struct {
		name        string
		repo        *string
		config      *conf.GitProviderConfig
		wantedError error
	}{
		{
			name: "Repo and orgWebHook enabled",
			repo: utils.SPtr("repo"),
			config: &conf.GitProviderConfig{
				OrgLevelWebhook: true,
				OrgName:         "test",
				RepoList:        "test-repo1",
				WebhookURL:      hookUrl,
			},
			wantedError: errors.New("error"),
		},
		{
			name: "New group webhook",
			repo: nil,
			config: &conf.GitProviderConfig{
				OrgLevelWebhook: true,
				OrgName:         "groupA",
				RepoList:        "test-repo1",
				WebhookURL:      hookUrl,
			},
			wantedError: nil,
		},
		{
			name: "Existing group webhook",
			repo: nil,
			config: &conf.GitProviderConfig{
				OrgLevelWebhook: true,
				OrgName:         "groupB",
				RepoList:        "test-repo1",
				WebhookURL:      hookUrl,
			},
			wantedError: nil,
		},
		{
			name: "New project webhook",
			repo: utils.SPtr("test-repo1"),
			config: &conf.GitProviderConfig{
				OrgLevelWebhook: false,
				OrgName:         "test",
				RepoList:        "test-repo1",
				WebhookURL:      hookUrl,
			},
			wantedError: nil,
		},
		{
			name: "Existing project webhook",
			repo: utils.SPtr("test-repo2"),
			config: &conf.GitProviderConfig{
				OrgLevelWebhook: false,
				OrgName:         "test",
				RepoList:        "test-repo2",
				WebhookURL:      hookUrl,
			},
			wantedError: nil,
		},
		
	}
	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c.cfg.GitProviderConfig = *test.config
			// Call the function being tested
			_, err := c.SetWebhook(&ctx, test.repo)

			// Use assert to check the equality of the error
			if test.wantedError != nil {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				//assert.Equal(hookUrl, hook.Config["url"])
			}
		})
	}
}
