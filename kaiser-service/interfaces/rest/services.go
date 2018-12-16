package rest

import (
	"github.com/gorilla/mux"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/kaiser-service/interfaces/rest/v1"
)

func GetRouter() *mux.Router {
	router := mux.NewRouter()
	logger.GetLogger().Info("Setting REST routes:")

	// Adding routes to the subrouter version 1
	v1Router := router.PathPrefix("/kaiser/v1").Subrouter()
	v1.AddRoutes(v1Router)
	printRoutes(v1Router)

	return router
}

func printRoutes(r *mux.Router) {
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		logger.GetLogger().Info(t)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
