package order

import (
	"concurrent-order-processing-system/utils/responses"
	"net/http"

	"github.com/gorilla/mux"
)

func DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	order, exists := Store.Get(orderID)
	if !exists {
		responses.WriteErrorResponse(w, "Order not found", http.StatusNotFound)
		return
	}

	if order.Status != "pending" {
		responses.WriteErrorResponse(w, "Cannot delete order, it's not in 'pending' status", http.StatusBadRequest)
		return
	}

	Store.Delete(orderID)

	responses.WriteSuccessResponse(w, "Order deleted successfully", nil)
}
