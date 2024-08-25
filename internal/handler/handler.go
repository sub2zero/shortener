package handler

import (
	"encoding/json"
	"net/http"
	store "shortener/internal/struct"

	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
)

type urlstore interface {
	Add(name string, shortUrl store.ShortUrls) error
	Get(name string) (store.ShortUrls, error)
	List() (map[string]store.ShortUrls, error)
	Update(name string, shortUrl store.ShortUrls) error
	Remove(name string) error
}
type UrlHandler struct {
	store urlstore
}

func NewUrlHandler(s urlstore) *UrlHandler {
	return &UrlHandler{
		store: s,
	}
}
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
func (h UrlHandler) CreateUrl(w http.ResponseWriter, r *http.Request) {
	var url store.ShortUrls

	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	resourceID := slug.Make(url.Full)

	if err := h.store.Add(resourceID, url); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (h UrlHandler) ListUrls(w http.ResponseWriter, r *http.Request) {
	url, err := h.store.List()

	jsonBytes, err := json.Marshal(url)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h UrlHandler) GetUrl(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	url, err := h.store.Get(id)
	if err != nil {
		if err == store.NotFoundErr {
			NotFoundHandler(w, r)
			return
		}

		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(url)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h UrlHandler) UpdateUrl(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Recipe object that will be populated from json payload
	var url store.ShortUrls
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.Update(id, url); err != nil {
		if err == store.NotFoundErr {
			NotFoundHandler(w, r)
			return
		}

		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(url)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h UrlHandler) DeleteUrl(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.store.Remove(id); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
