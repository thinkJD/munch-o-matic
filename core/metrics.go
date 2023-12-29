package core

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var dishOrders = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "total_orders_per_dish",
	Help: "Total number of orders per dish"},
	[]string{"order_id", "dish_name"},
)

var accountBalance = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "account_balance",
	Help: "Account balance in cent"},
	[]string{"user_id", "user_name"},
)

var paymentsTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "payments_total",
	Help: "Tracks total payments"},
	[]string{"user_id", "user_name"},
)

func init() {
	// Register metrics with Prometheus's default registry.
	prometheus.MustRegister(dishOrders)
	prometheus.MustRegister(accountBalance)
	prometheus.MustRegister(paymentsTotal)
}

func UpdateOrdersPlaced(OrderId int, DishName string, NumOrders int) {
	dishOrders.With(prometheus.Labels{"order_id": fmt.Sprint(OrderId), "dish_name": DishName}).Set(float64(NumOrders))
}

func UpdateAccountBalance(UserId int, UserName string, Balance int) {
	accountBalance.With(prometheus.Labels{"user_id": fmt.Sprint(UserId), "user_name": UserName}).Set(float64(Balance))
}

func UpdatePaymentsTotal(UserId int, UserName string, Value int) {
	paymentsTotal.With(prometheus.Labels{"user_id": fmt.Sprint(UserId), "user_name": UserName}).Set(float64(Value))
}

// GetPrometheusHandler returns the HTTP handler for Prometheus metrics.
func GetPrometheusHandler() http.Handler {
	return promhttp.Handler()
}
