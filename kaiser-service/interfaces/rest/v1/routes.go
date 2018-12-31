package v1

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/plopezm/cloud-kaiser/core/db"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/core/util"
	"github.com/plopezm/cloud-kaiser/kaiser-service/engine"
	"net/http"
)

func AddRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/health", healthStatusHandler).Methods("STATUS", "GET")
	router.HandleFunc("/job/{name}/{version}", executeJob).Methods("POST")
	return router
}

func healthStatusHandler(w http.ResponseWriter, r *http.Request) {
	logger.GetLogger().Debug("Called healthStatusHandler")
	util.ResponseOk(w, map[string]interface{}{
		"status": "UP",
	})
}

func executeJob(w http.ResponseWriter, r *http.Request) {
	logger.GetLogger().Debug("Called executeJob")
	vars := mux.Vars(r)

	job, err := db.FindJobByNameAndVersion(r.Context(), vars["name"], vars["version"])
	if err != nil {
		util.ResponseError(w, http.StatusNotFound, fmt.Sprintf("Job %s:%s not found", vars["name"], vars["version"]))
		return
	}
	engine.Execute(engine.CreateRunnable(*job))

	util.ResponseOk(w, map[string]interface{}{
		"status": "OK",
		"job":    fmt.Sprintf("%s:%s", vars["name"], vars["version"]),
		"msg":    "Job executed",
	})
}
