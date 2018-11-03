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

func NewRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/v1/tasks", listTasksHandler).
		Methods("GET", "OPTIONS")
	router.HandleFunc("/v1/search", searchTaskHandler).
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
