package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)


var upg = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}

var clients []websocket.Conn

func echoHandler(w http.ResponseWriter, r *http.Request)  {
    conn, _ := upg.Upgrade(w, r, nil)

    clients = append(clients, *conn)

    for {
        msgType, msg, err := conn.ReadMessage()
        if err != nil {
            return
        }
        fmt.Printf("%s sent : %s\n", conn.RemoteAddr(), string(msg))
        for _, client := range clients {
            if err = client.WriteMessage(msgType, msg); err != nil {
                return
            }
        }
    }
}
func main()  {
    http.HandleFunc("/echo", echoHandler)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "index.html")
    })
    http.ListenAndServe(":9090", nil)
}
