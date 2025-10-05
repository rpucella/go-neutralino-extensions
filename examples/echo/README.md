
# Example: Echo Extension

One thing:

1. this example requires Go and Nodejs

To run the example, first build the extension

    go build -o app/extensions/echo echo.go
    
Build the Neutralinojs app:

    cd app
    npx @neutralinojs/neu build

Run the created app (depending on your architecture):

    ./dist/echo-app/echo-app-mac_arm64
    
Profit.
