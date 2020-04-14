package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/plopezm/cloud-kaiser/core/db"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/core/search"
	"github.com/plopezm/cloud-kaiser/core/types"
	v1 "github.com/plopezm/cloud-kaiser/query-service/v1"
	"github.com/tinrab/retry"
)

//Config The configuration used by this service
type Config struct {
	PostgresDB           string `envconfig:"POSTGRES_DB"`
	PostgresUser         string `envconfig:"POSTGRES_USER"`
	PostgresPassword     string `envconfig:"POSTGRES_PASSWORD"`
	ElasticSearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS"`
	LogLevel             string `envconfig:"LOG_LEVEL"`
	PostgresAddr         string `envconfig:"POSTGRES_ADDR"`
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
		db.SetRepository(repo)
		return nil
	})
	defer db.Close()

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

	router := v1.NewRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
