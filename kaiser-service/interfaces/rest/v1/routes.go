package v1

import (
	"github.com/gorilla/mux"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/core/util"
	"net/http"
)

func AddRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/health", healthStatusHandler).Methods("STATUS", "GET")
	return router
}

func healthStatusHandler(w http.ResponseWriter, r *http.Request) {
	logger.GetLogger().Debug("Called healthStatusHandler")
	util.ResponseOk(w, map[string]interface{}{
		"status": "UP",
	})
}
