package engine

import (
	"fmt"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/core/types"
	"github.com/plopezm/cloud-kaiser/kaiser-service/contextvars"
	"golang.org/x/net/context"
	"strings"
)

func CreateRunnable(job types.Job) Runnable {
	return &JobRunner{
		Identifier: fmt.Sprintf("%s:%s", job.Name, job.Version),
		Job:        job,
	}
}

type JobRunner struct {
	types.Job
	Identifier   string
	StartTime    int64
	ResultStatus bool
}

func (r *JobRunner) GetIdentifier() string {
	return r.Identifier
}

func (r *JobRunner) SetStartTime(ts int64) {
	r.StartTime = ts
}

func (r *JobRunner) GetStartTime() int64 {
	return r.StartTime
}

func (r *JobRunner) GetStatus() RunnableStatus {
	return RunnableStatus(r.Status)
}

func (r *JobRunner) SetStatus(status RunnableStatus) {
	r.Status = types.JobStatus(status)
}

func (r *JobRunner) Run() {
	runnerContext := context.WithValue(context.Background(), contextvars.JobName, r.Name)
	runnerContext = context.WithValue(runnerContext, contextvars.JobVersion, r.Version)
	vm := NewVM(runnerContext)
	r.SetStatus(RunnableStatusRunning)
	task, found := r.Tasks[r.Entrypoint]
	logger.GetLogger().Debug(fmt.Sprintf("Job %s:%s, entrypoint %s started", r.Name, r.Version, task.Name))
	var result bool
	for found {
		logger.GetLogger().Debug(fmt.Sprintf("Executing script: \n=====\n%s\n=====", task.Script))
		_, err := vm.Run(task.Script)
		if err == nil {
			task, found = r.Tasks[getTaskName(task.OnSuccess)]
			result = true
		} else {
			logger.GetLogger().Debug(fmt.Sprintf("Task %s exited with error: %s", task.Name, err.Error()))
			task, found = r.Tasks[getTaskName(task.OnFailure)]
			result = false
		}
	}
	r.SetStatus(RunnableStatusStopped)
	r.SetResultStatus(result)
}

func getTaskName(taskid string) string {
	return strings.Split(taskid, ":")[0]
}

func (r *JobRunner) SetParameters(receivedParameters map[string]interface{}) {
	allParams := make([]types.JobArgs, 0)
	for _, parameter := range r.Parameters {
		if value, ok := receivedParameters[parameter.Name]; ok {
			parameter.Value = value
		}
		allParams = append(allParams, parameter)
	}
	r.Parameters = allParams
}

func (r *JobRunner) SetResultStatus(result bool) {
	r.ResultStatus = result
}

func (r *JobRunner) GetResultStatus() bool {
	return r.ResultStatus
}
