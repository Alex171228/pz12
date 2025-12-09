package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"example.com/notes-api/internal/core"
	"example.com/notes-api/internal/core/service"
	"example.com/notes-api/internal/repo"
)

// Handler содержит зависимости для HTTP-обработчиков.
type Handler struct {
	Service *service.NoteService
}

// NewHandler создаёт новый Handler.
func NewHandler(s *service.NoteService) *Handler {
	return &Handler{Service: s}
}

// CreateNoteRequest модель запроса на создание заметки.
// @Description Данные для создания новой заметки
type CreateNoteRequest struct {
	// Заголовок заметки (обязательное поле)
	Title string `json:"title" example:"Моя первая заметка"`
	// Содержимое заметки
	Content string `json:"content" example:"Текст заметки..."`
}

// UpdateNoteRequest модель запроса на обновление заметки.
// @Description Данные для частичного обновления заметки
type UpdateNoteRequest struct {
	// Новый заголовок (опционально)
	Title *string `json:"title,omitempty" example:"Обновлённый заголовок"`
	// Новое содержимое (опционально)
	Content *string `json:"content,omitempty" example:"Обновлённый текст"`
}

// ErrorResponse модель ответа с ошибкой.
// @Description Ответ сервера при возникновении ошибки
type ErrorResponse struct {
	// Описание ошибки
	Error string `json:"error" example:"something went wrong"`
}

// вспомогательная функция для ошибок.
func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}

// CreateNote создаёт новую заметку.
// @Summary Создать заметку
// @Description Создаёт новую заметку с указанным заголовком и содержимым
// @Tags notes
// @Accept json
// @Produce json
// @Param input body CreateNoteRequest true "Данные заметки"
// @Success 201 {object} core.Note "Созданная заметка"
// @Failure 400 {object} ErrorResponse "Ошибка валидации"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /notes [post]
func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var input CreateNoteRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	note, err := h.Service.CreateNote(input.Title, input.Content)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			writeError(w, http.StatusBadRequest, "title is required")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	_ = json.NewEncoder(w).Encode(note)
}

// ListNotes возвращает список всех заметок.
// @Summary Список заметок
// @Description Возвращает массив всех заметок
// @Tags notes
// @Produce json
// @Success 200 {array} core.Note "Список заметок"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /notes [get]
func (h *Handler) ListNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.Service.ListNotes()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	
	// Возвращаем пустой массив вместо null
	if notes == nil {
		notes = []core.Note{}
	}
	
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(notes)
}

// GetNote возвращает заметку по ID.
// @Summary Получить заметку
// @Description Возвращает заметку по её идентификатору
// @Tags notes
// @Produce json
// @Param id path int true "ID заметки"
// @Success 200 {object} core.Note "Найденная заметка"
// @Failure 400 {object} ErrorResponse "Некорректный ID"
// @Failure 404 {object} ErrorResponse "Заметка не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /notes/{id} [get]
func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	note, err := h.Service.GetNote(id)
	if err != nil {
		if errors.Is(err, repo.ErrNoteNotFound) {
			writeError(w, http.StatusNotFound, "note not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(note)
}

// UpdateNote частично обновляет заметку.
// @Summary Обновить заметку
// @Description Частично обновляет заметку (PATCH). Можно обновить только title, только content или оба поля.
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "ID заметки"
// @Param input body UpdateNoteRequest true "Данные для обновления"
// @Success 200 {object} core.Note "Обновлённая заметка"
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 404 {object} ErrorResponse "Заметка не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /notes/{id} [patch]
func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// Конвертируем в service.NoteUpdateInput
	updateInput := service.NoteUpdateInput{
		Title:   input.Title,
		Content: input.Content,
	}

	note, err := h.Service.UpdateNote(id, updateInput)
	if err != nil {
		if errors.Is(err, repo.ErrNoteNotFound) {
			writeError(w, http.StatusNotFound, "note not found")
			return
		}
		if errors.Is(err, service.ErrValidation) {
			writeError(w, http.StatusBadRequest, "invalid data")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(note)
}

// DeleteNote удаляет заметку.
// @Summary Удалить заметку
// @Description Удаляет заметку по ID. При успехе возвращает 204 No Content.
// @Tags notes
// @Param id path int true "ID заметки"
// @Success 204 "Заметка удалена"
// @Failure 400 {object} ErrorResponse "Некорректный ID"
// @Failure 404 {object} ErrorResponse "Заметка не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /notes/{id} [delete]
func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.Service.DeleteNote(id); err != nil {
		if errors.Is(err, repo.ErrNoteNotFound) {
			writeError(w, http.StatusNotFound, "note not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204, без тела
}
