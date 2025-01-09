package git_provider

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/quickube/piper/pkg/conf"
	"github.com/quickube/piper/pkg/utils"
	assertion "github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/net/context"
	"net/http"
	"testing"
)

func TestGitlabListFiles(t *testing.T) {
	// Prepare
	mux, client := setupGitlab(t)

	repoContent := gitlab.TreeNode{
		Type: "file",
		Name: "exit.yaml",
		Path: ".workflows/exit.yaml",
	}

	repoContent2 := gitlab.TreeNode{
		Type: "file",
		Name: "main.yaml",
		Path: ".workflows/main.yaml",
	}

	treeNodes := []gitlab.TreeNode{repoContent, repoContent2}
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
	mux.HandleFunc("/api/v4/projects/1/repository/tree", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		ref := r.URL.Query().Get("ref")

		// Check if the ref value matches the expected value
		if ref != expectedRef {
			http.Error(w, "Invalid ref value", http.StatusBadRequest)
			return
		}
		mockHTTPResponse(t, w, treeNodes)
	})
	url := fmt.Sprintf("/api/v4/projects/%s/%s", c.cfg.GitProviderConfig.OrgName, project)
	mockProject := gitlab.Project{ID: 1}
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		mockHTTPResponse(t, w, mockProject)
	})

	ctx := context.Background()

	// Execute
	actualContent, err := c.ListFiles(ctx, project, expectedRef, ".workflows")

	var expectedFilesNames []string
	for _, file := range treeNodes {
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
	project := "project1"
	fileName := "file.yaml"
	filePath := fmt.Sprintf(".workflows/%s", fileName)
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
	branch := "branch1"
	projectUrl := fmt.Sprintf("/api/v4/projects/%s/%s", c.cfg.GitProviderConfig.OrgName, project)
	mockProject := &gitlab.Project{ID: 1}
	mux.HandleFunc(projectUrl, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		mockHTTPResponse(t, w, mockProject)
	})
	decodedString := "file data"
	encoded := base64.StdEncoding.EncodeToString([]byte(decodedString))

	expectedFile := gitlab.File{
		Content:  encoded,
		FileName: fileName,
		CommitID: "1",
		FilePath: filePath,
	}
	url := fmt.Sprintf("/api/v4/projects/%d/repository/files/%s", mockProject.ID, filePath)
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("ref") != branch {
			t.Errorf("Unexpected request: %s", r.URL.String())
		}
		testMethod(t, r, "GET")
		mockHTTPResponse(t, w, expectedFile)
	})

	ctx := context.Background()

	// Execute
	actualFile, err := c.GetFile(ctx, project, branch, filePath)
	// Assert
	assert := assertion.New(t)
	assert.NotNil(t, err)
	assert.Equal(*actualFile.Path, expectedFile.FilePath)
	assert.Equal(*actualFile.Content, decodedString)
}

