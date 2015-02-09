package main

// QueryServer listens on a port, and sends Requests up its channel.
// Requests carry a link to a ClientConnection that will implement
// the communication
type QueryServer interface {
	Run(port int, rc chan<- Request) error
}

// RequestHandler handles a single connection with a client.
type ClientConnection interface {
	Serve(url string) ([]byte, error)
}
