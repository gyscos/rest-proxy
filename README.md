Rest-proxy
==========

Rest-proxy is a simple system to enable dynamic REST proxying.

It sets up a public server to which clients will connect to request temporary
redirrections.

The server part:
* Starts a TCP server and listen for redirrection requests
* Starts an HTTP server and redirrect calls based on current rules

The client part:
* Connects to the server to ask for a redirrection. Keep the connection open.
* Receive redirrected HTTP calls from the server through the connection, and
  redirrect it to an other locally visible webserver.

The main use case is when you want to develop a REST service on a machine that
is not publicly visible. Using the client will generate a random identifier,
and calls to http://publicserver.com/RandomID/Anything will be redirrected
to your local REST server.
