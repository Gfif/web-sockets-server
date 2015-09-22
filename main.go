package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var X int64 = 1000

func handler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Error("Failed to connect")
		return
	}

	// defer closing of connection
	defer func() {
		if err := conn.Close(); err != nil {
			log.WithError(err).Error("Failed to close connection")
		} else {
			log.Info("Connection closed")
		}
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.WithError(err).Error("Failed to read message")
			return
		}
		if err = conn.WriteMessage(messageType, p); err != nil {
			log.WithError(err).Error("Failed to write message")
			return
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/public/").
		Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	r.HandleFunc("/echo", handler)
	http.Handle("/", r)

	log.Info("Start server...")
	if err := http.ListenAndServe("127.0.0.1:80", nil); err != nil {
		log.WithError(err).Error("Can't listen and serve")
	}

}
