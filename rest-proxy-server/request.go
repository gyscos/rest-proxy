package main

// Represent a new client asking to be plugged
// Carry the ClientConnection to communicate, an
// answer channel to give it the response code,
// and a cancel channel that is used when the connection
// is closed.
type Request struct {
	conn ClientConnection

	confirmation func(token string) error
}
