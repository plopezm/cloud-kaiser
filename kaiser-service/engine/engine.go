package engine

import (
	"fmt"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"github.com/plopezm/cloud-kaiser/kaiser-service/contextvars"
	"io/ioutil"
	"sync"
	"time"
)

var threadFinishedChannel = make(chan Runnable)
var mutex = &sync.Mutex{}
var executions = make(map[string]interface{})

func Start() {
	for runnable := range threadFinishedChannel {
		mutex.Lock()
		processKey := fmt.Sprintf("%s-%d", runnable.GetIdentifier(), runnable.GetStartTime())
		delete(executions, processKey)
		mutex.Unlock()
		logger.GetLogger().Debug(fmt.Sprintf("Execution of process %s finished", processKey))
	}
}

func Execute(runnable Runnable, parameters map[string]interface{}, onProcessFinishes chan Runnable) {
	mutex.Lock()
	defer mutex.Unlock()
	thread := func() {
		runnable.Run()
		threadFinishedChannel <- runnable
		if onProcessFinishes != nil {
			onProcessFinishes <- runnable
		}
	}
	runnable.SetStartTime(time.Now().UnixNano())
	runnable.SetStatus(RunnableStatusRunning)
	runnable.SetParameters(parameters)
	processKey := fmt.Sprintf("%s-%d", runnable.GetIdentifier(), runnable.GetStartTime())
	executions[processKey] = runnable
	go thread()
	logger.GetLogger().Debug(fmt.Sprintf("Execution of process %s started", processKey))
}

func GetExecutions() []Runnable {
	mutex.Lock()
	defer mutex.Unlock()
	executions := make([]Runnable, 0)
	for _, runnable := range executions {
		executions = append(executions, runnable)
	}
	return executions
}

func GetLogs(jobname string, version string) (string, error) {
	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s_%s.log", contextvars.DefaultLogFolder, jobname, version))
	return string(bytes), err
}

func notifyExecutionLogs(jobname string, jobversion string, taskname string, taskversion string) {
	//bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s_%s.log", contextvars.DefaultLogFolder, jobname, jobversion))

	//search.InsertLog(context.Background(), jobname, jobversion, taskname, taskversion, )
}
