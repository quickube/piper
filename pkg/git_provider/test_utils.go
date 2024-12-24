package git_provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v52/github"
	"github.com/ktrysmt/go-bitbucket"
	"github.com/xanzy/go-gitlab"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints. See issue #752.
	baseURLPath          = "/api-v3"
	bitbucketBaseURLPath = "/2.0"
)

func setup() (client *github.Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		fmt.Fprintln(os.Stderr, "\tSee https://github.com/google/go-github/issues/752 for information.")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the GitHub client being tested and is
	// configured to use test server.
	client = github.NewClient(nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url
	client.UploadURL = url

	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	t.Helper()
	want := url.Values{}
	for k, v := range values {
		want.Set(k, v)
	}

	err := r.ParseForm()
	if err != nil {
		t.Errorf("Go error parsing form: %v", err)
	}
	if got := r.Form; !cmp.Equal(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func setupBitbucket() (client *bitbucket.Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(bitbucketBaseURLPath+"/", http.StripPrefix(bitbucketBaseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server := httptest.NewServer(apiHandler)
	url, _ := url.Parse(server.URL + bitbucketBaseURLPath)
	client = bitbucket.NewBasicAuth("username", "password")
	client.SetApiBaseURL(*url)

	return client, mux, server.URL, server.Close
}
func setupGitlab(t *testing.T) (*http.ServeMux, *gitlab.Client) {
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	// client is the Gitlab client being tested.
	client, err := gitlab.NewClient("",
	gitlab.WithBaseURL(server.URL),
		// Disable backoff to speed up tests that expect errors.
		gitlab.WithCustomBackoff(func(_, _ time.Duration, _ int, _ *http.Response) time.Duration {
			return 0
		}),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return mux, client
}