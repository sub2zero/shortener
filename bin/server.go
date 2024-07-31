package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	handler "shortener/internal/handler"
	router "shortener/internal/router"
	store "shortener/internal/struct"
)

type serverAddress string

const serverAdd serverAddress = "127.0.0.1"

func main() {
	mStore := new(store.UrlStore)
	mStore.Create()
	eh := handler.NewUrlHandler(*mStore)

	r := router.NewRouter(eh)

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
