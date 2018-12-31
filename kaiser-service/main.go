package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/plopezm/cloud-kaiser/core/db"
	"github.com/plopezm/cloud-kaiser/core/event"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/kaiser-service/engine"
	"github.com/plopezm/cloud-kaiser/kaiser-service/interfaces"
	"github.com/tinrab/retry"
	"time"
)

type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
	LogLevel         string `envconfig:"LOG_LEVEL"`
	ServicePort      int    `envconfig:"SERVICE_PORT"`
}

func main() {
	// Parse configuration from environment variables
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		panic(err)
	}

	logger.InitializeLogger(config.LogLevel)
	log := logger.GetLogger()

	// Connect to PostgreSQL and inject the repository. The code below retries connection every 2 seconds
	retry.ForeverSleep(2*time.Second, func(attempt int) error {
		addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresDB)
		repo, err := db.NewPostgres(addr)
		if err != nil {
			log.Warn(err)
			return err
		}
		db.SetRepository(repo)
		return nil
	})
	defer db.Close()

	// Connect to NATS
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		repo, err := event.NewNats(fmt.Sprintf("nats://%s", config.NatsAddress))
		if err != nil {
			log.Warn(err)
			return err
		}
		event.SetEventStore(repo)
		return nil
	})
	defer event.Close()

	servicePort := config.ServicePort
	if servicePort == 0 {
		servicePort = 8080
	}
	log.Debug("Starting engine")
	go engine.Start()
	log.Debug(fmt.Sprintf("Starting server at port %d", servicePort))
	interfaces.StartServer(servicePort)
}
