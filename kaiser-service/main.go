package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/plopezm/cloud-kaiser/core/event"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/core/search"
	"github.com/plopezm/cloud-kaiser/core/types"
	"github.com/plopezm/cloud-kaiser/kaiser-service/engine"
	"github.com/plopezm/cloud-kaiser/kaiser-service/interfaces"
	_ "github.com/plopezm/cloud-kaiser/kaiser-service/plugins/http"
	_ "github.com/plopezm/cloud-kaiser/kaiser-service/plugins/logger"
	_ "github.com/plopezm/cloud-kaiser/kaiser-service/plugins/system"
	"github.com/tinrab/retry"
)

//Config Contains the service configuration
type Config struct {
	ElasticSearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS"`
	NatsAddress          string `envconfig:"NATS_ADDRESS"`
	LogLevel             string `envconfig:"LOG_LEVEL"`
	ServicePort          int    `envconfig:"SERVICE_PORT"`
}

func main() {
	// Parse configuration from environment variables
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		panic(err)
	}

	logger.InitializeLogger(config.LogLevel)
	types.RegisterCoreTypes()

	log := logger.GetLogger()
	log.Debug("Starting service...")

	// Connect to NATS
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		messaging, err := event.NewNats(fmt.Sprintf("nats://%s", config.NatsAddress))
		if err != nil {
			logger.GetLogger().Error(err)
			return err
		}
		event.SetEventStore(messaging)
		return nil
	})
	defer event.Close()

	// Connect to ElasticSearch
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := search.NewElasticSearch(fmt.Sprintf("http://%s", config.ElasticSearchAddress))
		if err != nil {
			logger.GetLogger().Error(err)
			return err
		}
		logger.GetLogger().Info("ElasticSearch connected!")
		search.SetRepository(es)
		return nil
	})
	defer search.Close()

	servicePort := config.ServicePort
	if servicePort == 0 {
		servicePort = 8080
	}
	os.Mkdir("logs", 0755)
	log.Debug("Starting engine")
	go engine.Start()
	log.Debug(fmt.Sprintf("Starting server at port %d", servicePort))
	interfaces.StartServer(servicePort)
}
