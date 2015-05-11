package webchannel

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type WebChannel struct {
	In  chan []byte
	Out chan []byte
}

func writer(ws *websocket.Conn, webChannel *WebChannel) {
	pingTicker := time.NewTicker(time.Second * 55)

	for {
		select {
		case <-pingTicker.C:
			log.Println("Pinging websocket")
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case msg := <-webChannel.Out:
			ws.WriteMessage(websocket.TextMessage, msg)
		}
	}

	defer pingTicker.Stop()
}

func reader(ws *websocket.Conn, webChannel *WebChannel) {
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(time.Second * 60))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(time.Second * 60))
		return nil
	})

	for {
		_, p, err := ws.ReadMessage()

		if err != nil {
			break
		}

		webChannel.In <- p
	}
}

func handleHttp(webChannel *WebChannel) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		log.Println("Opening socket")
		ws, err := upgrader.Upgrade(response, request, nil)

		if err != nil {
			log.Fatalln(err)
			return
		}

		defer func() {
			log.Println("Closing socket")
			ws.Close()
		}()

		go writer(ws, webChannel)
		reader(ws, webChannel)
	}
}

func New(path string) (*WebChannel, error) {
	wc := &WebChannel{
		In:  make(chan []byte),
		Out: make(chan []byte),
	}

	http.HandleFunc(path, handleHttp(wc))

	return wc, nil
}
