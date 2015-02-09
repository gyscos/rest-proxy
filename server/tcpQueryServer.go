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

func (tcc *TcpClientConnection) Serve(url string) ([]byte, error) {
	log.Println("Serving tcp connection", url)
	err := tcc.enc.Encode(url)
	if err != nil {
		return nil, err
	}

	var input []byte
	log.Println("Scanning response...")
	err = tcc.dec.Decode(&input)
	if err != nil {
		return nil, err
	}

	return input, nil
}

type TcpQueryServer struct {
	rc chan<- Request
	cc chan<- string
}

func (tqs *TcpQueryServer) handleConnection(c net.Conn) {
	tcc := NewTcpClientConnection(c)
	a := make(chan string)
	tqs.rc <- Request{conn: tcc, answer: a}
	id := <-a
	// Return that to the connection
	tcc.enc.Encode(id)
	// Then
}

func (tqs *TcpQueryServer) Run(port int, rc chan<- Request, cc chan<- string) error {
	tqs.rc = rc
	tqs.cc = cc

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
		go tqs.handleConnection(conn)

	}
	return nil
}
