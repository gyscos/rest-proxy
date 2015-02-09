package main

import (
	"fmt"
	"log"
	"net/http"
)

type StubConnection struct {
}

func (sc *StubConnection) Serve(url string) ([]byte, error) {
	fmt.Println(url)
	return []byte(url), nil
}

type StubServer struct {
}

func (s *StubServer) Run(port int, rc chan<- Request, cc chan<- string) error {
	http.HandleFunc("/del/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/del/"):]
		cc <- id
	})

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		sc := &StubConnection{}
		a := make(chan string)
		rc <- Request{conn: sc, answer: a}
		id := <-a
		fmt.Fprintln(w, id)
	})
	log.Println("Listening for queries on port", port)
	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
