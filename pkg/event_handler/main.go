package event_handler

import (
	"context"
	"github.com/quickube/piper/pkg/clients"
	"github.com/quickube/piper/pkg/conf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func Start(ctx context.Context, stop context.CancelFunc, cfg *conf.GlobalConfig, clients *clients.Clients) {
	labelSelector := &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{Key: "piper.quickube.com/notified",
				Operator: metav1.LabelSelectorOpExists},
		},
	}
	watcher, err := clients.Workflows.Watch(&ctx, labelSelector)
	if err != nil {
		log.Printf("[event handler] Failed to watch workflow error:%s", err)
		return
	}

	notifier := NewGithubEventNotifier(cfg, clients)
	handler := &workflowEventHandler{
		Clients:  clients,
		Notifier: notifier,
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Print("[event handler] context canceled, stopping watcher")
				watcher.Stop()
				return
			case event, ok := <-watcher.ResultChan():
				if !ok {
					log.Print("[event handler] result channel closed")
					watcher.Stop()
					stop()
					return
				}
				if err2 := handler.Handle(ctx, &event); err2 != nil {
					log.Printf("[event handler] failed to Handle workflow event: %v", err2)
				}
			}
		}

	}()
}
