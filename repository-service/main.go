package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/plopezm/cloud-kaiser/core/db"
	"github.com/plopezm/cloud-kaiser/core/event"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/core/types"
	v1 "github.com/plopezm/cloud-kaiser/repository-service/v1"
	"github.com/tinrab/retry"
)

//Config The service configuration
type Config struct {
	PostgresAddr     string `envconfig:"POSTGRES_ADDR"`
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
		log.Fatal(err)
	}

	logger.InitializeLogger(config.LogLevel)
	types.RegisterCoreTypes()

	// Connect to PostgreSQL and inject the repository. The code below retries connection every 2 seconds
	retry.ForeverSleep(2*time.Second, func(attempt int) error {
		addr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresAddr, config.PostgresDB)
		repo, err := db.NewPostgres(addr)
		if err != nil {
			log.Println(err)
			return err
		}
		logger.GetLogger().Info("PostgreSQL connected!")
		db.SetRepository(repo)
		return nil
	})
	defer db.Close()

	// Connect to NATS
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		repo, err := event.NewNats(fmt.Sprintf("nats://%s", config.NatsAddress))
		if err != nil {
			log.Println(err)
			return err
		}
		logger.GetLogger().Info("NATS connected!")
		event.SetEventStore(repo)
		return nil
	})
	defer event.Close()

	logger.GetLogger().Info("Starting service: pusher-service on port %d", config.ServicePort)
	router := v1.NewRouter()
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.ServicePort), router); err != nil {
		log.Fatal(err)
	}
}
