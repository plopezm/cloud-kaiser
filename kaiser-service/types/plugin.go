package types

// Plugin This interface should be implemented for every plugin
type Plugin interface {
	// GetFunctions returns the functions and attributes to be registered in the VM
	GetFunctions() map[string]interface{}
}
