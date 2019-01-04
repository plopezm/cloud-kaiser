package logger

import (
	"context"
	"github.com/plopezm/cloud-kaiser/kaiser-service/contextvars"
	"github.com/plopezm/cloud-kaiser/kaiser-service/engine"
	"log"
	"os"
)

func init() {
	engine.RegisterPlugin(new(LogPlugin))
}

// LogPlugin is used to save process context
type LogPlugin struct {
	context context.Context
}

// GetInstance Creates a new plugin instance with a context
func (plugin *LogPlugin) SetContext(context context.Context) {
	plugin.context = context
}

// GetFunctions returns the functions to be registered in the VM
func (plugin *LogPlugin) GetFunctions() map[string]interface{} {
	functions := make(map[string]interface{})
	functions["Logger"] = map[string]interface{}{
		"info": plugin.Info,
	}
	return functions
}

// Info Prints objects or strings sent as parameters
func (plugin *LogPlugin) Info(args ...interface{}) {
	jobName := plugin.context.Value(contextvars.JobName).(string)
	jobVersion := plugin.context.Value(contextvars.JobVersion).(string)
	f, err := os.OpenFile("logs/"+jobName+"_"+jobVersion+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Error creating file logs/" + jobName + ".log")
		return
	}
	defer f.Close()
	logger := log.New(f, "", log.Ldate|log.Ltime|log.LUTC)
	for _, arg := range args {
		logger.Println(arg)
	}
}
