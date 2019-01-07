package system

import (
	"context"

	"github.com/plopezm/cloud-kaiser/core/db"
	"github.com/plopezm/cloud-kaiser/kaiser-service/engine"
	"github.com/robertkrimen/otto"
)

func init() {
	engine.RegisterPlugin(new(OSPlugin))
}

// OSPlugin is used to save process context
type OSPlugin struct {
}

// GetFunctions returns the functions to be registered in the VM
func (plugin *OSPlugin) GetFunctions() map[string]interface{} {
	functions := make(map[string]interface{})
	functions["Process"] = map[string]interface{}{
		"sleep": Sleep,
	}
	functions["System"] = map[string]interface{}{
		"call": Call,
	}
	return functions
}

// Call Calls an existing job
func Call(jobName string, version string, params map[string]interface{}) otto.Value {
	job, err := db.FindJobByNameAndVersion(context.Background(), jobName, version)
	if err != nil {
		res, _ := otto.ToValue(err.Error())
		return res
	}
	runnable := engine.CreateRunnable(*job)
	engine.Execute(runnable, params, nil)
	return otto.Value{}
}
