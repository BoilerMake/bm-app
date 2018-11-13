package app

import (
	"context"
	"fmt"
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
	srv := &Server{
		&http.Server{
			Addr:           ":" + os.Getenv("PORT"),
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    15 * time.Second,
			MaxHeaderBytes: 1 << 13, // 8 KiB
		},
	}

	// TODO init db here

	// Register routes
	// NOTE this doesn't seem like the cleanest way to do this, I'm open to solutions
	srv.routes()

	return srv
}

// Start begins listening for HTTP or HTTPS requests, depending on the mode given in config.
// It will attempt to gracefully shutdown on SIGINT.
func (s *Server) Start() (err error) {
	mode, ok := os.LookupEnv("ENV_MODE")
	if !ok {
		return fmt.Errorf("environment variable not set: %v", "ENV_MODE")
	}

	if mode == "production" {
		log.Printf("Starting server in PRODUCTION mode at %s", s.Addr)

		// NOTE this is likely not the best way to store our certs, using something like crypto/x509 may be better
		certsPath, ok := os.LookupEnv("CERTS_DIR")
		if !ok {
			return fmt.Errorf("environment variable not set: %v", "CERTS_DIR")
		}

		m := autocert.Manager{
			Cache:  autocert.DirCache(certsPath),
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
	} else if mode == "development" {
		log.Printf("Starting server in DEVELOPMENT mode at %s", s.Addr)

		go func() {
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Could not listen on %s: %v\n", s.Addr, err)
			}
		}()
	}

	// Fun fact, empty structs take up 0 bytes
	done := make(chan struct{})
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		<-stop
		log.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		s.SetKeepAlivesEnabled(false)
		if err := s.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}

		//if err := db.Shutdown(); err != nil {
		//	log.Fatal("Could not gracefully shutdown the database")
		//}
		close(done)
	}()

	<-done
	log.Println("Server stopped")

	return err
}
