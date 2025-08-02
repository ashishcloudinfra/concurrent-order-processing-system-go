package health

import (
	"concurrent-order-processing-system/utils/responses"
	"net/http"

	"github.com/gorilla/mux"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	responses.WriteSuccessResponse(w, "Service is running", nil)
}

func RegisterHealthRoutes(mux *mux.Router) {
	mux.HandleFunc("/health", HealthCheckHandler).Methods("GET")
}
