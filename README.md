# Go Neutralinojs extension package

A simple package for writing [Neutralinojs] extensions in Go.

## To use

Add to a new Go project:

    go get https://github.com/rpucella/go-neutralino-extensions
    
Import into your file:

	import "github.com/rpucella/go-neutralino-extensions"

The package imports as `neutralinoext`.

## API

Read the Neutralinojs connection information sent from the app at extension start up time. 

Takes an `io.Reader` as input, usually `os.Stdin`.

    func ReadConnInfo(r io.Reader) (ConnInfo, error)

Send a message to Neutralinojs app via `app.broadcast`.

Takes an event string to sent to the app, and an associated JSON object (as a `map[string]any` map). 

    func (ci ConnInfo) SendMessage(event string, data map[string]any) error
    
Start a message loop to receive messages from the Neutralinojs app.

Blocks until an interrupe is received or the app closes the connection.

Takes a function to process each message, with type:

    func(string, any) (map[string]any, error)

That function expects the name of the event field in the message, and the content  of the data field
(as a unmarshalled JSON object, so it could really be anything). If the processing function returns
a non-nil `map[string]any`, and the message received had a `_respId` field and a `_respEvent` string
field, then a message is sent back to the Neutralinojs app through `app.broadcast` with the content of
`_respEvent` as event and a JSON object marshalled from the `map[string]any` map to which `_respId` is
added with the value the field had in the incoming message. (This is to encode a simple RPC.)

    func (ci ConnInfo) StartMessageLoop(process ProcessFn) error

See the `examples/` folder for two examples: a simple echo extension, and a skeleton [SQLite](https://sqlite.org/) extension.

