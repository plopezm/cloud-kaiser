package system

import (
	"context"
	"github.com/plopezm/cloud-kaiser/core/db"
	"github.com/plopezm/cloud-kaiser/kaiser-service/engine"
	"github.com/plopezm/cloud-kaiser/kaiser-service/types"
	"github.com/robertkrimen/otto"
)

func init() {
	engine.RegisterPlugin(new(OSPlugin))
}

// OSPlugin is used to save process context
type OSPlugin struct {
	context context.Context
}

// GetFunctions returns the functions to be registered in the VM
func (plugin *OSPlugin) GetFunctions() map[string]interface{} {
	functions := make(map[string]interface{})
	functions["process"] = map[string]interface{}{
		"sleep": plugin.Sleep,
	}
	functions["system"] = map[string]interface{}{
		"call": plugin.Call,
	}
	return functions
}

// GetInstance Creates a new plugin instance with a context
func (plugin *OSPlugin) GetInstance(context context.Context) types.Plugin {
	newPluginInstance := new(OSPlugin)
	newPluginInstance.context = context
	return newPluginInstance
}

// Call Calls an existing job
func (plugin *OSPlugin) Call(jobName string, version string, params map[string]interface{}) otto.Value {
	job, err := db.FindJobByNameAndVersion(plugin.context, jobName, version)
	if err != nil {
		res, _ := otto.ToValue(err.Error())
		return res
	}
	runnable := types.CreateRunnable(*job)
	engine.Execute(runnable)
	return otto.Value{}
}
