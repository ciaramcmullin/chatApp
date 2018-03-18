
package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type connection struct {

	send chan []byte

	h *hub
}

func (c *connection) reader(wg *sync.WaitGroup, wsConn *websocket.Conn) {
	defer wg.Done()

	for {

		_, message, err := wsCon.ReadMessage()
		if err != nil {
			break
		}

		c.h.broadcast <- message
	}
}

func (c *connection) writer(wg *sync.WaitGroup, wsConn *websocket.Conn) {
	defer wg.Done()
	for message := range c.send {
		err := wsConn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

type wsHandler struct {
	h *hub
}

func (wsh wsHandler) ServeHTTP(w.http.RespondWriter, r *http.Request){

	wsConn, err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		log.Printf("error %s", err)
		return
	}

	c := &connection{send: make(chan []byte, 256), h: wsh.h}
	c.h.addConnection(c)
	defer c.h.removeConnection(c)
	var wg sync.WaitGroup
	wg.Add(2)
	go c.writer(&wg, wsConn)
	go c.reader(&wg, wsConn)
	wg.Wait()
	wsConn.Close()
}

