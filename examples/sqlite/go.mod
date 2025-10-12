module github.com/rpucella/go-neutralino-extension/examples/sqlite

go 1.23.3

require (
	github.com/mattn/go-sqlite3 v1.14.32
	github.com/rpucella/go-neutralino-extension v0.0.0-20251001050614-a8cde5c425f5
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
)

// Redirect module to use the one in this exact repo.
replace github.com/rpucella/go-neutralino-extension => ../..
