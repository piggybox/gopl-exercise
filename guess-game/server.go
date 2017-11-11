package main

import (
	"fmt"
	"golang.org/x/net/websocket" // go get -u golang.org/x/net
	"log"
	"net/http"
)

var (
	JSON          = websocket.JSON              // codec for JSON
	Message       = websocket.Message           // codec for string, []byte
	ActiveClients = make(map[ClientConn]string) // map containing clients
	User          = make(map[string]string)
)

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  string
}

func echoHandler(ws *websocket.Conn) {
	var err error
	for {
		var reply string
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("can't receive")
			break
		}

		client := ws.Request().RemoteAddr
		fmt.Println("Client connected:", client)

		sockCli := ClientConn{ws, client}
		ActiveClients[sockCli] = client
		fmt.Println("Number of clients connected:", len(ActiveClients))

		if ActiveClients[sockCli] != "" {
			for cs, na := range ActiveClients {
				if na != "" {
					if err = websocket.Message.Send(cs.websocket, reply); err != nil {
						log.Println("Could not send message to ", cs.clientIP, err.Error())
					}
				}
			}
		}
	}

}

func main() {
	// Routes
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
