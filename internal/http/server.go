package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

// A Server wraps an http.Server and provides some additional functionality.
type Server struct {
	*http.Server
}

// NewServer creates a new Server.
func NewServer(address string, h *Handler) (s *Server) {
	return &Server{
		&http.Server{
			Addr:           address,
			Handler:        h,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    15 * time.Second,
			MaxHeaderBytes: 1 << 13, // 8 KiB
		},
	}
}

// Start begins listening for HTTP or HTTPS requests, depending on the
// environment mode. It will attempt to gracefully shutdown on SIGINT.
// TODO Maybe this shouldn't block like it does right now.  If we go for that
// approach then we'll have to block execution some other way in main.go.
func (s *Server) Start() (err error) {
	mode, ok := os.LookupEnv("ENV_MODE")
	if !ok {
		return fmt.Errorf("environment variable not set: %v", "ENV_MODE")
	}

	if mode == "production" {
		certsDir, ok := os.LookupEnv("CERTS_DIR")
		if !ok {
			return fmt.Errorf("environment variable not set: %v", "CERTS_DIR")
		}

		allowedHosts, ok := os.LookupEnv("ALLOWED_HOSTS")
		if !ok {
			return fmt.Errorf("environment variable not set: %v", "ALLOWED_HOSTS")
		}
		splitHosts := strings.Split(allowedHosts, ", ")

		m := autocert.Manager{
			Cache:      autocert.DirCache(certsDir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(splitHosts...),
		}

		s.TLSConfig = m.TLSConfig()

		// Redirects all http traffic to https
		go http.ListenAndServe(":80", m.HTTPHandler(nil))

		go func() {
			err := s.ListenAndServeTLS("", "")
			if err != nil && err != http.ErrServerClosed {
				log.Fatalf("could not listen on %s: %v\n", s.Addr, err)
			}
		}()
	} else if mode == "development" {
		// Not using HTTPS in development.  Dealing with self signed certs is
		// kinda dog trash
		go func() {
			err := s.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.Fatalf("could not listen on %s: %v\n", s.Addr, err)
			}
		}()
	} else {
		return fmt.Errorf("uknown environment mode: %s", mode)
	}

	log.Printf("started server in %s mode at %s", strings.ToUpper(mode), s.Addr)

	done := make(chan struct{})
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// NOTE shutting the server down this way means we can't have a server.Stop().
	// We can rearchitecure the channels above to be members of a server, and then
	// server.Stop() would be possible. We could also always do
	// `syscall.Kill(syscall.Getpid(), syscall.SIGINT)` but that doesn't seem as
	// clean to me.
	go func() {
		<-stop
		log.Println("server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		s.SetKeepAlivesEnabled(false)
		if err := s.Shutdown(ctx); err != nil {
			log.Fatalf("could not gracefully shutdown the server: %v\n", err)
		}

		close(done)
	}()

	<-done
	log.Println("server stopped")

	return err
}
