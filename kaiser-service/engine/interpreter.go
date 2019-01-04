package engine

import (
	"context"
	"github.com/plopezm/cloud-kaiser/kaiser-service/types"
	"github.com/robertkrimen/otto"
)

var configuredPlugins = make([]types.Plugin, 0)

// RegisterPlugin Registers a plugin
func RegisterPlugin(plugin types.Plugin) {
	configuredPlugins = append(configuredPlugins, plugin)
}

// NewVMWithPlugins Creates a new VM instance using plugins.
// @Param context map[string]interface{} Contains information about the process who creates this VM
func NewVM(context context.Context) *otto.Otto {
	vm := otto.New()
	addRegistedPlugins(vm, context)
	return vm
}

func addRegistedPlugins(vm *otto.Otto, context context.Context) {
	for _, plugin := range configuredPlugins {
		plugin.SetContext(context)
		registerPlugin(vm, plugin)
	}
}

func registerPlugin(vm *otto.Otto, plugin types.Plugin) {
	for key, function := range plugin.GetFunctions() {
		vm.Set(key, function)
	}
}
