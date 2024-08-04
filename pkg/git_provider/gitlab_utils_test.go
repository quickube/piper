package git_provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/quickube/piper/pkg/conf"
	assertion "github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/net/context"
)

func mockHTTPResponse(t *testing.T, w io.Writer, response interface{}) {
	json.NewEncoder(w).Encode(response)
}


func TestValidateGitlabPermissions(t *testing.T){
	//
	// Prepare
	//
	type testData = struct {
		name string
		scopes []string
		raiseErr bool
	}
	var CurrentTest testData
	mux, client := setupGitlab(t)
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
	ctx := context.Background()
	mux.HandleFunc("/api/v4/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		mockHTTPResponse(t, w, gitlab.User{ID:1234})
	})
	mux.HandleFunc("/api/v4/personal_access_tokens", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		mockHTTPResponse(t, w, []gitlab.PersonalAccessToken{{Scopes: CurrentTest.scopes,}})
	})
	//
	// Execute
	//
	tests := []testData{
		{name:"validScope", scopes: []string{"api"}, raiseErr: false},
		{name:"invalidScope", scopes: []string{"invalid"}, raiseErr: true},
	}
	for _, test := range tests {
		CurrentTest = test
		t.Run(test.name, func(t *testing.T) {
			err := ValidateGitlabPermissions(ctx, c.client, c.cfg)
			//
			// Assert
			//
			assert := assertion.New(t)
			if test.raiseErr{
				assert.NotNil(err)
			}else{
				assert.Nil(err)
			}
		})
	}
}

func TestIsGroupWebhookEnabled(t *testing.T){
	//
	// Prepare
	//
	mux, client := setupGitlab(t)
	c := GitlabClientImpl{
		client: client,
		cfg: &conf.GlobalConfig{
			GitProviderConfig: conf.GitProviderConfig{
				OrgLevelWebhook: true,
				OrgName:         "group1",
				WebhookURL: "testing-url",
			},
		},
	}

	hook := []gitlab.GroupHook{{
		ID: 1234,
		URL: c.cfg.GitProviderConfig.WebhookURL,
	},}

	mux.HandleFunc(fmt.Sprintf("/api/v4/groups/%s/hooks", c.cfg.GitProviderConfig.OrgName), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		mockHTTPResponse(t, w, hook)
	})
	//
	// Execute
	//
	groupHook, isEnabled := IsGroupWebhookEnabled(&c)
	//
	// Assert
	//
	assert := assertion.New(t)
	assert.Equal(isEnabled, true)
	assert.Equal(groupHook.URL, c.cfg.GitProviderConfig.WebhookURL)
}

func TestIsProjectWebhookEnabled(t *testing.T){
	//
	// Prepare
	//
	mux, client := setupGitlab(t)
	project := "test-repo1"
	c := GitlabClientImpl{
		client: client,
		cfg: &conf.GlobalConfig{
			GitProviderConfig: conf.GitProviderConfig{
				OrgLevelWebhook: false,
				OrgName: "group1",
				WebhookURL: "testing-url",
				RepoList: project,
			},
		},
	}

	hook := []gitlab.ProjectHook{{
		ID: 1234,
		URL: c.cfg.GitProviderConfig.WebhookURL,
	},}
	
	hooksUrl := fmt.Sprintf("/api/v4/projects/%s/%s/hooks",c.cfg.GitProviderConfig.OrgName, project)
	mux.HandleFunc(hooksUrl, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		mockHTTPResponse(t, w, hook)
	})
	//
	// Execute
	//
	projectHook, isEnabled := IsProjectWebhookEnabled(&c, "test-repo1")
	//
	// Assert
	//
	assert := assertion.New(t)
	assert.Equal(isEnabled, true)
	assert.Equal(projectHook.URL, c.cfg.GitProviderConfig.WebhookURL)
}

