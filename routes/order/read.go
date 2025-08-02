package order

import (
	"concurrent-order-processing-system/utils/responses"
	"net/http"
	"time"
)

func GetOrderByIDHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Path[len("/orders/"):]

	order, exists := Store.Get(orderID)
	if !exists {
		responses.WriteErrorResponse(w, "Order not found", http.StatusNotFound)
		return
	}

	responses.WriteSuccessResponse(w, "Order fetched successfully", order)
}

func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var fromTime, toTime time.Time
	var err error

	if fromStr != "" {
		fromTime, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			responses.WriteErrorResponse(w, "Invalid 'from' timestamp", http.StatusBadRequest)
			return
		}
	}

	if toStr != "" {
		toTime, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			responses.WriteErrorResponse(w, "Invalid 'to' timestamp", http.StatusBadRequest)
			return
		}
	}

	allOrders := Store.GetAll()
	var orderList []Order

	for _, order := range allOrders {
		if (fromStr == "" || order.Timestamp.After(fromTime)) &&
			(toStr == "" || order.Timestamp.Before(toTime)) {
			orderList = append(orderList, order)
		}
	}

	responses.WriteSuccessResponse(w, "Orders fetched successfully", orderList)
}
