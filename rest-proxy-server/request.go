package main

// Represent a new client asking to be plugged
// Carry the ClientConnection to communicate, and a confirmation function that
// will be used by the webserver when the target token is ready.
type Request struct {
	conn ClientConnection

	confirmation func(token string) error
}
