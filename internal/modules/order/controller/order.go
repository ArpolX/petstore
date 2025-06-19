package controller

import (
	"net/http"
	"petstore/internal/logs"
	"petstore/internal/modules/order/service"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type OrderRespond struct {
	Log           logs.Logger
	OrderServicer service.OrderServicer
}

type OrderResponder interface {
	PlaceOrder(w http.ResponseWriter, r *http.Request)
	GetOrder(w http.ResponseWriter, r *http.Request)
	DeleteOrder(w http.ResponseWriter, r *http.Request)
}

func NewOrderRespond(logger logs.Logger, orderServicer service.OrderServicer) OrderResponder {
	return &OrderRespond{
		Log:           logger,
		OrderServicer: orderServicer,
	}
}

func (o *OrderRespond) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	order := Order{}

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	if err = IsValidStruct(order); err != nil {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	orderService := service.Order{
		Id:       order.Id,
		PetId:    order.PetId,
		Quantity: order.Quantity,
		ShipDate: order.ShipDate,
		Status:   order.Status,
		Complete: order.Complete,
	}

	resp, err := o.OrderServicer.PlaceOrderService(orderService)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

func (o *OrderRespond) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderId := chi.URLParam(r, "orderId")

	if orderId == "" {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	orderIdInt, _ := strconv.Atoi(orderId)

	resp, order, err := o.OrderServicer.GetOrderService(orderIdInt)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	if resp != "" {
		w.Write([]byte(resp))
		return
	}

	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}
}

func (o *OrderRespond) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderId := chi.URLParam(r, "orderId")

	if orderId == "" {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	orderIdInt, _ := strconv.Atoi(orderId)

	resp, err := o.OrderServicer.DeleteOrderService(orderIdInt)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

func IsValidStruct(o Order) error {
	valid := validator.New()
	return valid.Struct(o)
}
