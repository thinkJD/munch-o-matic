package core

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ordersPlaced is a Prometheus counter metric.
var dishOrders = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "total_orders_per_dish",
	Help: "Total number of orders per dish orderID"},
	[]string{"OrderId", "DishName"},
)

func init() {
	// Register metrics with Prometheus's default registry.
	prometheus.MustRegister(dishOrders)
}

func UpdateOrdersPlaced(OrderId int, DishName string, NumOrders int) {
	dishOrders.With(prometheus.Labels{"dish": DishName}).Set(float64(NumOrders))
}

// GetPrometheusHandler returns the HTTP handler for Prometheus metrics.
func GetPrometheusHandler() http.Handler {
	return promhttp.Handler()
}
