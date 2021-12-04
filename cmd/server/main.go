package main

import (
	"log"

	"github.com/kakoitouser/ftp-fileservice/internal/server"
)

func main() {
	srv, err := server.NewTcpServer("localhost:9090")
	log.Println("[INFO] TCP SERVER STARTED")
	if err != nil {
		panic(err)
	}
	defer srv.Listner.Close()
	for {
		conn, err := srv.Listner.Accept()
		if err != nil {
			return
		}
		client := srv.NewClient(conn)
		srv.AddClient(client)
		log.Println("new client added")
		go srv.HandleUserRequest(client)
	}
}
