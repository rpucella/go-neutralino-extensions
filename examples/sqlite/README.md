
# Example: SQLite Extension

Three things:

1. this example requires Go and Nodejs
2. this example requires building the `go-sqlite3` package via `CGO_ENABLED=1` and a C compiler; see the [`go-sqlite3`](https://github.com/mattn/go-sqlite3) link for more details
3. this example has only been tested on Mac OS.

To run the example, first build the extension

    go build -o app/extensions/sqlite sqlite.go
    
Build the Neutralinojs app:

    cd app
    npx @neutralinojs/neu build

Run the created app (depending on your architecture):

    ./dist/sqlite-app/echo-app-mac_arm64
    
Profit.

