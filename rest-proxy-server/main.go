package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Read flags
	var port int
	var webPort int
	flag.IntVar(&port, "p", 6666, "Port to listen to for requests")
	flag.IntVar(&webPort, "w", 80, "Port to listen to for the web server")
	flag.Parse()

	// Prepare communication
	requestChannels := make(chan Request, 5)

	// Prepare servers
	qs := TcpQueryServer{}
	ws := NewWebServer(requestChannels)

	// Run everything
	go webServer(webPort, ws)
	qs.Run(port, requestChannels)

	// Wait?
}

func webServer(port int, ws *WebServer) {
	log.Println("Starting web server on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), ws))
}
