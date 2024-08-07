package handler

import (
	"log"
	"net/http"
	store "shortener/internal/struct"
	"strconv"

	"github.com/gorilla/mux"
)

type UrlHandler struct {
	urlService store.UrlStore
}

func NewUrlHandler(ms store.UrlStore) UrlHandler {
	return UrlHandler{
		urlService: ms,
	}
}
func (mh UrlHandler) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	// typeParam := vars["type"]
	nameParam := vars["url"]
	// valueParam := vars["shorturl"]

	_, err := w.Write([]byte(nameParam))
	if err != nil {
		log.Println("Failed writing HTTP response")
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (mh UrlHandler) PostUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	// typeParam := vars["type"]
	nameParam := vars["url"]
	valueParam := vars["shorturl"]

	// if validateURLType(typeParam) {
	// 	http.Error(w, "Not implemented", http.StatusNotImplemented)
	// 	return
	// }

	if _, err := strconv.ParseFloat(valueParam, 64); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// if typeParam == config.Counter {
	// 	if _, err := strconv.Atoi(valueParam); err != nil {
	// 		http.Error(w, "Bad request", http.StatusBadRequest)
	// 		return
	// 	}
	// }

	if err := mh.urlService.AddUrl(nameParam, valueParam); err != nil {
		http.Error(w, "See other request", http.StatusSeeOther)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (mh UrlHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ms := mh.urlService.GetAll()
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println(r.Method)
		return
	}
	_, err := w.Write([]byte(ms))
	if err != nil {
		log.Println("Failed writing HTTP response")
		return
	}

}

// func validateURLType(typeParam string) bool {
// 	return typeParam != config.URL
// }
