package repository

import (
	"fmt"
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
	UpdatePhoto(petId int, filepath string) error
	Delete(petId int) error
	GetId(petId int) (Pet, error)
	GetStatus(status string) ([]Pet, error)
}

func NewPetRepository(logger logs.Logger, db *gorm.DB) PetRepositoryer {
	return &PetRepository{
		db:  db,
		Log: logger,
	}
}

func (pe *PetRepository) Create(p Pet) error {
	err := pe.db.Transaction(func(tx *gorm.DB) error {
		err := pe.db.Exec("insert into pets (id, name, category_id, status) values (?, ?, ?, ?)", p.Id, p.Name, p.Category.Id, p.Status).Error
		if err != nil {
			pe.Log.Error("Ошибка в Create запросе", zap.String("err", err.Error()))
			return err
		}

		for _, t := range p.Tag {
			err = tx.Exec("INSERT INTO pet_tags (pet_id, tag_id) VALUES (?, ?)", p.Id, t.Id).Error
			if err != nil {
				pe.Log.Error("Ошибка в Exec запросе", zap.String("err", err.Error()))
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (pe *PetRepository) Update(p Pet) error {
	err := pe.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Pet{}).Where("id = ?", p.Id).Updates(map[string]interface{}{
			"name":        p.Name,
			"status":      p.Status,
			"category_id": p.Category.Id,
		}).Error
		if err != nil {
			pe.Log.Error("Ошибка в Update запросе", zap.String("err", err.Error()))
			return err
		}

		err = tx.Exec("DELETE FROM pet_tags WHERE pet_id = ?", p.Id).Error
		if err != nil {
			pe.Log.Error("Ошибка в Exec запросе", zap.String("err", err.Error()))
			return err
		}

		for _, t := range p.Tag {
			err = tx.Exec("INSERT INTO pet_tags (pet_id, tag_id) VALUES (?, ?)", p.Id, t.Id).Error
			if err != nil {
				pe.Log.Error("Ошибка в Exec запросе", zap.String("err", err.Error()))
				return err
			}
		}

		return nil
	})
	if err != nil {
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
		pe.Log.Error("Ошибка в Update запросе", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func (pe *PetRepository) UpdatePhoto(petId int, filepath string) error {
	err := pe.db.Model(&Pet{}).Where("id = ?", petId).Updates(map[string]interface{}{
		"photo_url": filepath,
	}).Error
	if err != nil {
		pe.Log.Error("Ошибка в Update запросе", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func (pe *PetRepository) Delete(petId int) error {
	err := pe.db.Transaction(func(tx *gorm.DB) error {
		err := pe.db.Exec("delete from pet_tags where pet_id = ?", petId).Error
		if err != nil {
			pe.Log.Error("Ошибка в Exec запросе", zap.String("err", err.Error()))
			return err
		}

		err = pe.db.Delete(&Pet{}, "id = ?", petId).Error
		if err != nil {
			pe.Log.Error("Ошибка в Delete запросе", zap.String("err", err.Error()))
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (pe *PetRepository) GetId(petId int) (Pet, error) {
	p := Pet{}

	err := pe.db.Preload("Tag").First(&p, "id = ?", petId).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return Pet{}, err
	}
	fmt.Println(p)

	return p, nil
}

func (pe *PetRepository) GetStatus(status string) ([]Pet, error) {
	p := []Pet{}

	err := pe.db.Find(&p, "status = ?", status).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return p, nil
}
