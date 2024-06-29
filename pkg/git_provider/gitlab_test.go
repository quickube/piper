package git_provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/quickube/piper/pkg/conf"
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
