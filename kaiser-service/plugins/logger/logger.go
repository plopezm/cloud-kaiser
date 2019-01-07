package logger

import (
	"context"
	"fmt"
	"github.com/plopezm/cloud-kaiser/core/search"
	"github.com/plopezm/cloud-kaiser/kaiser-service/contextvars"
	"github.com/plopezm/cloud-kaiser/kaiser-service/engine"
	"github.com/robertkrimen/otto"
	"log"
	"os"
)

func init() {
	engine.RegisterPlugin(new(LogPlugin))
}

// LogPlugin is used to save process context
type LogPlugin struct {
}

// GetFunctions returns the functions to be registered in the VM
func (plugin *LogPlugin) GetFunctions() map[string]interface{} {
	functions := make(map[string]interface{})
	functions["Logger"] = map[string]interface{}{
		"info": Info,
	}
	return functions
}

// Info Prints objects or strings sent as parameters
func Info(call otto.FunctionCall) otto.Value {
	jobName, _ := call.Otto.Get(contextvars.JobName)
	jobVersion, _ := call.Otto.Get(contextvars.JobVersion)
	taskName, _ := call.Otto.Get(contextvars.TaskName)
	taskVersion, _ := call.Otto.Get(contextvars.TaskVersion)
	f, err := os.OpenFile(contextvars.DefaultLogFolder+"/"+jobName.String()+"_"+jobVersion.String()+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Error creating file " + contextvars.DefaultLogFolder + "/" + jobName.String() + ".log")
		return otto.Value{}
	}
	defer f.Close()
	logger := log.New(f, "", log.Ldate|log.Ltime|log.LUTC)
	for _, arg := range call.ArgumentList {
		logline := fmt.Sprintf("[%s:%s] %s", taskName, taskVersion, arg.String())
		logger.Println(logline)
		if search.IsConfigured() {
			err := search.InsertLog(context.Background(), jobName.String(), jobVersion.String(), taskName.String(), taskVersion.String(), logline)
			if err != nil {
				logger.Println("Error sending message to Elasticsearch: " + err.Error())
			}
		}
	}
	return otto.Value{}
}
