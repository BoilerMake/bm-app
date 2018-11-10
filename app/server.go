package app

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
)

// Server is an extended net/http server with handler dependencys injected.
// It uses dependency injection to allow for easier testing using mocks.
// See README.md for more on dependency injection
type Server struct {
	*http.Server
	db *sql.DB
}

// NewServer returns a new Server struct based on a given config.
func NewServer(c *config.Config) (*Server) {
	return &Server {
		Addr: ":" + conf.Port,
		Handler: srv.routes(),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 13, // 8 KiB
	}
}

// Start begins listening for HTTP requests on a server.
// It will attempt to gracefully shutdown on SIGINT or when Server.Stop() is called. // TODO graceful shutdown
func (s *Server) Start() {
	// TODO should this be an env var?
	if c.Mode == "prod" {
		m := autocert.Manager{
			Cache:      autocert.DirCache(conf.CertsPath),
			Prompt:     autocert.AcceptTOS,
			// FIXME put these into a hosts file, config, or env var?
			HostPolicy: autocert.HostWhitelist("boilermake.org", "www.boilermake.org"),
		}

		s.TLSConfig = m.TLSConfig()
		// Redirect all http traffic to https
		go http.ListenAndServe(":80", m.HTTPHandler(nil))

		s.ListenAndServeTLS("", "")
	} else {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", s.Addr, err)
		}
	}
}

// Stop will try to terminate the server gracefully
func (srv *Server) Stop() {

}
