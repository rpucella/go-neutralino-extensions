package main

import (
	"os"
	"errors"
	"github.com/rpucella/go-neutralino-extensions"
)

func main() {

	connInfo, err := neutralinoext.ReadConnInfo(os.Stdin)
	if err != nil {
		panic(err)
	}

	if err := connInfo.StartMessageLoop(processMsg); err != nil {
		panic(err)
	}
}

func processMsg(event string, data any) (map[string]any, error) {
	if event != "echo" {
		return nil, nil
	}
	dataObj, ok := data.(map[string]any)
	if !ok {
		return nil, errors.New("data not an object")
	}
	messageIfc, ok := dataObj["message"]
	if !ok {
		return nil, errors.New("no message field")
	}
	message, ok := messageIfc.(string)
	if !ok {
		return nil, errors.New("mesage not a string")
	}
	result := make(map[string]any)
	result["echo"] = "echoed value: " + message
	return result, nil
}
