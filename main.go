package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func createRouter() chi.Router {
	mux := chi.NewRouter()
	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return mux
}

func main() {
	l, err := netListen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failded to listen: %v", err)
	}

	mux := createRouter()
	server := http.Server{
		Handler: mux,
	}

	go func() {
		log.Printf("starting server on %s", l.Addr())
		if err := server.Serve(l); err != nil {
			log.Fatalf("server closed with %v", err)
		}
	}()
}
