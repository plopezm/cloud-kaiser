package interfaces

import (
	"fmt"
	"github.com/plopezm/cloud-kaiser/kaiser-service/interfaces/rest"
	"log"
	"net/http"
)

// StartServer Starts graphql api on selected port
func StartServer(port int) {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), rest.GetRouter()); err != nil {
		log.Fatal(err)
	}
}