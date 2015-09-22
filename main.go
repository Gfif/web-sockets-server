package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var X int64 = 1000

func readHeartbeat(conn *websocket.Conn) chan int {
	resChan := make(chan int)
	go func() {
		for {
			_, _, err := conn.NextReader()
			if err != nil {
				log.WithError(err).Warn("heartbeat is lost")
				return
			}
		}
		resChan <- 1
	}()

	return resChan
}

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

	hbChan := readHeartbeat(conn)
	select {
	case <-hbChan:
		log.Info("Heartbeat recieved")
	case <-time.After(time.Millisecond * 3 * time.Duration(X)):
		log.Error("Timeout")
		return
	}

	// for _ = range time.NewTicker(time.Second * 10).C {
	// 	if err = conn.WriteMessage(1, []byte("Some message")); err != nil {
	// 		log.WithError(err).Error("Failed to write message")
	// 		return
	// 	}
	// }
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
