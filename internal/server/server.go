package server

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/kakoitouser/ftp-fileservice/internal/db/inmemory"
	"github.com/kakoitouser/ftp-fileservice/internal/models"
)

const (
	Download = "download"
	Upload   = "upload"
)

type TcpServer struct {
	Listner     net.Listener
	Mu          *sync.Mutex
	Connections []*models.Client
}

func (s *TcpServer) NewClient(conn net.Conn) *models.Client {
	return &models.Client{
		Conn: conn,
	}
}

func NewTcpServer(Port string) (*TcpServer, error) {
	listner, err := net.Listen("tcp", Port)
	if err != nil {
		return nil, err
	}
	tcpServer := &TcpServer{
		Listner:     listner,
		Mu:          &sync.Mutex{},
		Connections: make([]*models.Client, 0),
	}
	return tcpServer, err
}
func (s *TcpServer) HandleUserRequest(client *models.Client) {
	defer client.Conn.Close()
	for {
		userInput, err := bufio.NewReader(client.Conn).ReadString('\n')
		if err != nil {
			s.RemoveClient(client)
			return
		}
		if strings.Contains(userInput, Download) {
			fileUID := strings.Split(userInput, " ")[1]
			file, err := inmemory.GetFileByUID(fileUID)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(file.UID, file.Name)
			go s.UploadFile(file, client)
		}
	}
}

func (s *TcpServer) AddClient(client *models.Client) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Connections = append(s.Connections, client)
}

func (s *TcpServer) RemoveClient(client *models.Client) {
	var i int
	for i = range s.Connections {
		if s.Connections[i] == client {
			break
		}
	}
	s.Connections = append(s.Connections[:i], s.Connections[i+1:]...)
}

func (s *TcpServer) UploadFile(file *models.File, client *models.Client) {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	enc.Encode(file)
	buff.Write([]byte("\n"))

	log.Println("file uploaded")
	client.Conn.Write(buff.Bytes())
}
