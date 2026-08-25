package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"status":"ok"`) {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestArcLengthEndpoint(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)
	body := `[[0,0],[0.45,1.4],[1.55,1.4],[2,0]]`
	req := httptest.NewRequest(http.MethodPost, "/api/arclength", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"arcLength"`) {
		t.Fatalf("missing arcLength in response: %s", w.Body.String())
	}
}

func TestValidateEndpointInvalid(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)
	body := `[[0,0],[1,1]]`
	req := httptest.NewRequest(http.MethodPost, "/api/validate", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"valid":false`) {
		t.Fatalf("expected invalid: %s", w.Body.String())
	}
}

func TestCheckEndpoint(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)
	body := `[[0,0],[0.45,1.4],[1.55,1.4],[2,0]]`
	req := httptest.NewRequest(http.MethodPost, "/api/check", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"allPass":true`) {
		t.Fatalf("invariants failed: %s", w.Body.String())
	}
}

func TestOffsetEndpoint(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)
	body := `{"controlPoints":[[0,0],[0.45,1.4],[1.55,1.4],[2,0]],"offsetDistance":0.1}`
	req := httptest.NewRequest(http.MethodPost, "/api/offset", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"points"`) {
		t.Fatalf("missing points: %s", w.Body.String())
	}
}

func TestSampleEndpoint(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)
	body := `[[0,0],[0.45,1.4],[1.55,1.4],[2,0]]`
	req := httptest.NewRequest(http.MethodPost, "/api/sample", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"arcLength"`) {
		t.Fatalf("missing arcLength: %s", w.Body.String())
	}
}

func TestMethodNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)
	req := httptest.NewRequest(http.MethodGet, "/api/arclength", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("want 405, got %d", w.Code)
	}
}
