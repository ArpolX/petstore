package repository

import (
	"petstore/internal/logs"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderRepository struct {
	Db  *gorm.DB
	Log logs.Logger
}

type OrderRepositorer interface {
	Create(o Order) error
	Get(orderId int) (Order, error)
	Delete(orderId int) error
}

func NewOrderRepository(logger logs.Logger, db *gorm.DB) OrderRepositorer {
	return &OrderRepository{
		Db:  db,
		Log: logger,
	}
}

func (or *OrderRepository) Create(o Order) error {
	err := or.Db.Create(&o).Error
	if err != nil {
		or.Log.Error("Ошибка в Create запросе", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func (or *OrderRepository) Get(orderId int) (Order, error) {
	o := Order{}

	err := or.Db.Find(&o, "id = ?", orderId).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return Order{}, err
	}

	return o, nil
}

func (or *OrderRepository) Delete(orderId int) error {
	err := or.Db.Delete(&Order{}, "id = ?", orderId).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	return nil
}
