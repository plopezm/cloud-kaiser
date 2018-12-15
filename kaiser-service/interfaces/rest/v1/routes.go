package v1

import (
	"github.com/gorilla/mux"
	"github.com/plopezm/cloud-kaiser/core/util"
	"net/http"
)

func AddRoutes(router *mux.Router) (*mux.Router) {
	router.HandleFunc("/health", createTaskHandler).Methods("STATUS", "GET")
	return router
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	util.ResponseOk(w, map[string]interface{}{
		"status": "UP",
	})
}
