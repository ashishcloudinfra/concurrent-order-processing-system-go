package order

import (
	"concurrent-order-processing-system/utils/responses"
	"encoding/json"
	"net/http"
)

func UpdateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Path[len("/orders/"):]

	var updatedOrder Order
	err := json.NewDecoder(r.Body).Decode(&updatedOrder)
	if err != nil {
		responses.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updatedOrder.Status == "" {
		responses.WriteErrorResponse(w, "Status cannot be empty", http.StatusBadRequest)
		return
	}

	order, exists := Store.Get(orderID)
	if !exists {
		responses.WriteErrorResponse(w, "Order not found", http.StatusNotFound)
		return
	}

	order.Status = updatedOrder.Status
	Store.Set(order)

	responses.WriteSuccessResponse(w, "Order status updated successfully", order)
}
