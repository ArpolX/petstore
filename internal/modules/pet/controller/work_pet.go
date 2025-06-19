package controller

import (
	"fmt"
	"net/http"
	"petstore/internal/logs"
	"petstore/internal/modules/pet/service"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type AnimalStore struct {
	Log        logs.Logger
	PetService service.PetServicer
}

type AnimalStorer interface {
	RegisterPet(w http.ResponseWriter, r *http.Request)
	UpdatePet(w http.ResponseWriter, r *http.Request)
	UpdateNameStatusPet(w http.ResponseWriter, r *http.Request)
	GetPetByStatus(w http.ResponseWriter, r *http.Request)
	GetPet(w http.ResponseWriter, r *http.Request)
	GetImagePet(w http.ResponseWriter, r *http.Request)
	DeletePet(w http.ResponseWriter, r *http.Request)
}

func NewAnimalStore(logger logs.Logger, petService service.PetServicer) AnimalStorer {
	return &AnimalStore{
		Log:        logger,
		PetService: petService,
	}
}

func (a *AnimalStore) RegisterPet(w http.ResponseWriter, r *http.Request) {
	p := Pet{}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	if err := IsValidStruct(p); err != nil {
		strErr := fmt.Sprintf("Невалидный запрос (неверный формат): %v", err.Error())
		http.Error(w, strErr, http.StatusBadRequest)
		return
	}

	tagService := []service.Tag{}
	for _, tag := range p.Tag {
		tags := service.Tag{
			Id: tag.Id,
		}
		tagService = append(tagService, tags)
	}

	petService := service.Pet{
		Id:        p.Id,
		Name:      p.Name,
		Category:  service.Category(p.Category),
		PhotoUrls: p.PhotoUrls,
		Tag:       tagService,
		Status:    p.Status,
	}

	resp, err := a.PetService.RegisterService(petService)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

func (a *AnimalStore) UpdatePet(w http.ResponseWriter, r *http.Request) {
	p := Pet{}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	if err := IsValidStruct(p); err != nil {
		strErr := fmt.Sprintf("Невалидный запрос (неверный формат): %v", err.Error())
		http.Error(w, strErr, http.StatusBadRequest)
		return
	}

	tagService := []service.Tag{}
	for _, tag := range p.Tag {
		tags := service.Tag{
			Id: tag.Id,
		}
		tagService = append(tagService, tags)
	}

	petService := service.Pet{
		Id:        p.Id,
		Name:      p.Name,
		Category:  service.Category(p.Category),
		PhotoUrls: p.PhotoUrls,
		Tag:       tagService,
		Status:    p.Status,
	}

	resp, err := a.PetService.UpdateService(petService)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

func (a *AnimalStore) UpdateNameStatusPet(w http.ResponseWriter, r *http.Request) {
	petId := chi.URLParam(r, "petId")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Неверный формат formData", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	status := r.FormValue("status")
	petIdInt, _ := strconv.Atoi(petId)

	resp, err := a.PetService.UpdateNameStatusService(petIdInt, name, status)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

func (a *AnimalStore) GetPetByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	pet, err := a.PetService.GetByStatusService(status)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(pet)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}
}

func (a *AnimalStore) GetPet(w http.ResponseWriter, r *http.Request) {
	petId := chi.URLParam(r, "petId")
	petIdInt, _ := strconv.Atoi(petId)

	pet, err := a.PetService.GetService(petIdInt)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(pet)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}
}

func (a *AnimalStore) DeletePet(w http.ResponseWriter, r *http.Request) {
	petId := chi.URLParam(r, "petId")
	petIdInt, _ := strconv.Atoi(petId)

	resp, err := a.PetService.DeleteService(petIdInt)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

func (a *AnimalStore) GetImagePet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("В разработке"))
}

func IsValidStruct(u Pet) error {
	valid := validator.New()
	return (valid.Struct(u))
}
