package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/plopezm/cloud-kaiser/core/db"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/kaiser-service/engine"
	"github.com/plopezm/cloud-kaiser/kaiser-service/interfaces"
	_ "github.com/plopezm/cloud-kaiser/kaiser-service/plugins/http"
	_ "github.com/plopezm/cloud-kaiser/kaiser-service/plugins/logger"
	_ "github.com/plopezm/cloud-kaiser/kaiser-service/plugins/system"
	"github.com/tinrab/retry"
	"os"
	"time"
)

type Config struct {
	PostgresAddress  string `envconfig:"POSTGRES_ADDR"`
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
	log.Debug("Starting service...")

	// Connect to PostgreSQL and inject the repository. The code below retries connection every 2 seconds
	retry.ForeverSleep(2*time.Second, func(attempt int) error {
		addr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresAddress, config.PostgresDB)
		repo, err := db.NewPostgres(addr)
		if err != nil {
			log.Warn(err)
			return err
		}
		db.SetRepository(repo)
		return nil
	})
	defer db.Close()

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
