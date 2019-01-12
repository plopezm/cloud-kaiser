package main

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/plopezm/cloud-kaiser/core/event"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/core/search"
	"github.com/plopezm/cloud-kaiser/core/types"
	"github.com/plopezm/cloud-kaiser/pusher-service/pusher"
	"github.com/tinrab/retry"
	"log"
	"net/http"
	"time"
)

type Config struct {
	NatsAddress          string `envconfig:"NATS_ADDRESS"`
	ElasticSearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS"`
	LogLevel             string `envconfig:"LOG_LEVEL"`
	ServicePort          int    `envconfig:"SERVICE_PORT"`
}

var hub *pusher.Hub

func main() {
	// Parse configuration from environment variables
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	logger.InitializeLogger(config.LogLevel)
	hub = pusher.NewHub()

	// Connect to ElasticSearch
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := search.NewElasticSearch(fmt.Sprintf("http://%s", config.ElasticSearchAddress))
		if err != nil {
			log.Println(err)
			return err
		}
		search.SetRepository(es)
		return nil
	})
	defer search.Close()

	// Connect to NATS
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		messaging, err := event.NewNats(fmt.Sprintf("nats://%s", config.NatsAddress))
		if err != nil {
			log.Println(err)
			return err
		}
		err = messaging.OnEvent(event.TaskCreated, onTaskCreated)
		if err != nil {
			log.Println(err)
			return err
		}
		event.SetEventStore(messaging)
		return nil
	})
	defer event.Close()

	// Run WebSocket server
	go hub.Run()
	http.HandleFunc("/pusher", hub.HandleWebSocket)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.ServicePort), nil); err != nil {
		logger.GetLogger().Fatal(err)
	}
}

func onTaskCreated(packet event.Envelope) {
	switch packet.Destination {
	case event.TaskCreated:
		if err := search.InsertTask(context.Background(), packet.Content.(types.Task)); err != nil {
			logger.GetLogger().Error(fmt.Sprintf("Error in event %s: %s", event.TaskCreated, err.Error()))
		}
	}
	hub.Broadcast(packet)
}
