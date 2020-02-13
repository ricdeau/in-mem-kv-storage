package tests

import (
	"bytes"
	"github.com/ricdeau/in-mem-kv-storage/middleware"
	"github.com/ricdeau/in-mem-kv-storage/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestId(t *testing.T) {
	buf := bytes.NewBufferString("")
	log.SetOutput(buf)
	request := httptest.NewRequest("GET", "/api/", nil)
	handler := func(rw http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(utils.RequestID)
		if requestID == nil {
			t.Errorf("RequestID hasn't been set")
		}
	}
	var logging http.Handler = middleware.RequestIDMiddleware(http.HandlerFunc(handler))
	logging.ServeHTTP(httptest.NewRecorder(), request)

	logLine := buf.String()
	if !strings.Contains(logLine, "RequestID") {
		t.Errorf("Log line doesn't contain 'RequestID'")
	}
	if !strings.Contains(logLine, "Method: GET, path: /api/, size: 0") {
		t.Errorf("Log line doesn't contain 'Method: GET, path: /api/, size: 0'")
	}
}

func TestLimitsKeyToLong(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 1000))
	request := httptest.NewRequest("PUT", "/api/key", buf)
	handler := func(http.ResponseWriter, *http.Request) {}
	var limits http.Handler = middleware.LimitsMiddleware(http.HandlerFunc(handler), "/api/", 2, 5000)
	recorder := httptest.NewRecorder()
	limits.ServeHTTP(recorder, request)
	const expectedBody = `{"error":"Provided key size (3 bytes) exceeds maximum allowed key size (2 bytes)."}`
	actualBody := strings.Trim(recorder.Body.String(), "\n")

	if recorder.Code != http.StatusRequestURITooLong {
		t.Errorf("Status error: want %d\nbut got %d", http.StatusRequestURITooLong, recorder.Code)
	}

	if expectedBody != actualBody {
		t.Errorf("Body error: want %s\nbut got %s", expectedBody, actualBody)
	}
}

func TestLimitsValueToLong(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 1000))
	request := httptest.NewRequest("PUT", "/api/key", buf)
	handler := func(http.ResponseWriter, *http.Request) {}
	var limits http.Handler = middleware.LimitsMiddleware(http.HandlerFunc(handler), "/api/", 10, 500)
	recorder := httptest.NewRecorder()
	limits.ServeHTTP(recorder, request)
	const expectedBody = `{"error":"Provided value size (1000 bytes) exceeds maximum allowed value size (500 bytes)."}`
	actualBody := strings.Trim(recorder.Body.String(), "\n")

	if recorder.Code != http.StatusRequestURITooLong {
		t.Errorf("Status error: want %d\nbut got %d", http.StatusRequestURITooLong, recorder.Code)
	}

	if expectedBody != actualBody {
		t.Errorf("Want %s\nbut got %s", expectedBody, actualBody)
	}
}

func TestLimitsValid(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 1000))
	request := httptest.NewRequest("PUT", "/api/key", buf)
	const expectedBody = "pass"
	handler := func(rw http.ResponseWriter, r *http.Request) {
		_, _ = rw.Write([]byte(expectedBody))
	}
	var limits http.Handler = middleware.LimitsMiddleware(http.HandlerFunc(handler), "/api/", 10, 5000)
	recorder := httptest.NewRecorder()
	limits.ServeHTTP(recorder, request)
	actualBody := strings.Trim(recorder.Body.String(), "\n")

	if recorder.Code != http.StatusOK {
		t.Errorf("Status error: want %d\nbut got %d", http.StatusOK, recorder.Code)
	}

	if expectedBody != actualBody {
		t.Errorf("Want %s\nbut got %s", expectedBody, actualBody)
	}
}
