# Rest-proxy

Rest-proxy is a simple system to enable dynamic REST reverse proxying.

It sets up a public server to which clients will connect to request temporary
redirections.

The server part:
* Starts a TCP server and listen for redirection requests
* Starts an HTTP server and redirect calls based on current rules

The client part:
* Connects to the server to ask for a redirection. Keep the connection open.
* Receive redirected HTTP calls from the server through the connection, and
  redirect it to an other locally visible webserver.

The main use case is when you want to develop a REST service on a machine that
is not publicly visible. Using the client will generate a random identifier,
and calls to http://publicserver.com/RandomID/Anything will be redirected
to your local REST server.

This is similar to `ssh -R`, but is simpler to use in some way:
* It only requires ssh access to start the server
* It can dispatch calls to multiple clients from the same server port

## Tutorial

### Client

To use the client, first get the code:

    go get github.com/Gyscos/rest-proxy/rest-proxy-client

Then move to the golang bin directory (or add it to your PATH) and run:

    ./rest-proxy-client -h localhost:80 mypublicserver.com:6666

This will connect to the rest-proxy server (hopefully) running on
mypublicserver.com:6666, ask for a temporary redirection, and print your
token ID. It will keep running to maintain the communication, so don't kill it
until you want to stop the redirection.
Now, calls to mypublicserver.com/token will be redirected to localhost:80 for
your local server to handle.

### Server

To use the server, as usual get the code:

    go get github.com/Gyscos/rest-proxy/rest-proxy-server

And to start the server:

    ./rest-proxy-server -p 6666 -w 80

Will run the web server on port 80 and listen for requests on port 6666
(these are the default options)

## Documentation

### Client

```
Usage: rest-proxy-client [-h webhost[:PORT]] target[:TARGET_PORT]

target[:TARGET_PORT]    Hostname and optionnal port to connect to.
                          Default port is 6666
Options:
  -h webhost[:PORT]     Hostname and optionnal port to redirect calls to.
                          Default is localhost:80
```

### Server

```
Usage: rest-proxy-server [OPTION]...

Options:
  -w WEB_PORT    Run the http server on port WEB_PORT.
                   Defaults to 80
  -p CMD_PORT    Run the request interface on port CMD_PORT.
                   Defaults to 6666
```
