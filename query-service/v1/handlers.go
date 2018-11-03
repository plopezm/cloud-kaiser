package v1

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/plopezm/cloud-kaiser/core/db"
	"github.com/plopezm/cloud-kaiser/core/search"
	"github.com/plopezm/cloud-kaiser/core/types"
	"github.com/plopezm/cloud-kaiser/core/util"
	"log"
	"net/http"
	"strconv"
)

const handlersPrefix = "/v1"

func NewRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc(handlersPrefix+"/tasks", listTasksHandler).
		Methods("GET", "OPTIONS")
	router.HandleFunc(handlersPrefix+"/tasks/{name}", getTaskByName).
		Methods("GET", "OPTIONS")
	router.HandleFunc(handlersPrefix+"/tasks/{name}/version/{version}", getTaskByNameAndVersion).
		Methods("GET", "OPTIONS")
	router.HandleFunc(handlersPrefix+"/search", searchTaskHandler).
		Methods("GET", "OPTIONS")
	return
}

func OnTaskCreated(task types.JobTask) {
	if err := search.InsertTask(context.Background(), task); err != nil {
		log.Println(err)
	}
}

func searchTaskHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	// Read parameters
	query := r.FormValue("query")
	if len(query) == 0 {
		util.ResponseError(w, http.StatusBadRequest, "Missing query parameter")
		return
	}
	offset := uint64(0)
	offsetStr := r.FormValue("offset")
	limit := uint64(100)
	takeStr := r.FormValue("limit")
	if len(offsetStr) != 0 {
		offset, err = strconv.ParseUint(offsetStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid offset parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		limit, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
	}

	// Search meows
	tasks, err := search.FindTasks(ctx, query, offset, limit)
	if err != nil {
		log.Println(err)
		util.ResponseOk(w, []types.JobTask{})
		return
	}

	util.ResponseOk(w, tasks)
}

func listTasksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	// Read parameters
	offset := uint64(0)
	offsetStr := r.FormValue("offset")
	limit := uint64(100)
	takeStr := r.FormValue("limit")
	if len(offsetStr) != 0 {
		offset, err = strconv.ParseUint(offsetStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid offset parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		limit, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
	}

	// Fetch meows
	meows, err := db.ListTasks(ctx, offset, limit)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch tasks")
		return
	}

	util.ResponseOk(w, meows)
}

func getTaskByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	vars := mux.Vars(r)
	taskName := vars["name"]

	// Fetch meows
	meows, err := db.FindTaskByName(ctx, taskName)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch tasks")
		return
	}

	util.ResponseOk(w, meows)
}

func getTaskByNameAndVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	vars := mux.Vars(r)
	taskName := vars["name"]
	version := vars["version"]

	// Fetch meows
	meows, err := db.FindTaskByNameAndVersion(ctx, taskName, version)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch tasks")
		return
	}

	util.ResponseOk(w, meows)
}
