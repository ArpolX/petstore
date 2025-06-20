package controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	AddPhotoPet(w http.ResponseWriter, r *http.Request)
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

// @Summary Добавить нового питомца в магазин
// @Description Создание и добавление нового питомца с различными полями. Отсчёт будет идти от id (не должен повторяться). В категориях и тегах указывается только id. Категории: 1 - dog, 2 - cat. Tags: 1 - friendly, 2 - wild, 3 - trained. Также доступно 3 статуса: available, sold и pending
// @Tags pet
// @Security BearerAuth
// @Accept json
// @Produce plain
// @Param pet body Pet true "Заполни все поля для добавления"
// @Success 200 {string} Info "Успешное добавление или не ошибочное сообщение"
// @Failure 401 {string} Err "Аутентификация не пройдена"
// @Failure 403 {string} Err "Авторизация не пройдена"
// @Failure 400 {string} Err "Неверной формат"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /pet/ [post]
func (a *AnimalStore) RegisterPet(w http.ResponseWriter, r *http.Request) {
	p := Pet{}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	StatusMap := map[string]struct{}{
		"available": {},
		"sold":      {},
		"pending":   {},
	}

	if _, ok := StatusMap[p.Status]; !ok {
		http.Error(w, "Невалидный статус", http.StatusBadRequest)
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
		Id:       p.Id,
		Name:     p.Name,
		Category: service.Category(p.Category),
		Tag:      tagService,
		Status:   p.Status,
	}

	resp, err := a.PetService.RegisterService(petService)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

// @Summary Добавить фото питомца
// @Description Добавить фото конкретному питомцу, будет сохранено локально, адрес можете узнать в поле photoUrl
// @Tags pet
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce plain
// @Param petId path string true "Введи id питомца"
// @Param photoFile formData file true "Добавь изображение животного"
// @Success 200 {string} Info "Успешное добавление фото или не ошибочное сообщение"
// @Failure 401 {string} Err "Аутентификация не пройдена"
// @Failure 403 {string} Err "Авторизация не пройдена"
// @Failure 400 {string} Err "Неверной формат"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /pet/photo/{petId} [post]
func (a *AnimalStore) AddPhotoPet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		fmt.Println("111", err)
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	petId := chi.URLParam(r, "petId")
	petIdInt, _ := strconv.Atoi(petId)

	file, handler, err := r.FormFile("photoFile")
	if err != nil {
		fmt.Println("222", err)
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	projectRoot := filepath.Dir(cwd)
	filePath := filepath.Join(projectRoot, "internal", "uploads", fmt.Sprintf("%v_%v", petId, handler.Filename))

	dst, err := os.Create(filePath)
	if err != nil {
		fmt.Println("333", err)
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	resp, err := a.PetService.AddPhotoPet(petIdInt, filePath)
	if err != nil {
		fmt.Println("444", err)
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

// @Summary Обновить информацию о питомце
// @Description Обновить информацию, id и фото изменить нельзя. Категории: 1 - dog, 2 - cat. Tags: 1 - friendly, 2 - wild, 3 - trained
// @Tags pet
// @Security BearerAuth
// @Accept json
// @Produce plain
// @Param pet body Pet true "Заполни все поля для изменения"
// @Success 200 {string} Info "Успешное обновление информации или не ошибочное сообщение"
// @Failure 400 {string} Err "Неверной формат"
// @Failure 401 {string} Err "Аутентификация не пройдена"
// @Failure 403 {string} Err "Авторизация не пройдена"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /pet/ [put]
func (a *AnimalStore) UpdatePet(w http.ResponseWriter, r *http.Request) {
	p := Pet{}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	StatusMap := map[string]struct{}{
		"available": {},
		"sold":      {},
		"pending":   {},
	}

	if _, ok := StatusMap[p.Status]; !ok {
		http.Error(w, "Невалидный статус", http.StatusBadRequest)
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
		Id:       p.Id,
		Name:     p.Name,
		Category: service.Category(p.Category),
		Tag:      tagService,
		Status:   p.Status,
	}

	resp, err := a.PetService.UpdateService(petService)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

// @Summary Обновить информацию о питомце
// @Description Обновить только name и status.
// @Tags pet
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce plain
// @Param petId path string true "Введи Id животного"
// @Param name formData string true "Имя питомца"
// @Param status formData string true "Статус питомца (available, sold, pending)"
// @Success 200 {string} Info "Успешное обновление информации или не ошибочное сообщение"
// @Failure 400 {string} Err "Неверной формат"
// @Failure 401 {string} Err "Аутентификация не пройдена"
// @Failure 403 {string} Err "Авторизация не пройдена"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /pet/{petId} [post]
func (a *AnimalStore) UpdateNameStatusPet(w http.ResponseWriter, r *http.Request) {
	petId := chi.URLParam(r, "petId")
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		http.Error(w, "Неверный формат formData", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	status := r.FormValue("status")
	if name == "" || status == "" {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}
	petIdInt, _ := strconv.Atoi(petId)

	StatusMap := map[string]struct{}{
		"available": {},
		"sold":      {},
		"pending":   {},
	}

	if _, ok := StatusMap[status]; !ok {
		http.Error(w, "Невалидный статус", http.StatusBadRequest)
		return
	}

	resp, err := a.PetService.UpdateNameStatusService(petIdInt, name, status)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

// @Summary Получить всех животных по статусу
// @Description Получить всех животных по статусу
// @Tags pet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param status query string true "Статус питомца (available, sold, pending)"
// @Success 200 {array} OutputPetArray "Успешное получение животных или не ошибочное сообщение"
// @Failure 400 {string} Err "Неверной формат"
// @Failure 401 {string} Err "Аутентификация не пройдена"
// @Failure 403 {string} Err "Авторизация не пройдена"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /pet/findByStatus [get]
func (a *AnimalStore) GetPetByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	pet, err := a.PetService.GetByStatusService(status)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	OutputPetArray := []OutputPet{}
	for _, p := range pet {
		outputTagSlice := []OutputTag{}
		for _, tags := range p.Tag {
			outputTag := OutputTag{}

			outputTag.Id = tags.Id

			outputTagSlice = append(outputTagSlice, outputTag)
		}
		outputPet := OutputPet{
			Id:          p.Id,
			Name:        p.Name,
			Category_id: p.CategoryID,
			Status:      p.Status,
			OutputTag:   outputTagSlice,
		}

		OutputPetArray = append(OutputPetArray, outputPet)
	}

	err = json.NewEncoder(w).Encode(OutputPetArray)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}
}

// @Summary Получить животного по id
// @Description Получить одного животного по id
// @Tags pet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param petId path string true "id питомца"
// @Success 200 {object} OutputPet "Получить одного животного или не ошибочное сообщение"
// @Failure 400 {string} Err "Неверной формат структуры"
// @Failure 401 {string} Err "Аутентификация не пройдена"
// @Failure 403 {string} Err "Авторизация не пройдена"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /pet/{petId} [get]
func (a *AnimalStore) GetPet(w http.ResponseWriter, r *http.Request) {
	petId := chi.URLParam(r, "petId")
	petIdInt, _ := strconv.Atoi(petId)

	pet, err := a.PetService.GetService(petIdInt)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}

	OutputTagSlice := []OutputTag{}
	for _, tags := range pet.Tag {
		outputTag := OutputTag{
			Id: tags.Id,
		}

		OutputTagSlice = append(OutputTagSlice, outputTag)
	}
	outputPet := OutputPet{
		Id:          pet.Id,
		Name:        pet.Name,
		Category_id: pet.CategoryID,
		PhotoUrl:    pet.PhotoUrl,
		Status:      pet.Status,
		OutputTag:   OutputTagSlice,
	}

	err = json.NewEncoder(w).Encode(outputPet)
	if err != nil {
		http.Error(w, "Неожиданная ошибка", http.StatusInternalServerError)
		return
	}
}

// @Summary Удалить животного по id
// @Description Удалить одного животного по id
// @Tags pet
// @Security BearerAuth
// @Accept json
// @Produce plain
// @Param petId path string true "id питомца"
// @Success 200 {string} Info "Успешное удаление или не ошибочное сообщение"
// @Failure 400 {string} Err "Неверной формат структуры"
// @Failure 401 {string} Err "Аутентификация не пройдена"
// @Failure 403 {string} Err "Авторизация не пройдена"
// @Failure 500 {string} Err "Внутренняя ошибка сервера"
// @Router /pet/{petId} [delete]
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
