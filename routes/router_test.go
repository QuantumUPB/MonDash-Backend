package routes

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	router := NewRouter(nil)

	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if strings.TrimSpace(resp.Body.String()) != "ok" {
		t.Fatalf("unexpected body: %q", resp.Body.String())
	}
}

func TestUpdateNodeRequiresAuth(t *testing.T) {
	router := NewRouter(nil)

	body := bytes.NewBufferString(`{"nodes":[{"name":"node","status":"up","stored_key_count":1,"current_key_rate":0.5}]}`)
	req := httptest.NewRequest(http.MethodPost, "/update-node", body)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.Code)
	}
}

func TestUpdateNodeAuthorized(t *testing.T) {
	router := NewRouter(nil)

	body := bytes.NewBufferString(`{"nodes":[{"name":"node","status":"up","stored_key_count":1,"current_key_rate":0.5}]}`)
	req := httptest.NewRequest(http.MethodPost, "/update-node", body)
	req.Header.Set("X-Auth-Token", "Bearer abc")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
}

func TestCORSMiddleware(t *testing.T) {
	router := NewRouter(nil)

	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if origin := resp.Header().Get("Access-Control-Allow-Origin"); origin != "http://localhost:3000" {
		t.Fatalf("expected Access-Control-Allow-Origin http://localhost:3000, got %q", origin)
	}
	if cred := resp.Header().Get("Access-Control-Allow-Credentials"); cred != "true" {
		t.Fatalf("expected Access-Control-Allow-Credentials true, got %q", cred)
	}
}

func TestCORSDisallowedOrigin(t *testing.T) {
	router := NewRouter(nil)

	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	req.Header.Set("Origin", "http://example.com")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if origin := resp.Header().Get("Access-Control-Allow-Origin"); origin != "http://example.com" {
		t.Fatalf("expected Access-Control-Allow-Origin http://example.com, got %q", origin)
	}
	if cred := resp.Header().Get("Access-Control-Allow-Credentials"); cred != "true" {
		t.Fatalf("expected Access-Control-Allow-Credentials true, got %q", cred)
	}
}

func TestCORSPreflight(t *testing.T) {
	router := NewRouter(nil)

	req := httptest.NewRequest(http.MethodOptions, "/healthcheck", nil)
	req.Header.Set("Origin", "http://localhost:4000")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if origin := resp.Header().Get("Access-Control-Allow-Origin"); origin != "http://localhost:4000" {
		t.Fatalf("expected Access-Control-Allow-Origin http://localhost:4000, got %q", origin)
	}
	if cred := resp.Header().Get("Access-Control-Allow-Credentials"); cred != "true" {
		t.Fatalf("expected Access-Control-Allow-Credentials true, got %q", cred)
	}
}

func TestLoginSetsCookie(t *testing.T) {
	router := NewRouter(nil)

	body := bytes.NewBufferString(`{"username":"admin","password":"admin"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/login", body)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	cookie := resp.Result().Cookies()
	if len(cookie) == 0 || cookie[0].Name != "auth_token" {
		t.Fatalf("expected auth_token cookie")
	}
	token := os.Getenv("AUTH_TOKEN")
	if token == "" {
		token = "abc"
	}
	if cookie[0].Value != token {
		t.Fatalf("expected cookie value %s, got %s", token, cookie[0].Value)
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	router := NewRouter(nil)

	body := bytes.NewBufferString(`{"username":"wrong","password":"bad"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/login", body)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.Code)
	}
}

func TestAPIRequiresAuthCookie(t *testing.T) {
	router := NewRouter(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/apps", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.Code)
	}
}

func TestAPIWithValidCookie(t *testing.T) {
	router := NewRouter(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/apps", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: "abc"})
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
}

func TestAppsTimelineRequiresAuthCookie(t *testing.T) {
	router := NewRouter(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/apps-timeline?startTimestamp=2024-01-01T00:00:00Z&endTimestamp=2024-01-02T00:00:00Z", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.Code)
	}
}

func TestAppsTimelineWithValidCookie(t *testing.T) {
	router := NewRouter(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/apps-timeline?startTimestamp=2024-01-01T00:00:00Z&endTimestamp=2024-01-02T00:00:00Z", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: "abc"})
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
}
