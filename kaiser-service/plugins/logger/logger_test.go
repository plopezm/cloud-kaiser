package logger

import (
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
	os.Mkdir(contextvars.DefaultLogFolder, 0755)
}

func cleanTestFolder() {
	os.RemoveAll(contextvars.DefaultLogFolder)
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
	vm.Set(contextvars.JobName, "test")
	vm.Set(contextvars.JobVersion, "v1")
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
	_, err := vm.Run("Logger.info(\"hello world\")")

	// Then
	assert.Equal(t, err, nil, "Error should be null")
	assert.Equal(t, fileExists(contextvars.DefaultLogFolder, "test_v1.log"), true, "The logfile was not created")
	assert.Equal(t, fileContainsLine(contextvars.DefaultLogFolder+"/test_v1.log", "hello world"), true, "The log message was not written")

	cleanTestFolder()
}
