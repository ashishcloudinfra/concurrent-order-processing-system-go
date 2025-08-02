package order

import (
	"concurrent-order-processing-system/utils/responses"
	"net/http"

	"github.com/gorilla/mux"
)

func GetOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	order, exists := Store.Get(orderID)
	if !exists {
		responses.WriteErrorResponse(w, "Order not found", http.StatusNotFound)
		return
	}

	statusResponse := struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}{
		ID:     order.ID,
		Status: order.Status,
	}

	responses.WriteSuccessResponse(w, "Order status fetched successfully", statusResponse)
}
