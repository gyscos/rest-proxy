# Rest-proxy

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

## Tutorial

### Client

To use the client, first get the code:
    go get github.com/Gyscos/rest-proxy/rest-proxy-client

Then move to the golang bin directory (or add it to your PATH) and run:
    ./rest-proxy-client -h localhost:80 mypublicserver.com:6666

This will connect to the rest-proxy server (hopefully) running on
mypublicserver.com:6666, ask for a temporary redirrection, and print your
token ID. It will keep running to maintain the communication, so don't kill it
until you want to stop the redirrection.
Now, calls to mypublicserver.com/token will be redirrected to localhost:80 for
your local server to handle.

### Server

To use the server, as usual get the code:
    go get github.com/Gyscos/rest-proxy/rest-proxy-server

And to start the server:
    ./rest-proxy-server -p 6666 -w 80

Will run the web server on port 80 and listen for requests on port 6666
(these are the default options)
