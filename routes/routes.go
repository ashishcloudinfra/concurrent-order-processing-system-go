package routes

import (
	"concurrent-order-processing-system/routes/health"
	"concurrent-order-processing-system/routes/order"

	"github.com/gorilla/mux"
)

func InitializeRoutes() *mux.Router {
	mux := mux.NewRouter()

	order.RegisterOrderRoutes(mux)
	health.RegisterHealthRoutes(mux)

	return mux
}
