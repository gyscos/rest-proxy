package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

// A Target represents a client
type Target struct {
	// It has a connection
	conn ClientConnection
}

func NewTarget(conn ClientConnection) *Target {
	return &Target{conn: conn}
}

// Frontend for the rest-proxy server.
// Redirrects http calls to the correct target
type WebServer struct {
	rc <-chan Request
	cc chan string

	targets map[string]*Target

	mutex sync.RWMutex
}

// Main control loop. Listens for requests and update the rules.
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

// Add the given client to the target pool
func (ws *WebServer) handleRequest(rq Request) {
	var id string
	// Try to find a correct id
	for {
		id = randomId(4)
		if _, ok := ws.targets[id]; !ok {
			break
		}
	}
	err := rq.confirmation(id)
	if err != nil {
		// Error occured?
		log.Println("Error on request confirmation:", err)
	} else {
		ws.AddTarget(id, NewTarget(rq.conn))
	}
}

// Remove the given target
func (ws *WebServer) handleCancel(id string) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	// fmt.Println("Removing", id)
	log.Println("Removing target", id)
	delete(ws.targets, id)
}

func (ws *WebServer) AddTarget(id string, target *Target) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	// fmt.Println("Adding", id)
	log.Println("New target:", id)
	ws.targets[id] = target
}

func (ws *WebServer) GetTarget(id string) *Target {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()

	return ws.targets[id]
}

// Create a new web handler, and start goroutines to listen on the channels
func NewWebServer(rc <-chan Request) *WebServer {
	ws := &WebServer{
		rc:      rc,
		cc:      make(chan string, 5),
		targets: make(map[string]*Target),
	}

	go ws.ListenForRequests()

	return ws
}

func handleNotFound(w http.ResponseWriter) {
	fmt.Fprintln(w, "Not found.")
}

// Main frontent event handler.
func (ws *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
			if err == io.EOF {
				// Connection lost. Let's get rid of the guy
				ws.cc <- id
			}
			// Error occured
			handleNotFound(w)
		} else {
			w.Write(b)
		}
	}
}
