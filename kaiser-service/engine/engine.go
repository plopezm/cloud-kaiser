package engine

import (
	"fmt"
	"github.com/plopezm/cloud-kaiser/core/logger"
	"sync"
	"time"
)

var threadFinishedChannel = make(chan Runnable)
var mutex = &sync.Mutex{}
var executions = make(map[string]interface{})

func Start() {
	for runnable := range threadFinishedChannel {
		mutex.Lock()
		delete(executions, fmt.Sprintf("%s-%d", runnable.GetIdentifier(), runnable.GetStartTime()))
		mutex.Unlock()
		logger.GetLogger().Debug(fmt.Sprintf("Execution of process %s-%d finished", runnable.GetIdentifier(), runnable.GetStartTime()))
	}
}

func Execute(runnable Runnable, parameters map[string]interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	thread := func() {
		runnable.Run()
		threadFinishedChannel <- runnable
	}
	runnable.SetStartTime(time.Now().UnixNano())
	runnable.SetStatus(RunnableStatusRunning)
	runnable.SetParameters(parameters)
	executions[fmt.Sprintf("%s-%d", runnable.GetIdentifier(), runnable.GetStartTime())] = runnable
	go thread()
	logger.GetLogger().Debug(fmt.Sprintf("Execution of process %s-%d started", runnable.GetIdentifier(), runnable.GetStartTime()))
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
