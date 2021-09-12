package main

import (
	"fmt"
	"log"
	"net"

	"github.com/lestrrat-go/server-starter/listener"
)

const port = 8080

func netListen(network, addr string) (net.Listener, error) {
	ls, err := listener.ListenAll()
	if err != nil {
		net.Listen(network, addr)
	}
	return ls[0], nil
}

func main() {
	l, err := netListen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failded to listen: %v", err)
	}

	go func() {
		log.Printf("starting server on %s", l.Addr())
		if err := server.Serve(l); err != nil {
			log.Fatalf("server closed with %v", err)
		}
	}()
}
