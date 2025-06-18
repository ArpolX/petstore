package repository

import (
	"petstore/internal/logs"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PetRepository struct {
	db  *gorm.DB
	Log logs.Logger
}

type PetRepositoryer interface {
	Create(p Pet) error
	Update(p Pet) error
	UpdateNameStatus(petId int, name, status string) error
	Delete(petId int) error
	GetId(petId int) (Pet, error)
	GetStatus(status string) ([]Pet, error)
	GetName(name string) (Pet, error)
}

func NewPetRepository(logger logs.Logger, db *gorm.DB) PetRepositoryer {
	return &PetRepository{
		db:  db,
		Log: logger,
	}
}

func (pe *PetRepository) Create(p Pet) error {
	err := pe.db.Create(&p).Error
	if err != nil {
		pe.Log.Error("Ошибка в Create запросе", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func (pe *PetRepository) Update(p Pet) error {
	err := pe.db.Model(&Pet{}).Where("id = ?", p.Id).Updates(p).Error
	if err != nil {
		pe.Log.Error("Ошибка в Update запросе", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func (pe *PetRepository) UpdateNameStatus(petId int, name, status string) error {
	err := pe.db.Model(&Pet{}).Where("id = ?", petId).Updates(map[string]interface{}{
		"name":   name,
		"status": status,
	}).Error
	if err != nil {
		return err
	}

	return nil
}

func (pe *PetRepository) Delete(petId int) error {
	err := pe.db.Delete(&Pet{}, "id = ?", petId).Error
	if err != nil {
		return err
	}

	return nil
}

func (pe *PetRepository) GetId(petId int) (Pet, error) {
	p := Pet{}

	err := pe.db.First(&p, "id = ?", petId).Error
	if err != nil {
		return Pet{}, err
	}

	return p, nil
}

func (pe *PetRepository) GetStatus(status string) ([]Pet, error) {
	p := []Pet{}

	err := pe.db.Find(&p, "status = ?", status).Error
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (pe *PetRepository) GetName(name string) (Pet, error) {
	p := Pet{}

	err := pe.db.First(&p, "name = ?", name).Error
	if err != nil {
		return Pet{}, err
	}

	return p, nil
}
