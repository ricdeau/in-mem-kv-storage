package service

import (
	"github.com/ricdeau/in-mem-kv-storage/contracts"
	"github.com/ricdeau/in-mem-kv-storage/logger"
	"github.com/ricdeau/in-mem-kv-storage/storage"
	"github.com/ricdeau/in-mem-kv-storage/utils"
	"io/ioutil"
	"net/http"
)

// Headers
const (
	contentType = "Content-Type"
	allow       = "Allow"
)

type service struct {
	stor storage.Storage
}

// New creates http.Handler that performs operations with storage.Storage depending on http method.
// Allowed methods: GET, PUT, DELETE, OPTIONS
func New(route string, stor storage.Storage) http.Handler {
	srv := &service{stor}
	return http.StripPrefix(route, srv)
}

func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleGet(w, r)
	case "PUT":
		s.handlePut(w, r)
	case "DELETE":
		s.handleDelete(w, r)
	case "OPTIONS":
		s.handleOptions(w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *service) handleGet(rw http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	requestIDPrefix := utils.RequestIDPrefix(r.Context())
	value, exists := s.stor.Get(key)
	if !exists {
		logger.Infof("%s; Value with key=%s doesn't exist", requestIDPrefix, key)
		utils.ErrorResponse(http.StatusNotFound, contracts.NotExistsError(key), rw)
		return
	}
	rw.Header().Set(contentType, value.Type)
	_, err := rw.Write(value.Payload)
	if err != nil {
		logger.Errorf("%s; Error while writing response: %v", requestIDPrefix, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Infof("%s; Value with key=%s has been retrieved", requestIDPrefix, key)
}

func (s *service) handlePut(rw http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	requestIDPrefix := utils.RequestIDPrefix(r.Context())
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("%s; Error while writing response: %v", requestIDPrefix, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	contentType := r.Header.Get(contentType)
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	isNew := s.stor.Set(key, &storage.Data{Type: contentType, Payload: data})
	if isNew {
		rw.WriteHeader(http.StatusCreated)
	}
	logger.Infof("%s; Value with key=%s has been saved", requestIDPrefix, key)
}

func (s *service) handleDelete(rw http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	requestIDPrefix := utils.RequestIDPrefix(r.Context())
	s.stor.Delete(key)
	rw.WriteHeader(http.StatusNoContent)
	logger.Infof("%s; Value with key=%s has been deleted", requestIDPrefix, key)
}

func (s *service) handleOptions(rw http.ResponseWriter) {
	rw.Header().Set(allow, "GET, PUT, DELETE, OPTIONS")
}
