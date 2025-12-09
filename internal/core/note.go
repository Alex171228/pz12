package core

import "time"

// Note — доменная модель заметки.
// @Description Заметка с заголовком и содержимым
type Note struct {
	// Уникальный идентификатор заметки
	ID int64 `json:"id" example:"1"`
	// Заголовок заметки
	Title string `json:"title" example:"Моя заметка"`
	// Содержимое заметки
	Content string `json:"content" example:"Текст заметки..."`
	// Дата и время создания
	CreatedAt time.Time `json:"createdAt" example:"2024-12-08T12:00:00Z"`
	// Дата и время последнего обновления
	UpdatedAt *time.Time `json:"updatedAt,omitempty" example:"2024-12-08T13:00:00Z"`
}
