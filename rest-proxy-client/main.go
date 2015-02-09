package main

import (
	"flag"
	"log"
)

func main() {
	// Test
	var target string
	var host string

	// Read arguments
	flag.StringVar(&host, "h", "localhost:8080", "Local web server to proxify")
	flag.Parse()

	if flag.NArg() == 0 {
		// Error! Needs a target.
		log.Fatal("No target found")
	}
	target = flag.Arg(0)

	// And actually link to the target
	client := newClient(host)
	err := client.Connect(target)
	if err != nil {
		log.Fatal(err)
	}
}
