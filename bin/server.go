package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	handler "shortener/internal/handler"
	store "shortener/internal/struct"

	"github.com/gorilla/mux"
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
	UrlHandler := handler.NewUrlHandler(store)
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
