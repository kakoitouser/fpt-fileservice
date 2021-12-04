package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/kakoitouser/ftp-fileservice/internal/models"
)

const (
	CONN_PORT = ":9090"
	CONN_TYPE = "tcp"

	MSG_DISCONNECT = "Disconnected from the server.\n"

	DOWNLOAD = "download"
	UPLOAD   = "upload"
)

var wg sync.WaitGroup

func Download(conn net.Conn) {
	reader := bufio.NewReader(conn)
	byte, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
		wg.Done()
		return
	}
	file := &models.File{}
	dec := gob.NewDecoder(bytes.NewReader(byte))
	log.Println("i have file")
	dec.Decode(file)
	f, err := os.Create(file.Name)
	if err != nil {
		panic("panic")
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	writer.Write(file.Data)
	defer wg.Done()
}

// Starts up a read and write thread which connect to the server through the
// a socket connection.
func main() {
	wg.Add(1)

	conn, err := net.Dial(CONN_TYPE, CONN_PORT)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}

	var input int
	fmt.Println("you want ? 1-download, 2-upload")
	fmt.Scan(&input)
	switch input {
	case 1:
		var fileUID string
		fmt.Println("file uid: ")
		fmt.Scan(&fileUID)
		conn.Write([]byte(fmt.Sprintf("%s %s\n", DOWNLOAD, fileUID)))

		go Download(conn)
	}

	wg.Wait()
}
