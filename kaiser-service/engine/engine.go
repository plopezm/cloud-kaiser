package engine

import (
	"fmt"
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
	}
}

func Execute(runnable Runnable) {
	mutex.Lock()
	defer mutex.Unlock()
	thread := func() {
		runnable.Run()
		threadFinishedChannel <- runnable
	}
	runnable.SetStartTime(time.Now().UnixNano())
	runnable.SetStatus(RunnableStatusRunning)
	executions[fmt.Sprintf("%s-%d", runnable.GetIdentifier(), runnable.GetStartTime())] = runnable
	go thread()
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