func TestGitlabSetStatus(t *testing.T) {
	// Prepare
	ctx := context.Background()
	assert := assertion.New(t)
	mux, client := setupGitlab(t)

	project := "test-repo1"
	commit := "test-commit"
	c := GitlabClientImpl{
		client: client,
		cfg: &conf.GlobalConfig{
			GitProviderConfig: conf.GitProviderConfig{
				OrgLevelWebhook: false,
				OrgName:         "test",
				RepoList:        project,
			},
		},
	}
	projectUrl := fmt.Sprintf("/api/v4/projects/%s/%s", c.cfg.GitProviderConfig.OrgName, project)
	mockProject := &gitlab.Project{ID: 1}
	mux.HandleFunc(projectUrl, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		mockHTTPResponse(t, w, mockProject)
	})
	currCommitUrl := fmt.Sprintf("/api/v4/projects/%d/repository/commits/%s/statuses", mockProject.ID, commit)

	mux.HandleFunc(currCommitUrl, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		mockHTTPResponse(t, w, []gitlab.CommitStatus{})
	})
	mux.HandleFunc("/api/v4/projects/1/statuses/test-commit", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		w.WriteHeader(http.StatusCreated)
		jsonBytes := []byte(`{"status": "ok"}`)
		_, _ = fmt.Fprint(w, string(jsonBytes))
	})

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
			repo:        utils.SPtr(project),
			commit:      utils.SPtr(commit),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("success"),
			message:     utils.SPtr(""),
			wantedError: nil,
		},
		{
			name:        "Notify pending",
			repo:        utils.SPtr(project),
			commit:      utils.SPtr(commit),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("pending"),
			message:     utils.SPtr(""),
			wantedError: nil,
		},
		{
			name:        "Notify error",
			repo:        utils.SPtr(project),
			commit:      utils.SPtr(commit),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("error"),
			message:     utils.SPtr("some message"),
			wantedError: nil,
		},
		{
			name:        "Notify failure",
			repo:        utils.SPtr(project),
			commit:      utils.SPtr(commit),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("failure"),
			message:     utils.SPtr(""),
			wantedError: nil,
		},
		{
			name:        "Non managed repo",
			repo:        utils.SPtr("non-existing-repo"),
			commit:      utils.SPtr(commit),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("error"),
			message:     utils.SPtr(""),
			wantedError: errors.New("404 Not Found"),
		},
		{
			name:        "Non existing commit",
			repo:        utils.SPtr(project),
			commit:      utils.SPtr("not-exists"),
			linkURL:     utils.SPtr("https://argo"),
			status:      utils.SPtr("error"),
			message:     utils.SPtr(""),
			wantedError: errors.New("404 Not Found"),
		},
		{
			name:        "Wrong URL",
			repo:        utils.SPtr(project),
			commit:      utils.SPtr(commit),
			linkURL:     utils.SPtr("argo"),
			status:      utils.SPtr("error"),
			message:     utils.SPtr(""),
			wantedError: errors.New("invalid linkURL"),
		},
	}
	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Call the function being tested
			err := c.SetStatus(ctx, test.repo, test.commit, test.linkURL, test.status, test.message)

			if test.wantedError != nil {
				assert.NotNil(err)
				assert.Equal(test.wantedError.Error(), err.Error())
			} else {
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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.String() == "/api/v4/groups/groupA/hooks" {
			// new group webhook check for existing
			mockHTTPResponse(t, w, []*gitlab.GroupHook{})
		} else if r.Method == "POST" && r.URL.String() == "/api/v4/groups/groupA/hooks" {
			// new group webhook creation of new webhook
			w.WriteHeader(http.StatusCreated)
			mockHTTPResponse(t, w, gitlab.GroupHook{ID: 123, URL: hookUrl})
		} else if r.Method == "GET" && r.URL.String() == "/api/v4/groups/groupB/hooks" {
			// existing group Webhook check for existing
			mockHTTPResponse(t, w, []*gitlab.GroupHook{{ID: 123, URL: hookUrl}})
		} else if r.Method == "PUT" && r.URL.String() == "/api/v4/groups/groupB/hooks/123" {
			// existing group Webhook editing the existing one
			w.WriteHeader(http.StatusOK)
			mockHTTPResponse(t, w, gitlab.GroupHook{ID: 123, URL: hookUrl})
		} else if r.Method == "GET" && r.URL.String() == "/api/v4/projects/test%2Ftest-repo1" {
			// new project Webhook get project id
			mockHTTPResponse(t, w, &gitlab.Project{ID: 1})
		} else if r.Method == "GET" && r.URL.String() == "/api/v4/projects/test%2Ftest-repo2" {
			// new project Webhook get project id
			mockHTTPResponse(t, w, &gitlab.Project{ID: 2})
		} else if r.Method == "GET" && r.URL.String() == "/api/v4/projects/1/hooks" {
			// new project Webhook check for existing
			mockHTTPResponse(t, w, []*gitlab.ProjectHook{{}})
		} else if r.Method == "POST" && r.URL.String() == "/api/v4/projects/1/hooks" {
			// new project Webhook create new webhook
			w.WriteHeader(http.StatusCreated)
			mockHTTPResponse(t, w, gitlab.ProjectHook{ID: 123, URL: hookUrl})
		} else if r.Method == "GET" && r.URL.String() == "/api/v4/projects/2/hooks" {
			// new project Webhook check for existing
			mockHTTPResponse(t, w, []*gitlab.ProjectHook{{ID: 123, URL: hookUrl}})
		} else if r.Method == "PUT" && r.URL.String() == "/api/v4/projects/2/hooks/123" {
			// new project Webhook edit existing webhook
			w.WriteHeader(http.StatusOK)
			mockHTTPResponse(t, w, gitlab.ProjectHook{ID: 123, URL: hookUrl})
		} else {
			fmt.Println("unhandled ", r.Method, " route: ", r.URL.String())
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
		name   string
		repo   *string
		config conf.GitProviderConfig
	}{
		{
			name: "New group webhook",
			repo: nil,
			config: conf.GitProviderConfig{
				OrgLevelWebhook: true,
				OrgName:         "groupA",
				RepoList:        "",
				WebhookURL:      hookUrl,
			},
		},
		{
			name: "Existing group webhook",
			repo: nil,
			config: conf.GitProviderConfig{
				OrgLevelWebhook: true,
				OrgName:         "groupB",
				RepoList:        "",
				WebhookURL:      hookUrl,
			},
		},
		{
			name: "New project webhook",
			repo: utils.SPtr("test-repo1"),
			config: conf.GitProviderConfig{
				OrgLevelWebhook: false,
				OrgName:         "test",
				RepoList:        "test-repo1",
				WebhookURL:      hookUrl,
			},
		},
		{
			name: "Existing project webhook",
			repo: utils.SPtr("test-repo2"),
			config: conf.GitProviderConfig{
				OrgLevelWebhook: false,
				OrgName:         "test",
				RepoList:        "test-repo2",
				WebhookURL:      hookUrl,
			},
		},
	}
	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c.cfg.GitProviderConfig = test.config
			_, err := c.SetWebhook(ctx, &c.cfg.GitProviderConfig.RepoList)

			// Use assert to check the equality of the error
			assert.Nil(err)
		})
	}
}
