package http

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/plopezm/cloud-kaiser/kaiser-service/engine"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	gohttp "net/http"
	"time"
)

func init() {
	engine.RegisterPlugin(new(Plugin))
}

var netClient = &gohttp.Client{
	Timeout: time.Second * 10,
}

// Response The structure of an Http response
type Response struct {
	StatusCode int
	Body       string
	Headers    map[string][]string
}

// Plugin Implements http function
type Plugin struct {
	context context.Context
}

// GetInstance Creates a new plugin instance with a context
func (plugin *Plugin) SetContext(context context.Context) {
	plugin.context = context
}

// GetFunctions returns the functions to be registered in the VM
func (plugin *Plugin) GetFunctions() map[string]interface{} {
	functions := make(map[string]interface{})
	functions["Http"] = map[string]interface{}{
		"get": plugin.Get,
	}
	return functions
}

// Get Performs an http get request
func (plugin *Plugin) Get(call otto.FunctionCall) otto.Value {
	urlArgument := call.Argument(0)
	headers := call.Argument(1).Object()

	url, err := urlArgument.ToString()
	if err != nil {
		res, _ := call.Otto.ToValue(err.Error())
		return res
	}

	req, err := createRequest("GET", url, nil, headers)
	if err != nil {
		res, _ := call.Otto.ToValue(err.Error())
		return res
	}

	resp, err := sendRequest(req)
	if err != nil {
		res, _ := call.Otto.ToValue(err.Error())
		return res
	}

	valueToReturn, err := call.Otto.ToValue(resp)
	if err != nil {
		res, _ := call.Otto.ToValue(err.Error())
		return res
	}
	return valueToReturn
}

// Post Performs an http post request
func (plugin *Plugin) Post(call otto.FunctionCall) otto.Value {
	urlArgument := call.Argument(0)
	body := call.Argument(1).Object()
	headers := call.Argument(2).Object()

	url, err := urlArgument.ToString()
	if err != nil {
		res, _ := call.Otto.ToValue(err.Error())
		return res
	}

	req, err := createRequest("POST", url, body, headers)
	if err != nil {
		res, _ := call.Otto.ToValue(err.Error())
		return res
	}

	resp, err := sendRequest(req)
	if err != nil {
		res, _ := call.Otto.ToValue(err.Error())
		return res
	}

	valueToReturn, err := call.Otto.ToValue(resp)
	if err != nil {
		res, _ := call.Otto.ToValue(err.Error())
		return res
	}
	return valueToReturn
}

func createRequest(method string, url string, body *otto.Object, headers *otto.Object) (*gohttp.Request, error) {

	var bodyBuffer *bytes.Buffer
	if body != nil {
		var err error
		bodyBytes, err := convertInterfaceToBytes(body)
		if err != nil {
			return nil, err
		}
		bodyBuffer = bytes.NewBuffer(bodyBytes)
	} else {
		bodyBuffer = bytes.NewBuffer([]byte{})
	}
	req, err := gohttp.NewRequest(method, url, bodyBuffer)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for _, headerKey := range headers.Keys() {
			value, _ := headers.Get(headerKey)
			valueString, _ := value.ToString()
			req.Header.Set(headerKey, valueString)
		}
	}

	return req, nil
}

func sendRequest(req *gohttp.Request) (*Response, error) {
	resp, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsResponse = &Response{
		StatusCode: resp.StatusCode,
		Body:       string(respBody),
		Headers:    resp.Header,
	}

	return jsResponse, nil
}

func convertInterfaceToBytes(object interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if object == nil {
		return nil, nil
	}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(object)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
