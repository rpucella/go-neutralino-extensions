package neutralinoext

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"os"
	"os/signal"
)

// Mostly adapted from https://github.com/gorilla/websocket/blob/main/examples/echo/client.go

type processFn func(string, any) (map[string]any, error)

type ConnInfo struct {
	Port string `json:"nlPort"`
	Token  string `json:"nlToken"`
	ConnToken string `json:"nlConnectToken"`
	ExtId string `json:"nlExtensionId"`
	conn *websocket.Conn
	process processFn
}

/*
   const NL_PORT = processInput.nlPort
   const NL_TOKEN = processInput.nlToken
   const NL_CTOKEN = processInput.nlConnectToken
   const NL_EXTID = processInput.nlExtensionId
*/

func ReadConnInfo(r io.Reader) (ConnInfo, error) {
	br := bufio.NewReader(r)
	connInfoStr, err := br.ReadString('\n')
	log.Println(connInfoStr)
	if err != nil && err != io.EOF {
		return ConnInfo{}, fmt.Errorf("cannot read connection info: %w", err)
	}
	connInfo := ConnInfo{}
	if err := json.Unmarshal([]byte(connInfoStr), &connInfo); err != nil {
		return ConnInfo{}, fmt.Errorf("cannot parse connection info: %w", err)
	}
	return connInfo, nil
}

// Send a message to the app using app.broadcast.
//
// Format of message:
// {
//   "id": <uuid>,
//   "method": "app.broadcast",
//   "accessToken": <token from connInfo>,
//   "data": {
//      "event": <event name>,
//      "data": <data object to broadcast>
//   }
// }
func (ci ConnInfo) SendMessage(event string, data map[string]any) error {
	// Buiild a message.
	dataObj := make(map[string]any)
	dataObj["event"] = event
	dataObj["data"] = data
	msgObj := make(map[string]any)
	msgObj["id"] = uuid.NewString()
	msgObj["method"] = "app.broadcast"
	msgObj["accessToken"] = ci.Token
	msgObj["data"] = dataObj
	// Send it.
	msg, err := json.Marshal(msgObj)
	if err != nil {
		return fmt.Errorf("cannot marshal result: %w", err)
	}
	if ci.conn == nil {
		return fmt.Errorf("in SendMessage: connection is nil")
	}
	ci.conn.WriteMessage(websocket.BinaryMessage, msg)
	return nil
}

func (ci ConnInfo) StartMessageLoop(process processFn) error {

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)

	urlString := fmt.Sprintf("ws://localhost:%s?extensionId=%s&connectToken=%s",
		ci.Port, ci.ExtId, ci.ConnToken)

	conn, _, err := websocket.DefaultDialer.Dial(urlString, nil)
	if err != nil {
		return fmt.Errorf("error connecting to websocket: %w", err)
	}
	defer conn.Close()
	ci.conn = conn
	ci.process = process

	doneCh := make(chan struct{}, 1)
	defer close(doneCh)
	msgCh := make(chan map[string]any, 1)
	defer close(msgCh)

	go ci.readMessages(doneCh, msgCh)

	for {
		select {
		case <-doneCh:
			log.Println("done")
			return nil

		case <-interruptCh:
			log.Println("interrupt received")
			// Tell the app we're closing.
			msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
			err := conn.WriteMessage(websocket.CloseMessage, msg)
			if err != nil {
				return fmt.Errorf("write close: %w", err)
			}
			continue
			// Let the CloseHandler handle it.

		case msgObj := <-msgCh:
			go ci.processMessage(msgObj)
		}
	}
}

func (ci *ConnInfo) processMessage(msgObj map[string]any) {
	eventIfc, ok := msgObj["event"]
	if !ok {
		// This may be a response to an app.broadcast sent in response to a message!
		log.Printf("no event field: %s\n", msgObj)
		return
	}
	event, ok := eventIfc.(string)
	if !ok {
		log.Println("event field not a string")
		return
	}
	dataIfc, ok := msgObj["data"]
	if !ok {
		log.Printf("no data field: %s\n", msgObj)
		return
	}
	msgResult, err := (ci.process)(event, dataIfc)
	if err != nil {
		log.Printf("cannot process message: %v\n", err)
		return
	}
	if len(msgResult) == 0 {
		// No response.
		return
	}
	data, ok := dataIfc.(map[string]any)
	if !ok {
		log.Println("data not an object")
		return
	}
	respIdIfc, ok := data["_respId"]
	if !ok {
		log.Println("no _respId for response")
		return
	}
	msgResult["_respId"] = respIdIfc
	respEventIfc, ok := data["_respEvent"]
	if !ok {
		log.Println("no _respEvent for response")
		return
	}
	respEvent, ok := respEventIfc.(string)
	if !ok {
		log.Printf("_respEvent not a string: %v\n", respEventIfc)
		return
	}
	if err := ci.SendMessage(respEvent, msgResult); err != nil {
		log.Printf("cannot send response: %v\n", err)
		return
	}
}

func (ci *ConnInfo) readMessages(doneCh chan struct{}, msgCh chan map[string]any) error {
	for {
		_, message, err := ci.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("connection closed: %v\n", err)
				doneCh <- struct{}{}
				return nil
			}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("unexpected connection closed: %v\n", err)
				doneCh <- struct{}{}
				return nil
			}
			// What do we do in case of error!? skip + log?
			log.Printf("read: %v\n", err)
			continue
		}
		msgObj := make(map[string]any)
		if err := json.Unmarshal(message, &msgObj); err != nil {
			log.Printf("cannot parse message: %v\n", err)
			continue
		}
		msgCh <- msgObj
	}
}
