This package wraps the [Gorilla Websocket library](https://github.com/gorilla/websocket) to provide very simple io over a websocket connection using Go channels.

### Install

`go get github.com/emarref/webchannel`

### Usage

```go
// main.go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/emarref/webchannel"
)

func main() {
    // Open the websocket
    wc, err := webchannel.New("/socket")

    if err != nil {
        log.Fatalln(err)
    }

    // Example: Echo back any message received on the websocket
    go func() {
        for msg := range wc.In {
            wc.Out <- msg
        }
    }()

    // Example: Send the time down the websocket every 3 seconds
    pinger := time.NewTicker(time.Second * 3)

    go func() {
        for t := range pinger.C {
            wc.Out <- []byte(fmt.Sprintf("%s", t))
        }
    }()

    // Serve a dummy webpage with a websocket implemented in JS
    http.Handle("/", http.FileServer(http.Dir("./public")))
    http.ListenAndServe(":8080", nil)
}
```

```html
<!-- public/index.html -->
<!DOCTYPE html>
<html lang="en">
    <head>
        <title>Websocket</title>
    </head>
    <body>
        <script>
var connection = new WebSocket('ws://localhost:8080/socket');

// When the connection is open, send some data to the server
connection.onopen = function () {
    console.info("Websocket opened");
};

// Log errors
connection.onerror = function (error) {
    console.err('WebSocket error', error);
};

// Log messages from the server
connection.onmessage = function (e) {
    console.log('Websocket Message', e.data);
};

connection.onclose = function (e) {
    console.info('Websocket closed')
}

        </script>
    </body>
</html>
```
