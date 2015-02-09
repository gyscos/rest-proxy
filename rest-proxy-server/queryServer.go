package main

// Represent a new client asking to be plugged
// Carry the ClientConnection to communicate, an
// answer channel to give it the response code,
// and a cancel channel that is used when the connection
// is closed.
type Request struct {
	conn   ClientConnection
	answer chan<- string
}

// QueryServer listens on a port, and sends Requests up its channel.
// Requests carry a link to a ClientConnection that will implement
// the communication
type QueryServer interface {
	Run(port int, rc chan<- Request, cc chan<- string) error
}

// RequestHandler handles a single connection with a client.
type ClientConnection interface {
	Serve(url string) ([]byte, error)
}
