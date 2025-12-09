// Notes API server for практическое занятие №12.
package main

import (
	"log"
	"net/http"

	_ "example.com/notes-api/docs" // swagger docs

	"example.com/notes-api/internal/core/service"
	httpx "example.com/notes-api/internal/http"
	"example.com/notes-api/internal/http/handlers"
	"example.com/notes-api/internal/repo"
)

// @title Notes API
// @version 1.0
// @description Учебный REST API для заметок (CRUD) для практического занятия №12.
// @description Демонстрация code-first подхода с генерацией Swagger документации через swag.

// @host localhost:8080
// @BasePath /api/v1

// @schemes http

func main() {
	// Инициализация репозитория и сервиса
	rp := repo.NewNoteRepoMem()
	svc := service.NewNoteService(rp)
	h := handlers.NewHandler(svc)

	router := httpx.NewRouter(h)

	addr := ":8080" // слушаем на всех интерфейсах
	log.Println("Server started at", addr)
	log.Println("Swagger UI: http://localhost:8080/docs/")
	log.Fatal(http.ListenAndServe(addr, router))
}
