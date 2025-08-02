package order

import (
	"concurrent-order-processing-system/utils/responses"
	"context"
	"net/http"
	"time"
)

func ProcessOrdersHandler(w http.ResponseWriter, r *http.Request) {
	pendingOrders := Store.GetPendingOrders()

	if len(pendingOrders) == 0 {
		responses.WriteErrorResponse(w, "No pending orders to process", http.StatusNotFound)
		return
	}

	for _, order := range pendingOrders {
		go processOrderWithTimeout(order)
	}

	responses.WriteSuccessResponse(w, "Processing started", nil)
}

func processOrderWithTimeout(order Order) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-time.After(3 * time.Second):
		if ctx.Err() == context.DeadlineExceeded {
			order.Status = "canceled"
		} else {
			order.Status = "shipped"
		}
		Store.Set(order)
	case <-ctx.Done():
		order.Status = "canceled"
		Store.Set(order)
	}
}
