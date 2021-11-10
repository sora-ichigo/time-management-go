package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time_management_slackapp/app/di"
	"time_management_slackapp/app/server"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/lestrrat-go/server-starter/listener"
)

const port = 8080

func netListen(network, addr string) (net.Listener, error) {
	ls, err := listener.ListenAll()
	if err == nil {
		return ls[0], nil
	}
	return net.Listen(network, addr)
}

func createRouter(contentHandler server.ContentServer) chi.Router {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders: []string{"Accepts", "Authorization", "Content-Type", "X-CSRF-Token"},
	}).Handler)

	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Route("/contents", func(mux chi.Router) {
		mux.Get("/", contentHandler.GetContents)
		mux.Post("/", contentHandler.CreateContent)
		mux.Route("/{contentID}", func(mux chi.Router) {
			mux.Put("/", contentHandler.UpdateContent)
			mux.Delete("/", contentHandler.DeleteContent)
		})
	})
	return mux
}

func main() {
	run()
}

func run() {
	ctx := context.Background()

	app, cleanup, err := di.NewApp(ctx)
	if err != nil {
		log.Fatalf("server closed with %v", err)
		return
	}
	defer cleanup()

	mux := createRouter(app.ContentServer)
	svr := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		log.Printf("starting server on %s", svr.Addr)
		if err := svr.ListenAndServe(); err != nil {
			log.Fatalf("server closed with %v", err)
			return
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)

	log.Printf("SIGNAL %v received, then shutting down...", <-quit)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		log.Fatalf("failed to graceful shutdown: %v", err)
	}
	log.Printf("server shutdown")
}
