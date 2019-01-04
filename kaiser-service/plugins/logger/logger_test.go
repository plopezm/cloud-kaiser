package logger

import (
	"context"
	"github.com/magiconair/properties/assert"
	"github.com/plopezm/cloud-kaiser/kaiser-service/contextvars"
	"github.com/plopezm/cloud-kaiser/kaiser-service/types"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func createTestFolder() {
	os.Mkdir("logs", 0755)
}

func cleanTestFolder() {
	os.RemoveAll("logs")
}

func fileExists(folder string, filename string) bool {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return false
	}

	for _, f := range files {
		if filename == f.Name() {
			return true
		}
	}
	return false
}

func fileContainsLine(filepath string, line string) bool {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return false
	}
	content := string(bytes[:len(bytes)-1])
	return strings.Contains(content, line)
}

func initializeVM(vm *otto.Otto, plugin types.Plugin) {
	ctx := context.WithValue(context.Background(), contextvars.JobName, "test")
	ctx = context.WithValue(ctx, contextvars.JobVersion, "v1")
	plugin.SetContext(ctx)
	for key, function := range plugin.GetFunctions() {
		vm.Set(key, function)
	}
}

func TestLogPlugin_Info(t *testing.T) {
	// Given
	vm := otto.New()
	createTestFolder()
	initializeVM(vm, new(LogPlugin))

	// When
	_, err := vm.Run("Logger.info('hello world')")

	// Then
	assert.Equal(t, err, nil)
	assert.Equal(t, fileExists("logs", "test_v1.log"), true)
	assert.Equal(t, fileContainsLine("logs/test_v1.log", "hello world"), true)

	cleanTestFolder()
}
