package v1

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/plopezm/cloud-kaiser/core/db"
	"github.com/plopezm/cloud-kaiser/core/search"
	"github.com/plopezm/cloud-kaiser/core/types"
	"github.com/plopezm/cloud-kaiser/core/util"
)

const handlersPrefix = "/v1"

//NewRouter Creates a new mux router with paths
func NewRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc(handlersPrefix+"/search/tasks", searchTasks).
		Methods("GET", "OPTIONS")
	router.HandleFunc(handlersPrefix+"/search/jobs", searchJobs).
		Methods("GET", "OPTIONS")
	router.HandleFunc(handlersPrefix+"/search/logs", searchLogs).
		Methods("GET", "OPTIONS")
	router.HandleFunc(handlersPrefix+"/tasks", listTasksHandler).
		Methods("GET", "OPTIONS")
	router.HandleFunc(handlersPrefix+"/tasks/{name}", getTaskByName).
		Methods("GET", "OPTIONS")
	router.HandleFunc(handlersPrefix+"/tasks/{name}/version/{version}", getTaskByNameAndVersion).
		Methods("GET", "OPTIONS")
	router.HandleFunc(handlersPrefix+"/jobs", listJobsHandler).
		Methods("GET", "OPTIONS")
	router.HandleFunc(handlersPrefix+"/jobs/{name}", getJobByName).
		Methods("GET", "OPTIONS")
	router.HandleFunc(handlersPrefix+"/jobs/{name}/version/{version}", getJobByNameAndVersion).
		Methods("GET", "OPTIONS")
	return
}

func searchLogs(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	// Read parameters
	query := r.FormValue("query")
	if len(query) == 0 {
		util.ResponseError(w, http.StatusBadRequest, "Missing 'query' parameter")
		return
	}

	offset := uint64(0)
	offsetStr := r.FormValue("offset")
	limit := uint64(100)
	takeStr := r.FormValue("limit")
	if len(offsetStr) != 0 {
		offset, err = strconv.ParseUint(offsetStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid 'offset' parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		limit, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid 'limit' parameter")
			return
		}
	}

	// Search logs
	logs, err := search.FindLogs(ctx, query, offset, limit)
	if err != nil {
		log.Println(err)
		util.ResponseOk(w, []types.JobTask{})
		return
	}

	util.ResponseOk(w, logs)
}

func searchTasks(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	// Read parameters
	query := r.FormValue("query")
	if len(query) == 0 {
		util.ResponseError(w, http.StatusBadRequest, "Missing 'query' parameter")
		return
	}
	offset := uint64(0)
	offsetStr := r.FormValue("offset")
	limit := uint64(100)
	takeStr := r.FormValue("limit")
	if len(offsetStr) != 0 {
		offset, err = strconv.ParseUint(offsetStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid 'offset' parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		limit, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid 'limit' parameter")
			return
		}
	}

	// Search tasks
	tasks, err := search.FindTasks(ctx, query, offset, limit)
	if err != nil {
		log.Println(err)
		util.ResponseOk(w, []types.JobTask{})
		return
	}

	util.ResponseOk(w, tasks)
}

func searchJobs(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	// Read parameters
	query := r.FormValue("query")
	if len(query) == 0 {
		util.ResponseError(w, http.StatusBadRequest, "Missing 'query' parameter")
		return
	}
	offset := uint64(0)
	offsetStr := r.FormValue("offset")
	limit := uint64(100)
	takeStr := r.FormValue("limit")
	if len(offsetStr) != 0 {
		offset, err = strconv.ParseUint(offsetStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid 'offset' parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		limit, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid 'limit' parameter")
			return
		}
	}

	// Search tasks
	tasks, err := search.FindJobs(ctx, query, offset, limit)
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

	// Fetch tasks
	tasks, err := db.ListTasks(ctx, offset, limit)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch tasks")
		return
	}

	util.ResponseOk(w, tasks)
}

func getTaskByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	vars := mux.Vars(r)
	taskName := vars["name"]

	// Fetch task
	task, err := db.FindTaskByName(ctx, taskName)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch tasks")
		return
	}

	util.ResponseOk(w, task)
}

func getTaskByNameAndVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	vars := mux.Vars(r)
	taskName := vars["name"]
	version := vars["version"]

	// Fetch task
	task, err := db.FindTaskByNameAndVersion(ctx, taskName, version)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch tasks")
		return
	}

	util.ResponseOk(w, task)
}

func listJobsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Fetch jobs
	jobs, err := db.ListJobs(ctx, offset, limit)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch tasks")
		return
	}

	util.ResponseOk(w, jobs)
}

func getJobByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	vars := mux.Vars(r)
	taskName := vars["name"]

	// Fetch job
	job, err := db.FindJobByName(ctx, taskName)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch tasks")
		return
	}

	util.ResponseOk(w, job)
}

func getJobByNameAndVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	vars := mux.Vars(r)
	taskName := vars["name"]
	version := vars["version"]

	// Fetch job
	job, err := db.FindJobByNameAndVersion(ctx, taskName, version)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch tasks")
		return
	}

	util.ResponseOk(w, job)
}
