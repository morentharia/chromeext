package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type msg map[string]interface{}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", rootHandler)
	logrus.Info("gogogo")

	panic(http.ListenAndServe(":1337", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// content, err := ioutil.ReadFile("index.html")
	// if err != nil {
	// 	fmt.Println("Could not open file.", err)
	// }
	fmt.Fprintf(w, "%s", "<pre>just test</pre>")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Header.Get("Origin") != "http://"+r.Host {
	// 	http.Error(w, "Origin not allowed", 403)
	// 	return
	// }
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	go echo(conn)
}

func echo(conn *websocket.Conn) {
	for {
		m := msg{
			"message": "eval",
			"code":    "alert(\"wow\")",
		}

		fmt.Printf("send message: %#v\n", m)
		if err := conn.WriteJSON(m); err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(2 * time.Second)

		// err := conn.ReadJSON(&m)
		// if err != nil {
		// 	fmt.Println("Error reading json.", err)
		// 	return
		// }

	}
}
