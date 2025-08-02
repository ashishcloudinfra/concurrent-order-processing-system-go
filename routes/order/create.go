package order

import (
	"concurrent-order-processing-system/utils/responses"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var newOrder Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		responses.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newOrder.ID = uuid.New().String()
	newOrder.Timestamp = time.Now()
	newOrder.Status = "pending"

	if newOrder.Customer == "" {
		responses.WriteErrorResponse(w, "Customer name is required", http.StatusBadRequest)
		return
	}

	Store.Set(newOrder)

	responses.WriteSuccessResponseWithCode(w, "Order created successfully", &newOrder, http.StatusCreated)
}
