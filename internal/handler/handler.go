package handler

import (
	"log"
	"net/http"
	store "shortener/internal/struct"
)

type UrlHandler struct {
	urlService store.UrlStore
}

func NewUrlHandler(ms store.UrlStore) UrlHandler {
	return UrlHandler{
		urlService: ms,
	}
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
