package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

type Target struct {
	// It has a connection
	conn ClientConnection
}

func NewTarget(conn ClientConnection) *Target {
	return &Target{conn: conn}
}

type WebServer struct {
	rc <-chan Request
	cc <-chan string

	targets map[string]*Target
}

func (ws *WebServer) ListenForRequests() {
	for {
		select {
		case rq := <-ws.rc:
			// Read the request
			ws.handleRequest(rq)
		case id := <-ws.cc:
			// Read the cancel
			ws.handleCancel(id)
		}
	}
}

func (ws *WebServer) handleRequest(rq Request) {
	var id string
	for {
		id = randomId(4)
		if _, ok := ws.targets[id]; !ok {
			break
		}
	}
	rq.answer <- id
	ws.AddTarget(id, NewTarget(rq.conn))
}

func (ws *WebServer) handleCancel(id string) {
	// TODO: mutex
	// fmt.Println("Removing", id)
	delete(ws.targets, id)
}

func (ws *WebServer) AddTarget(id string, target *Target) {
	// TODO: mutex
	// fmt.Println("Adding", id)
	ws.targets[id] = target
}

func (ws *WebServer) GetTarget(id string) *Target {
	// TODO: mutex
	return ws.targets[id]
}

// Create a new web handler, and start goroutines to listen on the channels
func NewWebServer(rc <-chan Request, cc <-chan string) *WebServer {
	ws := &WebServer{
		targets: make(map[string]*Target),

		rc: rc,
		cc: cc}

	go ws.ListenForRequests()

	return ws
}

func handleNotFound(w http.ResponseWriter) {
	fmt.Fprintln(w, "Not found.")
}

func (ws *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving http")
	// Find asked ID
	url := r.URL.Path[1:]
	id := strings.Split(url, "/")[0]
	target := ws.GetTarget(id)
	if target == nil {
		// No such client
		handleNotFound(w)
	} else {
		// Ok, we're ready!
		// Find the requested URL
		targetUrl := url[len(id):]

		b, err := target.conn.Serve(targetUrl)
		if err != nil {
			// Error occured
			handleNotFound(w)
		} else {
			w.Write(b)
		}
	}
}

var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomId(l int) string {
	result := make([]rune, l)

	for i := 0; i < l; i++ {
		result[i] = rune(letters[rand.Intn(len(letters))])
	}

	return string(result)
}
