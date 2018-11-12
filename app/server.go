package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

// Server is an extended net/http server with handler dependencys injected.
// It uses dependency injection to allow for easier testing using mocks.
// See README.md for more on dependency injection
type Server struct {
	*http.Server
	//db *sql.DB
}

// NewServer returns a new Server struct based on a given config.
func NewServer() *Server {
	return &Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        srv.routes(),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 13, // 8 KiB
	}
}

// Start begins listening for HTTP or HTTPS requests, depending on the mode given in config.
// It will attempt to gracefully shutdown on SIGINT.
func (s *Server) Start() {
	mode := os.Getenv("ENV_MODE")

	if mode == "production" {
		log.Printf("Starting server in PRODUCTION mode at %s", s.Addr)

		m := autocert.Manager{
			Cache:  autocert.DirCache(conf.CertsPath),
			Prompt: autocert.AcceptTOS,
			// FIXME put these into a hosts file, config, or env var? idk probs one of those
			HostPolicy: autocert.HostWhitelist("boilermake.org", "www.boilermake.org"),
		}

		s.TLSConfig = m.TLSConfig()

		// Redirect all http traffic to https
		go http.ListenAndServe(":80", m.HTTPHandler(nil))

		go func() {
			if err := s.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Could not listen on %s: %v\n", s.Addr, err)
			}
		}()
	} else {
		log.Printf("Starting server in DEVELOPMENT mode at %s", s.Addr)

		go func() {
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Could not listen on %s: %v\n", s.Addr, err)
			}
		}()
	}

	stop := make(chan os.Signal, 1)
	// Fun fact, empty structs take up 0 bytes
	done := make(chan struct{})

	signal.Notify(stop, os.Interrupt)

	go func() {
		<-quit
		log.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}

		//if err := db.Shutdown(); err != nil {
		//	log.Fatal("Could not gracefully shutdown the database")
		//}
		close(done)
	}()

	<-done
	log.Println("Server stopped")
}
