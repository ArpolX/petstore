package controller

import (
	"fmt"
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

// @Summary Разместить заказ на животное
// @Description Создание заказа на животное. Отсчёт по id (не должен повторяться). Нулевые значения у полей integer не допускаются
// @Tags order
// @Accept json
// @Produce plain
// @Param order body Order true "Заполни все поля для размещения заказа"
// @Success 200 {string} Info "Успешное размещение заказа или не ошибочное сообщение"
// @Failure 400 {string} Err "Неверной формат"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /store/order [post]
func (o *OrderRespond) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	order := Order{}

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		fmt.Println(err)
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

// @Summary Получить заказ по Id
// @Description Получить заказ по Id
// @Tags order
// @Accept json
// @Produce json
// @Param orderId path string true "Введите id заказа"
// @Success 200 {object} Order "Структура заказа или не ошибочное сообщение"
// @Failure 400 {string} Err "Неверной формат структуры"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /store/order/{orderId} [get]
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

// @Summary Удалить заказ по Id
// @Description Удалить заказ по Id
// @Tags order
// @Accept json
// @Produce plain
// @Param orderId path string true "Введите id заказа"
// @Success 200 {string} Info "Успешное удаление заказа или не ошибочное сообщение"
// @Failure 400 {string} Err "Неверной формат структуры"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /store/order/{orderId} [delete]
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
