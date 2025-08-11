package main

import (
	"log"

	"forte/internal/activities"
	"forte/internal/config"
	"forte/internal/constants"
	"forte/internal/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config  %v", err)
	}

	c, err := client.Dial(client.Options{HostPort: cfg.TemporalHost})
	if err != nil {
		log.Fatalf("failed to create temporal client %v", err)
	}
	defer c.Close()

	w := worker.New(c, constants.ORDER_TASK_QUEUE, worker.Options{})

	w.RegisterWorkflow(workflows.OrderWorkflow)

	orderActivities := activities.NewOrderActivities()
	for _, activity := range orderActivities.GetActivities() {
		w.RegisterActivity(activity)
	}

	log.Printf("Starting temporal worker on task queue: %s", constants.ORDER_TASK_QUEUE)
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("failed to start worker %v", err)
	}
}
