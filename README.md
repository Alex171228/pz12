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

## 2. Фрагменты аннотаций над методами

### Главные аннотации API

> **Файл:** `cmd/api/main.go`

```go
// @title Notes API
// @version 1.0
// @description Учебный REST API для заметок (CRUD) для практического занятия №12.
// @description Демонстрация code-first подхода с генерацией Swagger документации через swag.

// @host localhost:8080
// @BasePath /api/v1

// @schemes http
```

**Разбор аннотаций:**
| Аннотация | Назначение |
|-----------|------------|
| `@title` | Название API (отображается в шапке Swagger UI) |
| `@version` | Версия API |
| `@description` | Описание API (можно несколько строк) |
| `@host` | Хост для запросов |
| `@BasePath` | Базовый путь API |
| `@schemes` | Поддерживаемые протоколы (http, https) |

---

### Аннотации над методами

> **Файл:** `internal/http/handlers/notes.go`

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
	
	// Возвращаем пустой массив вместо null
	if notes == nil {
		notes = []core.Note{}
	}
	
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(notes)
}
```

**Разбор аннотаций:**
| Аннотация | Значение |
|-----------|----------|
| `@Summary` | Краткое название в Swagger UI |
| `@Description` | Подробное описание |
| `@Tags notes` | Группировка в секцию "notes" |
| `@Produce json` | Возвращает JSON |
| `@Success 200` | HTTP 200, массив объектов Note |
| `@Failure 500` | HTTP 500, объект ErrorResponse |
| `@Router /notes [get]` | GET /notes |

---

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
```

**Разбор аннотаций:**
| Аннотация | Значение |
|-----------|----------|
| `@Accept json` | Принимает JSON в теле запроса |
| `@Param input body CreateNoteRequest true` | Обязательный параметр в body |
| `@Success 201` | HTTP 201 Created |
| `@Router /notes [post]` | POST /notes |

---

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
```

**Разбор аннотаций:**
| Аннотация | Значение |
|-----------|----------|
| `@Param id path int true` | Обязательный параметр `id` в URL path, тип int |
| `@Failure 404` | HTTP 404 если заметка не найдена |
| `@Router /notes/{id} [get]` | GET /notes/{id} |

---

### Модели данных с аннотациями

> **Файл:** `internal/http/handlers/notes.go`

```go
// CreateNoteRequest модель запроса на создание заметки.
// @Description Данные для создания новой заметки
type CreateNoteRequest struct {
	// Заголовок заметки (обязательное поле)
	Title string `json:"title" example:"Моя первая заметка"`
	// Содержимое заметки
	Content string `json:"content" example:"Текст заметки..."`
}

// ErrorResponse модель ответа с ошибкой.
// @Description Ответ сервера при возникновении ошибки
type ErrorResponse struct {
	// Описание ошибки
	Error string `json:"error" example:"something went wrong"`
}
```

> **Файл:** `internal/core/note.go`

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

**Теги в структурах:**
| Тег | Назначение |
|-----|------------|
| `json:"id"` | Имя поля в JSON |
| `example:"1"` | Пример значения для Swagger UI |
| `// комментарий` | Описание поля в документации |

---

## 3. Скриншот Swagger UI

После запуска сервера Swagger UI доступен по адресу:

```
http://109.237.98.39:8080/docs/
```

<img width="1816" height="1284" alt="image" src="https://github.com/user-attachments/assets/007cada6-5274-4623-b6e7-926d6629b7db" /> 
 
*Что видно в интерфейсе:*
- Название и описание API
- Список всех эндпоинтов в группе `notes`
- Кнопка "Try it out" для интерактивного тестирования
- Схемы моделей данных (Note, CreateNoteRequest, UpdateNoteRequest, ErrorResponse)

---

## 4. Команда генерации документации

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

### Альтернатива (без установки в PATH)

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -o docs
```

---

## 5. Структура проекта

<img width="311" height="595" alt="image" src="https://github.com/user-attachments/assets/2c85bfdd-a52b-4917-ae1b-167cd88227b8" /> 

### Ключевые файлы

| Файл | Назначение |
|------|------------|
| `cmd/api/main.go` | Главные аннотации API (@title, @version, @host, @BasePath) |
| `internal/http/handlers/notes.go` | Аннотации для каждого эндпоинта (@Summary, @Router и т.д.) |
| `internal/core/note.go` | Модель данных с аннотациями (@Description, example) |
| `docs/docs.go` | Сгенерированный Go-пакет (импортируется в main.go) |
| `docs/swagger.json` | OpenAPI спецификация для Swagger UI |

---

## Запуск проекта

### Требования

| Компонент | Версия | Проверка |
|-----------|--------|----------|
| Go | 1.22+ | `go version` |
| Git | любая | `git --version` |

### Шаг 1. Клонирование репозитория

```bash
git clone <repo-url>
cd notes-api-pz12
```

### Шаг 2. Установка зависимостей

```bash
go mod download
```

### Шаг 3. Генерация документации (если docs/ отсутствует)

```bash
# Установить swag (один раз)
go install github.com/swaggo/swag/cmd/swag@latest

# Сгенерировать документацию
swag init -g cmd/api/main.go -o docs
```

### Шаг 4. Запуск сервера

```bash
go run ./cmd/api
```

После запуска в консоли появится:

```
2024/12/08 12:00:00 Server started at :8080
2024/12/08 12:00:00 Swagger UI: http://localhost:8080/docs/
```

### Шаг 5. Открытие Swagger UI

Откройте в браузере:

```
http://109.237.98.39:8080/docs/
```

---

## Доступные URL-адреса

| URL | Описание |
|-----|----------|
| http://109.237.98.39:8080/health | Healthcheck (возвращает `OK`) |
| http://109.237.98.39:8080/docs/ | Swagger UI — интерактивная документация |
| http://109.237.98.39:8080/docs/doc.json | OpenAPI спецификация в формате JSON |
| http://109.237.98.39:8080/api/v1/notes | API заметок |

---

## Примеры запросов

### Использование curl

```bash
# Создать заметку
curl -X POST http://109.237.98.39:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -d '{"title": "Первая заметка", "content": "Текст заметки"}'

# Получить все заметки
curl http://109.237.98.39:8080/api/v1/notes

# Получить заметку по ID
curl http://109.237.98.39:8080/api/v1/notes/1

# Обновить заметку
curl -X PATCH http://109.237.98.39:8080/api/v1/notes/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Обновлённый заголовок"}'

# Удалить заметку
curl -X DELETE http://109.237.98.39:8080/api/v1/notes/1
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


## Зависимости

```
github.com/go-chi/chi/v5 v5.0.12           # HTTP роутер
github.com/swaggo/http-swagger/v2 v2.0.2   # Swagger UI middleware
github.com/swaggo/swag v1.16.4             # Генератор документации (dev)
```

---

