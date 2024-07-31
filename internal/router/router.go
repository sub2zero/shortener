package remote

import (
	"log"
	"net/http"
	handler "shortener/internal/handler"

	"github.com/gorilla/mux"
)

func NewRouter(mh handler.UrlHandler) *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/", mh.GetAll).
		Methods(http.MethodGet)

	// r.HandleFunc("/value/{type}/{name}", mh.GetByName).
	// 	Methods(http.MethodGet)
	// r.HandleFunc("/value/", mh.GetByNameJSON).MatcherFunc(matchJSON).
	// 	Methods(http.MethodPost)

	// r.HandleFunc("/update/", mh.JSONUpdate).MatcherFunc(matchJSON).
	// 	Methods(http.MethodPost)
	// r.HandleFunc("/update/{type}/{name}/{value}", mh.PostUpdate).
	// 	Methods(http.MethodPost)
	return r
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.Method + " requestURI: " + r.RequestURI + " from: " + r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// func matchJSON(r *http.Request, rm *mux.RouteMatch) bool {
// 	ctHdr := r.Header.Get("Content-Type")
// 	contentType, params, err := mime.ParseMediaType(ctHdr)
// 	if err != nil {
// 		return false
// 	}
// 	charset := params["charset"]
// 	return contentType == "application/json" && (charset == "" || strings.EqualFold(charset, "UTF-8"))
// }
