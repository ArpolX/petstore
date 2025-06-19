package service

import (
	"petstore/internal/logs"
	"petstore/internal/modules/order/repository"
)

type OrderService struct {
	Log          logs.Logger
	OrderReposit repository.OrderRepositorer
}

type OrderServicer interface {
	PlaceOrderService(o Order) (string, error)
	GetOrderService(orderId int) (string, repository.Order, error)
	DeleteOrderService(orderId int) (string, error)
}

func NewOrderService(logger logs.Logger, orderReposit repository.OrderRepositorer) OrderServicer {
	return &OrderService{
		Log:          logger,
		OrderReposit: orderReposit,
	}
}

func (os *OrderService) PlaceOrderService(o Order) (string, error) {
	order, err := os.OrderReposit.Get(o.Id)
	if err != nil {
		return "", err
	}

	if order.PetId != 0 {
		return "Такой id заказа уже существует", nil
	}

	orderRep := repository.Order{
		Id:       o.Id,
		PetId:    o.PetId,
		Quantity: o.Quantity,
		ShipDate: o.ShipDate,
		Status:   o.Status,
		Complete: o.Complete,
	}

	err = os.OrderReposit.Create(orderRep)
	if err != nil {
		return "", err
	}

	return "Заказ сформирован", nil
}

func (os *OrderService) GetOrderService(orderId int) (string, repository.Order, error) {
	order, err := os.OrderReposit.Get(orderId)
	if err != nil {
		return "", repository.Order{}, err
	}

	if order.PetId == 0 {
		return "Такого Id заказа нет", repository.Order{}, nil
	}

	return "", order, nil
}

func (os *OrderService) DeleteOrderService(orderId int) (string, error) {
	order, err := os.OrderReposit.Get(orderId)
	if err != nil {
		return "", err
	}

	if order.PetId == 0 {
		return "Такого Id заказа нет", nil
	}

	err = os.OrderReposit.Delete(orderId)
	if err != nil {
		return "", err
	}

	return "Заказ удалён", nil
}
