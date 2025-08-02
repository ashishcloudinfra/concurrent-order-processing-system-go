package order

import (
	"github.com/gorilla/mux"
)

func RegisterOrderRoutes(mux *mux.Router) {
	mux.HandleFunc("/orders", CreateOrderHandler).Methods("POST")
	mux.HandleFunc("/orders/{id:[a-zA-Z0-9-_]+}", GetOrderByIDHandler).Methods("GET")
	mux.HandleFunc("/orders", GetOrdersHandler).Methods("GET")
	mux.HandleFunc("/orders/{id:[a-zA-Z0-9-_]+}", UpdateOrderStatusHandler).Methods("PUT")
	mux.HandleFunc("/orders/{id:[a-zA-Z0-9-_]+}", DeleteOrderHandler).Methods("DELETE")
	mux.HandleFunc("/orders/status/{id:[a-zA-Z0-9-_]+}", GetOrderStatusHandler).Methods("GET")
	mux.HandleFunc("/orders/cancel/{id:[a-zA-Z0-9-_]+}", CancelOrderHandler).Methods("POST")
	mux.HandleFunc("/orders/process", ProcessOrdersHandler).Methods("POST")
}
