# практическое задание 12
## Шишков А.Д. ЭФМО-02-22
## Тема
Подключение Swagger/OpenAPI. Автоматическая генерация документации 
## Цели
- Освоить основы спецификации OpenAPI (Swagger) для REST API.
- Подключить автогенерацию документации к проекту из ПЗ 11 (notes-api).
- Научиться публиковать интерактивную документацию (Swagger UI / ReDoc) на эндпоинте GET /docs.
- Синхронизировать код и спецификацию (комментарии-аннотации → генерация) и/или «schema-first» (генерация кода из openapi.yaml).
- одготовить процесс обновления документации (Makefile/скрипт).
### Подход: Code-First

В данном проекте используется **code-first подход** с инструментом **swag** — Swagger-документация автоматически генерируется из специальных аннотаций в Go-коде.

| Подход | Описание | Используется в проекте |
|--------|----------|------------------------|
| **Code-first** | Аннотации в коде → `swag init` → генерация `docs/` | да |
| **Schema-first** | OpenAPI YAML/JSON → `oapi-codegen` → генерация кода | нет |
| **Embedded Spec** | Ручное написание спецификации в коде | нет |

**Преимущества code-first подхода:**
- Документация всегда соответствует коду (Single Source of Truth)
- Аннотации находятся рядом с кодом — легко поддерживать
- Автоматическая генерация типов и примеров из Go-структур
- Поддержка валидации и примеров через теги

**Инструменты:**
- [swaggo/swag](https://github.com/swaggo/swag) — генератор документации
- [swaggo/http-swagger](https://github.com/swaggo/http-swagger) — middleware для Swagger UI

---
### Фрагменты кода методов API

### Общая информация об API (`cmd/api/main.go`)

```go
// @title Notes API
// @version 1.0
// @description Учебный REST API для заметок (CRUD) для практического занятия №12.
// @description Демонстрация code-first подхода с генерацией Swagger документации через swag.

// @host localhost:8080
// @BasePath /api/v1

// @schemes http
```

### Метод `ListNotes` — получение списка заметок

```go
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
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(notes)
}
```

### Метод `CreateNote` — создание заметки

```go
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
    // ... реализация
}
```

### Метод `GetNote` — получение заметки по ID

```go
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
    // ... реализация
}
```

### Модели данных с аннотациями

```go
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
```

---
### Скриншот Swagger UI

<img width="1857" height="920" alt="image" src="https://github.com/user-attachments/assets/d941710a-bbe8-489a-aa27-77806841f9c6" /> 

### Команды для запуска

### Запуск сервера

```bash
go run ./cmd/api
```

Сервер запустится на порту `8080`.
### Доступные эндпоинты

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/health` | Healthcheck |
| GET | `/docs/*` | Swagger UI |
| GET | `/swagger/doc.json` | Swagger JSON спецификация |
| GET | `/api/v1/notes` | Список заметок |
| POST | `/api/v1/notes` | Создать заметку |
| GET | `/api/v1/notes/{id}` | Получить заметку |
| PATCH | `/api/v1/notes/{id}` | Обновить заметку |
| DELETE | `/api/v1/notes/{id}` | Удалить заметку |

### Генерация документации

В данном проекте **генерация не требуется** — спецификация встроена в код.

Для справки, при использовании других подходов команды были бы следующими:

**Code-first (swag):**
```bash
# Установка
go install github.com/swaggo/swag/cmd/swag@latest

# Генерация docs/ из аннотаций
swag init -g cmd/api/main.go -o docs
```

**Schema-first (oapi-codegen):**
```bash
# Установка
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Генерация Go-кода из OpenAPI спецификации
oapi-codegen -generate types,server -package api openapi.yaml > internal/api/api.gen.go
```

---
## Запуск проекта

### Требования

| Компонент | Версия | Проверка |
|-----------|--------|----------|
| Go | 1.22+ | `go version` |
| Git | любая | `git --version` |

Клонировать репозиторий

```bash
git clone https://github.com/Alex171228/pz12
cd pz12
```

Установка зависимостей

```bash
go mod download
```

Эта команда скачает все необходимые пакеты:
- `github.com/go-chi/chi/v5` — HTTP роутер
- `github.com/swaggo/http-swagger/v2` — Swagger UI middleware

Для проверки установленных зависимостей:

```bash
go mod verify
```

Запуск сервера

```bash
go run ./cmd/api
```
### Структура проекта

<img width="303" height="515" alt="image" src="https://github.com/user-attachments/assets/630f3775-f5b7-47b8-af8e-0222d558441f" /> 

### Установка swag CLI

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Генерация документации

```bash
swag init -g cmd/api/main.go -o docs
```

**Параметры:**
- `-g cmd/api/main.go` — файл с главными аннотациями (@title, @version и т.д.)
- `-o docs` — папка для вывода сгенерированных файлов
**Результат выполнения:**

```
2025/12/08 22:46:07 Generate swagger docs....
2025/12/08 22:46:07 Generate general API Info, search dir:./
2025/12/08 22:46:08 Generating handlers.CreateNoteRequest
2025/12/08 22:46:08 Generating core.Note
2025/12/08 22:46:08 Generating handlers.ErrorResponse
2025/12/08 22:46:08 Generating handlers.UpdateNoteRequest
2025/12/08 22:46:08 create docs.go at docs/docs.go
2025/12/08 22:46:08 create swagger.json at docs/swagger.json
2025/12/08 22:46:08 create swagger.yaml at docs/swagger.yaml
```
## 6. Выводы

### Что удалось

1. **Реализован полноценный REST API** с CRUD-операциями для заметок
2. **Использован code-first подход** — документация генерируется из аннотаций в коде
3. **Автоматическая генерация** OpenAPI спецификации через `swag init`
4. **Интерактивный Swagger UI** по адресу `/docs/`
5. **Чистая архитектура** с разделением на слои (core, http, repo)

### Что автоматизировано

| Компонент | Способ автоматизации |
|-----------|----------------------|
| Swagger спецификация | `swag init` генерирует `docs/swagger.json` |
| Модели данных в docs | Автоматически из Go-структур с тегами |
| Примеры значений | Теги `example:"..."` в структурах |
| Swagger UI | Middleware `http-swagger` |

### Основные Swagger-аннотации

| Аннотация | Назначение | Пример |
|-----------|------------|--------|
| `@Summary` | Краткое описание метода | `@Summary Список заметок` |
| `@Description` | Подробное описание | `@Description Возвращает массив всех заметок` |
| `@Tags` | Группировка методов | `@Tags notes` |
| `@Accept` | Входной формат | `@Accept json` |
| `@Produce` | Выходной формат | `@Produce json` |
| `@Param` | Параметр запроса | `@Param id path int true "ID заметки"` |
| `@Success` | Успешный ответ | `@Success 200 {object} Note` |
| `@Failure` | Ответ с ошибкой | `@Failure 404 {object} ErrorResponse` |
| `@Router` | Путь и метод | `@Router /notes/{id} [get]` |
