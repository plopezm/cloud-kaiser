package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/plopezm/cloud-kaiser/core/event"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/core/search"
	"github.com/plopezm/cloud-kaiser/core/types"
	"github.com/plopezm/cloud-kaiser/pusher-service/wsocket"
	"github.com/tinrab/retry"
)

//Config Represents the service configuration
type Config struct {
	NatsAddress          string `envconfig:"NATS_ADDRESS"`
	ElasticSearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS"`
	LogLevel             string `envconfig:"LOG_LEVEL"`
	ServicePort          int    `envconfig:"SERVICE_PORT"`
}

func main() {
	// Parse configuration from environment variables
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	logger.InitializeLogger(config.LogLevel)
	types.RegisterCoreTypes()

	// Connect to ElasticSearch
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := search.NewElasticSearch(fmt.Sprintf("http://%s", config.ElasticSearchAddress))
		if err != nil {
			log.Println(err)
			return err
		}
		logger.GetLogger().Info("ElasticSearch connected!")
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
		err = messaging.OnQueuedEvent("pusher-service", event.TaskCreated, reflect.TypeOf(types.Task{}), onEventReceived)
		if err != nil {
			logger.GetLogger().Error(err)
			return err
		}
		err = messaging.OnQueuedEvent("pusher-service", event.TaskExecutionLog, reflect.TypeOf(types.TaskExecutionLog{}), onEventReceived)
		if err != nil {
			logger.GetLogger().Error(err)
			return err
		}
		err = messaging.OnEvent(event.NotifyUI, reflect.TypeOf(event.UINotification{}), onUIEvent)
		if err != nil {
			logger.GetLogger().Error(err)
			return err
		}
		event.SetEventStore(messaging)
		return nil
	})
	defer event.Close()

	// Run WebSocket server
	go wsocket.Start()
	http.HandleFunc("/ws/", wsocket.WsPage)
	logger.GetLogger().Debug("Starting service: pusher-service on port %d", config.ServicePort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.ServicePort), nil); err != nil {
		logger.GetLogger().Fatal(err)
	}
}

func onEventReceived(packet event.Envelope) {
	logger.GetLogger().Debug(fmt.Sprintf("Event received: %s", packet.Subject))
	switch packet.Subject {
	case event.TaskCreated:
		if err := search.InsertTask(context.Background(), packet.Content.(types.Task)); err != nil {
			logger.GetLogger().Error(fmt.Sprintf("Error in event %s: %s", event.TaskCreated, err.Error()))
		}
		logger.GetLogger().Debug(fmt.Sprintf("Event %s processed", packet.Subject))
	case event.TaskExecutionLog:
		taskExecutionLog := packet.Content.(types.TaskExecutionLog)
		taskExecutionLog.Ts = time.Now().UnixNano()
		err := search.InsertLog(context.Background(), taskExecutionLog)
		if err != nil {
			logger.GetLogger().Error("Error sending message to Elasticsearch: " + err.Error())
		}
		logger.GetLogger().Debug(fmt.Sprintf("Event %s processed", packet.Subject))
	}
}

func onUIEvent(packet event.Envelope) {
	logger.GetLogger().Debug(
		fmt.Sprintf("Event received: %s for type %s", packet.Subject, packet.Content.(event.UINotification).Type))
	wsocket.Broadcast(packet.Content, nil)
}
