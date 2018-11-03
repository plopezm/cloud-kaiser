package v1

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/plopezm/cloud-kaiser/core/db"
	"github.com/plopezm/cloud-kaiser/core/event"
	"github.com/plopezm/cloud-kaiser/core/types"
	"github.com/plopezm/cloud-kaiser/core/util"
	"log"
	"net/http"
)

func NewRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/v1/tasks", createTaskHandler).
		Methods("POST")
	return
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var task types.JobTask
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		util.ResponseError(w, http.StatusBadRequest, err.Error())
	}

	if err := db.InsertTask(ctx, task); err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	// Publish event
	if err := event.PublishTaskCreated(task); err != nil {
		log.Println(err)
	}

	// Return new meow
	util.ResponseOk(w, task)
}
