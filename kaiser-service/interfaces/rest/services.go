package rest

import (
	"github.com/gorilla/mux"
	"github.com/plopezm/cloud-kaiser/kaiser-service/interfaces/rest/v1"
)

func GetRouter() *mux.Router {
	router := mux.NewRouter()

	// Adding routes to the subrouter
	v1Router := router.PathPrefix("/kaiser/v1").Subrouter()
	v1.AddRoutes(v1Router)

	return router
}
