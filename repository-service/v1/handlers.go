package v1

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/plopezm/cloud-kaiser/core/db"
	"github.com/plopezm/cloud-kaiser/core/event"
	"github.com/plopezm/cloud-kaiser/core/types"
	"github.com/plopezm/cloud-kaiser/core/util"
	"golang.org/x/net/context"
)

const handlersPrefix = "/v1"

//NewRouter Creates a new router instance with the defined routes
func NewRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc(handlersPrefix+"/tasks", createTaskHandler).
		Methods("POST")
	router.HandleFunc(handlersPrefix+"/jobs", createJobHandler).
		Methods("POST")
	return
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var task types.Task
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		util.ResponseError(w, http.StatusBadRequest, err.Error())
	}

	if err := db.Tx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted}, func(tx *sql.Tx) error {
		ctxWithTX := context.WithValue(ctx, db.ContextTX, tx)

		if err := db.InsertTask(ctxWithTX, task); err != nil {
			log.Println(err)
			return err
		}

		// Publish event
		if err := event.PublishEvent(event.TaskCreated, task); err != nil {
			log.Println(err)
			return err
		}
		return nil
	}); err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusBadRequest, "Create task transaction error: "+err.Error())
		return
	}

	if err := event.PublishEvent(event.NotifyUI, event.UINotification{
		Type:    event.UITaskCreated,
		Content: task,
	}); err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusBadRequest, "Create task error: "+err.Error())
		return
	}

	// Return new meow
	util.ResponseOk(w, task)
}

func createJobHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var job types.Job
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&job)
	if err != nil {
		util.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.InsertJob(ctx, &job); err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusBadRequest, "Create job error: "+err.Error())
		return
	}

	// Return new meow
	util.ResponseOk(w, job)
}
