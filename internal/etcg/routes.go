package etcg

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handle maps the API URIs with the application controller
func Handle(r *mux.Router, c Controller) error {
	r.Handle("/cards/get", c.GetCards()).Methods(http.MethodGet)
	// r.Handle("/sets/get", c.GetSets()).Methods(http.MethodGet)

	return nil
}
