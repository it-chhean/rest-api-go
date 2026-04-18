package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"email-api/handlers"
)

func RegisterRoutes(h *handlers.Handler) *mux.Router {
	r := mux.NewRouter()
	r.Use(jsonContentTypeMiddleware)

	v1 := r.PathPrefix("/api/v1").Subrouter()
	registerEmailRoutes(v1, h)

	return r
}

func registerEmailRoutes(r *mux.Router, h *handlers.Handler) {
	r.HandleFunc("/emails",      h.GetEmails).Methods(http.MethodGet)
	r.HandleFunc("/emails/{id}", h.GetEmailByID).Methods(http.MethodGet)
	r.HandleFunc("/emails",      h.CreateEmail).Methods(http.MethodPost)
	r.HandleFunc("/emails/{id}", h.UpdateEmail).Methods(http.MethodPut)
	r.HandleFunc("/emails/{id}", h.DeleteEmail).Methods(http.MethodDelete)
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
