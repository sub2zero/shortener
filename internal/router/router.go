package remote

import (
	"mime"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func matchJSON(r *http.Request, rm *mux.RouteMatch) bool {
	ctHdr := r.Header.Get("Content-Type")
	contentType, params, err := mime.ParseMediaType(ctHdr)
	if err != nil {
		return false
	}
	charset := params["charset"]
	return contentType == "application/json" && (charset == "" || strings.EqualFold(charset, "UTF-8"))
}
