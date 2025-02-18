package route

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Load() *mux.Router {
	r := mux.NewRouter()
	v1 := r.PathPrefix("/v1").Subrouter()

	loadV1Routes(v1)


	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error": "404 Page Not Found", "message": "The requested resource could not be found"}`)
	})

	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	})

	return r
}