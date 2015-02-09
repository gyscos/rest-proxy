package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type TcpClientConnection struct {
	conn net.Conn
	enc  *gob.Encoder
	dec  *gob.Decoder
}

func NewTcpClientConnection(conn net.Conn) *TcpClientConnection {
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	return &TcpClientConnection{conn, enc, dec}
}

// Called by the WebServer when it needs a page served.
func (tcc *TcpClientConnection) Serve(url string) ([]byte, error) {
	// Send the GET request
	err := tcc.enc.Encode(url)
	if err != nil {
		return nil, err
	}

	// Read the answer
	var input []byte
	err = tcc.dec.Decode(&input)
	if err != nil {
		return nil, err
	}

	return input, nil
}

// Main backend server. Listens for incoming connections, and pipe requests
// through its channels to the WebServer.
type TcpQueryServer struct {
	// Request Channel. Used to add a new target.
	rc chan<- Request
}

// Some new guy just connected. We'll take care of him.
func (tqs *TcpQueryServer) handleConnection(c net.Conn) {
	tcc := NewTcpClientConnection(c)

	// Ask a token to the server
	tqs.rc <- Request{conn: tcc, confirmation: func(token string) error {
		// Return that to the connection
		return tcc.enc.Encode(token)
	}}
}

// Main backend loop. Listens for incoming connection and serve them in
// goroutines.
func (tqs *TcpQueryServer) Run(port int, rc chan<- Request) error {
	tqs.rc = rc

	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}

	log.Println("Listening for queries on port", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		tqs.handleConnection(conn)

	}
	return nil
}
