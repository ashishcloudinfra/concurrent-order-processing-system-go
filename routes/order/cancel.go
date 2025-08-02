package order

import (
	"concurrent-order-processing-system/utils/responses"
	"net/http"

	"github.com/gorilla/mux"
)

func CancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	order, exists := Store.Get(orderID)
	if !exists {
		responses.WriteErrorResponse(w, "Order not found", http.StatusNotFound)
		return
	}

	if order.Status != "pending" {
		responses.WriteErrorResponse(w, "Cannot cancel order, it's not in 'pending' status", http.StatusBadRequest)
		return
	}

	order.Status = "canceled"
	Store.Set(order)

	responses.WriteSuccessResponse(w, "Order canceled successfully", nil)
}
