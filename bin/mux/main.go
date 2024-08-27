package main

//need to implement the same as for gin in terms of ids
import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	store "shortener/internal/struct"

	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
)

type serverAddress string

const serverAdd serverAddress = "127.0.0.1"

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.Method + " requestURI: " + r.RequestURI + " from: " + r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
func main() {
	store := store.NewUrlStore()
	UrlHandler := NewUrlHandler(store)
	home := homeHandler{}

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	// Register the routes
	r.HandleFunc("/", home.ServeHTTP)
	r.HandleFunc("/url", UrlHandler.ListUrls).Methods("GET")
	r.HandleFunc("/url", UrlHandler.CreateUrl).Methods("POST")
	r.HandleFunc("/url/{id}", UrlHandler.GetUrl).Methods("GET")
	r.HandleFunc("/url/{id}", UrlHandler.UpdateUrl).Methods("PUT")
	r.HandleFunc("/url/{id}", UrlHandler.DeleteUrl).Methods("DELETE")

	ctx, cancelCtx := context.WithCancel(context.Background())
	serverOne := &http.Server{
		Addr:    ":8080",
		Handler: r,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, serverAdd, l.Addr().String())
			return ctx
		},
	}
	defer cancelCtx()
	go func() {
		err := serverOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("server one closed")
		} else if err != nil {
			log.Println("error listening for server one: ", err)
		}
		cancelCtx()
	}()
	<-ctx.Done()
}

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
