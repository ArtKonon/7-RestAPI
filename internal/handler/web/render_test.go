package web

import (
	"net/http/httptest"
	"testing"
)

func TestProfileCookieRoundTrip(t *testing.T) {
	profile := Profile{Name: "Alice", Email: "alice@example.com", Feedback: "Love it!"}
	rec := httptest.NewRecorder()
	saveProfileCookie(rec, profile)

	respCookies := rec.Result().Cookies()
	if len(respCookies) == 0 {
		t.Fatal("expected profile cookie to be set")
	}

	req := httptest.NewRequest("GET", "/login", nil)
	req.AddCookie(respCookies[0])

	got := readProfileCookie(req)
	if got.Name != profile.Name || got.Email != profile.Email || got.Feedback != profile.Feedback {
		t.Fatalf("profile round trip mismatch: got %#v want %#v", got, profile)
	}
}
