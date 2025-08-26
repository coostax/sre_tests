package keyserver_test

import (
	"github.com/coostax/togetherai_test/keyserver"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestKeyserver(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(keyserver.KeyHandler))
	defer testServer.Close()
    tests := []struct {
        name           string
        method         string
        path           string
        wantStatus     int
        wantBodyLength int
    }{
        {
            name:           "Valid length within max",
            path:           "/key/16",
            wantStatus:     http.StatusOK,
        },
        {
            name:       "Length exceeds max-size",
            path:       "/key/1028",
            wantStatus: http.StatusBadRequest,
        },
		{
            name:       "Non-integer length",
            path:       "/key/abc",
            wantStatus: http.StatusBadRequest,
        },
		{
            name:       "Negative length",
            path:       "/key/-5",
            wantStatus: http.StatusBadRequest,
        },
		{
            name:       "Missing length segment",
            path:       "/key/",
            wantStatus: http.StatusBadRequest,
        },
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			resp, err := http.Get(testServer.URL+tst.path)
			if err != nil {
				t.Errorf("Got error: %v", err)
			}
			if resp.StatusCode != tst.wantStatus {
				t.Errorf("response code is not %d, was: %d", tst.wantStatus, resp.StatusCode)
			}
		})    
	}
}

func TestContentHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/key/8", nil)
    w := httptest.NewRecorder()

    keyserver.KeyHandler(w, req)
    res := w.Result()
    defer res.Body.Close()

    if ct := res.Header.Get("Content-Type"); ct != "text/plain" {
        t.Errorf("Content-Type = %q; want %q", ct, "text/plain")
    }
    if cl := res.Header.Get("Content-Length"); cl != "8" {
        t.Errorf("Content-Length = %q; want %q", cl, "8")
    }
}
