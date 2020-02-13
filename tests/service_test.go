package tests

import (
	"bytes"
	"encoding/json"
	"github.com/ricdeau/in-mem-kv-storage/contracts"
	"github.com/ricdeau/in-mem-kv-storage/service"
	"github.com/ricdeau/in-mem-kv-storage/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type storageMock struct {
	get func(string) (*storage.Data, bool)
	set func(string, *storage.Data) bool
	del func(string)
}

func (s *storageMock) Get(key string) (data *storage.Data, exists bool) {
	return s.get(key)
}

func (s *storageMock) Set(key string, value *storage.Data) (new bool) {
	return s.set(key, value)
}

func (s *storageMock) Delete(key string) {
	s.del(key)
}

type testCase struct {
	expectedBody, expectedType string
	expectedCode               int
	key, method                string
	payload                    io.Reader
	stor                       storage.Storage
}

func TestGetExist(t *testing.T) {
	const (
		expectedBody = "expected-data"
		expectedType = "text/plain"
	)
	mock := new(storageMock)
	mock.get = func(s string) (data *storage.Data, b bool) {
		return &storage.Data{Type: expectedType, Payload: []byte(expectedBody)}, true
	}
	tc := &testCase{
		expectedBody: expectedBody,
		expectedType: expectedType,
		expectedCode: http.StatusOK,
		key:          "key",
		method:       "GET",
		payload:      nil,
		stor:         mock,
	}
	testAction(t, tc)
}

func TestGetNotExist(t *testing.T) {
	const (
		key          = "some-key"
		expectedType = "application/json"
	)
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(contracts.NotExistsError(key)); err != nil {
		t.Error(err)
	}

	mock := new(storageMock)
	mock.get = func(s string) (data *storage.Data, b bool) {
		return &storage.Data{}, false
	}
	tc := &testCase{
		expectedBody: buf.String(),
		expectedType: expectedType,
		expectedCode: http.StatusNotFound,
		key:          key,
		method:       "GET",
		payload:      nil,
		stor:         mock,
	}
	testAction(t, tc)
}

func TestPutNew(t *testing.T) {
	const (
		key = "some-key"
	)
	buf := &bytes.Buffer{}
	mock := new(storageMock)
	mock.set = func(string, *storage.Data) bool {
		return true
	}
	tc := &testCase{
		expectedBody: "",
		expectedType: "",
		expectedCode: http.StatusCreated,
		key:          key,
		method:       "PUT",
		payload:      buf,
		stor:         mock,
	}
	testAction(t, tc)
}

func TestPutUpdate(t *testing.T) {
	const (
		key = "some-key"
	)
	buf := &bytes.Buffer{}
	mock := new(storageMock)
	mock.set = func(string, *storage.Data) bool {
		return false
	}
	tc := &testCase{
		expectedBody: "",
		expectedType: "",
		expectedCode: http.StatusOK,
		key:          key,
		method:       "PUT",
		payload:      buf,
		stor:         mock,
	}
	testAction(t, tc)
}

func TestDelete(t *testing.T) {
	const (
		key = "some-key"
	)
	mock := new(storageMock)
	mock.del = func(string) {}
	tc := &testCase{
		expectedBody: "",
		expectedType: "",
		expectedCode: http.StatusNoContent,
		key:          key,
		method:       "DELETE",
		payload:      nil,
		stor:         mock,
	}
	testAction(t, tc)
}

func TestOptions(t *testing.T) {
	srv := service.New("/api/", nil)
	request := httptest.NewRequest("OPTIONS", "/api/", nil)
	recorder := httptest.NewRecorder()
	srv.ServeHTTP(recorder, request)

	const expectedAllow = "GET, PUT, DELETE, OPTIONS"
	actualAllow := recorder.Header().Get("Allow")

	if expectedAllow != actualAllow {
		t.Errorf("Allow header error: want %s but got %s", expectedAllow, actualAllow)
	}
}

func TestInvalidMethod(t *testing.T) {
	mock := new(storageMock)
	tc := &testCase{
		expectedBody: "",
		expectedType: "",
		expectedCode: http.StatusMethodNotAllowed,
		key:          "key",
		method:       "POST",
		payload:      nil,
		stor:         mock,
	}
	testAction(t, tc)
}

func testAction(t *testing.T, tc *testCase) {
	srv := service.New("/api/", tc.stor)
	request := httptest.NewRequest(tc.method, "/api/"+tc.key, tc.payload)
	recorder := httptest.NewRecorder()
	srv.ServeHTTP(recorder, request)

	if tc.expectedCode != recorder.Code {
		t.Errorf("Status error: want %d but got %d", tc.expectedCode, recorder.Code)
	}
	actualType := recorder.Header().Get("Content-Type")
	if tc.expectedType != actualType {
		t.Errorf("Type error: want %s but got %s", tc.expectedType, actualType)
	}
	actualBody := recorder.Body.String()
	if tc.expectedBody != actualBody {
		t.Errorf("Body error: want %s but got %s", tc.expectedBody, actualBody)
	}
}
