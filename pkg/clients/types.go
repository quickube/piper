package clients

import (
	"github.com/quickube/piper/pkg/git_provider"
	"github.com/quickube/piper/pkg/workflow_handler"
)

type Clients struct {
	GitProvider git_provider.Client
	Workflows   workflow_handler.WorkflowsClient
}
