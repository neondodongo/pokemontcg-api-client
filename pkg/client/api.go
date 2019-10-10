package client

import (
	"encoding/json"
	"net/http"
)

func Health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body HealthResponse
		body.Status = http.StatusOK
		body.Message = "Application is up and running"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(body)
	})
}
