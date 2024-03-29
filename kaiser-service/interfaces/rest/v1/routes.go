package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/core/search"
	"github.com/plopezm/cloud-kaiser/core/util"
	"github.com/plopezm/cloud-kaiser/kaiser-service/engine"
)

func AddRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/health", healthStatusHandler).Methods("STATUS", "GET")
	router.HandleFunc("/jobs/{name}/{version}", executeJob).Methods("POST")
	router.HandleFunc("/jobs/{name}/{version}", getJobLogs).Methods("STATUS", "GET")
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

	jobs, err := search.FindJobs(r.Context(), fmt.Sprintf("%s:%s", vars["name"], vars["version"]), 0, 1)
	if err != nil || len(jobs) == 0 {
		util.ResponseError(w, http.StatusNotFound, fmt.Sprintf("Job %s:%s not found", vars["name"], vars["version"]))
		return
	}
	parameters := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&parameters)
	if err != nil {
		util.ResponseError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding parameters: %s", err.Error()))
		return
	}

	onJobFinish := make(chan engine.Runnable)

	engine.Execute(engine.CreateRunnable(jobs[0]), parameters, onJobFinish)
	runnable := <-onJobFinish

	util.ResponseOk(w, map[string]interface{}{
		"status": runnable.GetResultStatus(),
		"job":    fmt.Sprintf("%s:%s", vars["name"], vars["version"]),
		"msg":    "Job executed",
	})
}

func getJobLogs(w http.ResponseWriter, r *http.Request) {
	logger.GetLogger().Debug("Called getJobLogs")
	vars := mux.Vars(r)

	content, err := engine.GetLogs(vars["name"], vars["version"])
	if err != nil {
		util.ResponseError(w, http.StatusNotFound, fmt.Sprintf("Error found %s", err.Error()))
		return
	}
	util.ResponseOk(w, map[string]interface{}{
		"status": true,
		"job":    fmt.Sprintf("%s:%s", vars["name"], vars["version"]),
		"logs":   strings.Split(content, "\n"),
	})
}
