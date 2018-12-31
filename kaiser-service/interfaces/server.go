package interfaces

import (
	"fmt"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/kaiser-service/interfaces/rest"
	"net/http"
)

// StartServer Starts graphql api on selected port
func StartServer(port int) {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), rest.GetRouter()); err != nil {
		logger.GetLogger().Fatal(err)
	}
}
