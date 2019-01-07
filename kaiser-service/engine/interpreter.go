package engine

import (
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
func NewVM() *otto.Otto {
	vm := otto.New()
	addRegistedPlugins(vm)
	return vm
}

func addRegistedPlugins(vm *otto.Otto) {
	for _, plugin := range configuredPlugins {
		registerPlugin(vm, plugin)
	}
}

func registerPlugin(vm *otto.Otto, plugin types.Plugin) {
	for key, function := range plugin.GetFunctions() {
		vm.Set(key, function)
	}
}
