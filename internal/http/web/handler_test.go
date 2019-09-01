package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

// TestWalkRoutesTemplates mocks a request at each GET endpoint available to
// makes sure they do not error.
func TestWalkRoutesTemplates(t *testing.T) {
	// Set up env config vars
	os.Setenv("JWT_ISSUER", "test")
	os.Setenv("JWT_COOKIE_NAME", "test")
	os.Setenv("JWT_SIGNING_KEY", "test")
	os.Setenv("ENV_MODE", "development")
	os.Setenv("DOMAIN", "testhost")
	os.Setenv("PORT", "8080")
	os.Setenv("WEB_PATH", "../../../web")
	os.Setenv("TEMPLATES_PATH", "../../../templates")
	os.Setenv("SESSION_SECRET", "secretsssh")
	os.Setenv("SESSION_COOKIE_NAME", "session_name")
	handler := NewHandler(nil, nil, nil, nil)

	// This func will be called at every end point in the handler
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		// Only GET handlers render templates, so only check those
		// Also ignore pages that use a {URL param}
		if method == "GET" && !strings.ContainsAny(route, "{}") {
			t.Run(strings.TrimPrefix(route, "/"), func(t *testing.T) {
				// Let tests run in parallel
				t.Parallel()

				req, err := http.NewRequest("GET", route, nil)
				if err != nil {
					t.Errorf("error creating request: %s", err)
				}

				// Record and serve mock request
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)

				if rr.Code == http.StatusInternalServerError {
					t.Errorf("failed to render route: %v", route)
				}
			})
		}

		return nil
	}

	// Runs the walkFunc on every registered route
	if err := chi.Walk(handler, walkFunc); err != nil {
		t.Fatalf("failed while walking: %v", err)
	}
}
