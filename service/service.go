package service

import (
	"encoding/json"
	"github.com/ricdeau/in-mem-kv-storage/contracts"
	"github.com/ricdeau/in-mem-kv-storage/storage"
	"io/ioutil"
	"net/http"
)

const (
	GET    = "GET"
	PUT    = "PUT"
	DELETE = "DELETE"
)

const (
	contentType        = "Content-Type"
	defaultContentType = "application/octet-stream"
	jsonContentType    = "application/json"
)

type service struct {
	stor storage.Storage
}

func New(route string, stor storage.Storage) http.Handler {
	srv := &service{stor}
	return http.StripPrefix(route, srv)
}

func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case GET:
		s.handleGet(w, r)
	case PUT:
		s.handlePut(w, r)
	case DELETE:
		s.handleDelete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *service) handleGet(rw http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	value, exists := s.stor.Get(key)
	if !exists {
		errorResponse(http.StatusNotFound, contracts.NotExistsError(key), rw)
		return
	}
	rw.Header().Set(contentType, value.Type)
	_, err := rw.Write(value.Payload)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *service) handlePut(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	contentType := r.Header.Get(contentType)
	if contentType == "" {
		contentType = defaultContentType
	}
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	isNew := s.stor.Set(key, storage.Data{Type: contentType, Payload: data})
	if isNew {
		w.WriteHeader(http.StatusCreated)
	}
}

func (s *service) handleDelete(rw http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	s.stor.Delete(key)
	rw.WriteHeader(http.StatusNoContent)
}

func errorResponse(status int, error *contracts.Error, rw http.ResponseWriter) {
	rw.Header().Set(contentType, jsonContentType)
	rw.WriteHeader(status)
	err := json.NewEncoder(rw).Encode(error)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
