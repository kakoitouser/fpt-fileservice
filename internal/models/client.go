package models

import "net"

type Client struct {
	Conn     net.Conn
	Password string
}
