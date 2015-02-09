package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	var target string
	var host string

	flag.StringVar(&host, "h", "localhost:8080", "Local web server to proxify")
	flag.Parse()

	if flag.NArg() == 0 {
		// Error! Needs a target
		log.Fatal("No target found")
	}
	target = flag.Arg(0)

	log.Fatal(link(host, target))
}

func ensureHasPort(host string, defaultPort int) string {
	id := strings.Index(host, ":")
	if id < 0 {
		return fmt.Sprintf("%v:%v", host, defaultPort)
	} else if id == len(host)-1 {
		// : as last character. weird.
		return fmt.Sprintf("%v%v", host, defaultPort)
	} else {
		return host
	}
}

func link(host string, target string) error {
	host = ensureHasPort(host, 8080)
	target = ensureHasPort(target, 6666)
	log.Println(host, target)

	conn, err := net.Dial("tcp", target)
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

	log.Println("Decoding...")
	var id string
	err = dec.Decode(&id)
	if err != nil {
		return err
	}
	fmt.Println(id)

	log.Println("Now Serving http")
	return serveHttp(host, dec, enc)
}

func serveHttp(host string, dec *gob.Decoder, enc *gob.Encoder) error {
	for {
		var url string
		err := dec.Decode(&url)
		if err != nil {
			return err
		}
		resp, err := http.Get("http://" + host + url)
		if err != nil {
			return err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		enc.Encode(body)
	}

	return nil
}
