package types

import "context"

// Plugin This interface should be implemented for every plugin
type Plugin interface {
	// GetFunctions returns the functions to be registered in the VM
	GetFunctions() map[string]interface{}
	SetContext(context context.Context)
}
