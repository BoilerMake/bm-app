package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRoot(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(NewHandler(nil).getRoot())
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler gave wrong status code: got %v expected %v", status, http.StatusOK)
	}
}
