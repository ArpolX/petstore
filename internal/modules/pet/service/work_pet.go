package service

import (
	"petstore/internal/logs"
	"petstore/internal/modules/pet/repository"
)

type PetService struct {
	Log           logs.Logger
	PetRepository repository.PetRepositoryer
}

type PetServicer interface {
	RegisterService(p Pet) (string, error)
	UpdateService(p Pet) (string, error)
	UpdateNameStatusService(petId int, name, status string) (string, error)
	GetByStatusService(status string) ([]repository.Pet, error)
	GetService(petId int) (repository.Pet, error)
	DeleteService(petId int) (string, error)
}

func NewPetService(logger logs.Logger, petRepository repository.PetRepositoryer) PetServicer {
	return &PetService{
		Log:           logger,
		PetRepository: petRepository,
	}
}

func (pe *PetService) RegisterService(p Pet) (string, error) {
	pet, err := pe.PetRepository.GetId(p.Id)
	if err != nil {
		return "", err
	}

	if pet.Id != 0 {
		return "Такой id уже существует", nil
	}

	tagRepository := []repository.Tag{}
	for _, tag := range p.Tag {
		tags := repository.Tag{
			Id: tag.Id,
		}
		tagRepository = append(tagRepository, tags)
	}

	petRepository := repository.Pet{
		Id:        p.Id,
		Name:      p.Name,
		Category:  repository.Category(p.Category),
		PhotoUrls: p.PhotoUrls,
		Tag:       tagRepository,
		Status:    p.Status,
	}

	err = pe.PetRepository.Create(petRepository)
	if err != nil {
		return "", err
	}

	return "Регистрация животного прошла успешно", nil
}

func (pe *PetService) UpdateService(p Pet) (string, error) {
	pet, err := pe.PetRepository.GetId(p.Id)
	if err != nil {
		return "", err
	}

	if pet.Id == 0 {
		return "Такого id не существует", nil
	}

	tagRepository := []repository.Tag{}
	for _, tag := range p.Tag {
		tags := repository.Tag{
			Id: tag.Id,
		}
		tagRepository = append(tagRepository, tags)
	}

	petRepository := repository.Pet{
		Id:        p.Id,
		Name:      p.Name,
		Category:  repository.Category(p.Category),
		PhotoUrls: p.PhotoUrls,
		Tag:       tagRepository,
		Status:    p.Status,
	}

	err = pe.PetRepository.Update(petRepository)
	if err != nil {
		return "", err
	}

	return "Обновление информации о животном прошло успешно", nil
}

func (pe *PetService) UpdateNameStatusService(petId int, name, status string) (string, error) {
	pet, err := pe.PetRepository.GetId(petId)
	if err != nil {
		return "", err
	}

	if pet.Id == 0 {
		return "ID животного не найдено", nil
	}

	err = pe.PetRepository.UpdateNameStatus(petId, name, status)
	if err != nil {
		return "", err
	}

	return "Имя и статус обновлены", nil
}

func (pe *PetService) GetByStatusService(status string) ([]repository.Pet, error) {
	pet, err := pe.PetRepository.GetStatus(status)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

func (pe *PetService) GetService(petId int) (repository.Pet, error) {
	pet, err := pe.PetRepository.GetId(petId)
	if err != nil {
		return repository.Pet{}, err
	}

	return pet, nil
}

func (pe *PetService) DeleteService(petId int) (string, error) {
	pet, err := pe.PetRepository.GetId(petId)
	if err != nil {
		return "", err
	}

	if pet.Name == "" {
		return "Животное не найдено", nil
	}

	err = pe.PetRepository.Delete(petId)
	if err != nil {
		return "", err
	}

	return "Животное успешно удалено", nil
}
