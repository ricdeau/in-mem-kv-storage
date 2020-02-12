package service

import (
	"github.com/ricdeau/in-mem-kv-storage/contracts"
	"github.com/ricdeau/in-mem-kv-storage/logger"
	"github.com/ricdeau/in-mem-kv-storage/storage"
	"github.com/ricdeau/in-mem-kv-storage/utils"
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
	requestId := r.Context().Value(utils.RequestId)
	key := r.URL.Path
	value, exists := s.stor.Get(key)
	if !exists {
		logger.Infof("RequestId: %v; Value with key=%s doesn't exist", requestId, key)
		utils.ErrorResponse(http.StatusNotFound, contracts.NotExistsError(key), rw)
		return
	}
	rw.Header().Set(contentType, value.Type)
	_, err := rw.Write(value.Payload)
	if err != nil {
		logger.Errorf("RequestId: %v; Error while writing response: %v", requestId, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Infof("RequestId: %v; Value with key=%s has been retrieved", requestId, key)
}

func (s *service) handlePut(rw http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(utils.RequestId)
	key := r.URL.Path
	contentType := r.Header.Get(contentType)
	if contentType == "" {
		contentType = defaultContentType
	}
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Errorf("RequestId: %v; Error while writing response: %v", requestId, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	isNew := s.stor.Set(key, storage.Data{Type: contentType, Payload: data})
	if isNew {
		rw.WriteHeader(http.StatusCreated)
	}
	logger.Infof("RequestId: %v; Value with key=%s has been saved", requestId, key)
}

func (s *service) handleDelete(rw http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(utils.RequestId)
	key := r.URL.Path
	s.stor.Delete(key)
	rw.WriteHeader(http.StatusNoContent)
	logger.Infof("RequestId: %v; Value with key=%s has been deleted", requestId, key)
}
