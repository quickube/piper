package event_handler

import (
	"context"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/quickube/piper/pkg/clients"
	"github.com/quickube/piper/pkg/conf"
	"github.com/quickube/piper/pkg/utils"
)

type eventNotifier struct {
	cfg     *conf.GlobalConfig
	clients *clients.Clients
}

func NewEventNotifier(cfg *conf.GlobalConfig, clients *clients.Clients) EventNotifier {
	return &eventNotifier{
		cfg:     cfg,
		clients: clients,
	}
}

func (en *eventNotifier) Notify(ctx context.Context, workflow *v1alpha1.Workflow) error {
	fmt.Printf("Notifing workflow, %s\n", workflow.GetName())

	repo, ok := workflow.GetLabels()["repo"]
	if !ok {
		return fmt.Errorf("failed get repo label for workflow: %s", workflow.GetName())
	}
	commit, ok := workflow.GetLabels()["commit"]
	if !ok {
		return fmt.Errorf("failed get commit label for workflow: %s", workflow.GetName())
	}

	workflowLink := fmt.Sprintf("%s/workflows/%s/%s", en.cfg.WorkflowServerConfig.ArgoAddress, en.cfg.Namespace, workflow.GetName())

	status, err := en.clients.GitProvider.GetCorrelatingEvent(ctx, &workflow.Status.Phase)
	if err != nil {
		return fmt.Errorf("failed to translate workflow status for phase: %s status: %s", string(workflow.Status.Phase), status)
	}

	message := utils.TrimString(workflow.Status.Message, 140) // Max length of message is 140 characters
	err = en.clients.GitProvider.SetStatus(ctx, &repo, &commit, &workflowLink, &status, &message)
	if err != nil {
		return fmt.Errorf("failed to set status for workflow %s: %s", workflow.GetName(), err)
	}

	return nil
}
