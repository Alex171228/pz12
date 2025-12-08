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
### Подход: Ручная спецификация (Embedded Spec)

В данном проекте использован **гибридный подход** — OpenAPI 2.0 (Swagger) спецификация написана **вручную** и встроена непосредственно в исходный код как Go-константа.

| Подход | Описание | Используется в проекте |
|--------|----------|------------------------|
| **Code-first** | Аннотации в коде → `swag init` → генерация `docs/` | нет |
| **Schema-first** | OpenAPI YAML/JSON → `oapi-codegen` → генерация кода | нет |
| **Embedded Spec** | Ручное написание спецификации в коде | да |

**Преимущества выбранного подхода:**
- Полный контроль над спецификацией
- Нет зависимости от генераторов
- Спецификация всегда синхронизирована с приложением

**Недостатки:**
- Требуется ручное обновление при изменении API
- Нет автоматической валидации соответствия кода и спецификации

---
### Фрагменты кода методов API

### Метод `ListNotes` — получение списка заметок

```go
// GET /api/v1/notes
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

**Соответствующий фрагмент Swagger-спецификации:**

```json
"/notes": {
  "get": {
    "summary": "Список заметок",
    "tags": ["notes"],
    "produces": ["application/json"],
    "responses": {
      "200": {
        "description": "OK",
        "schema": {
          "type": "array",
          "items": { "$ref": "#/definitions/Note" }
        }
      },
      "500": {
        "description": "Internal error",
        "schema": { "$ref": "#/definitions/ErrorResponse" }
      }
    }
  }
}
```

---

### Метод `CreateNote` — создание заметки

```go
// POST /api/v1/notes
func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Title   string `json:"title"`
        Content string `json:"content"`
    }

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

**Соответствующий фрагмент Swagger-спецификации:**

```json
"post": {
  "summary": "Создать заметку",
  "tags": ["notes"],
  "consumes": ["application/json"],
  "produces": ["application/json"],
  "parameters": [{
    "in": "body",
    "name": "input",
    "required": true,
    "schema": { "$ref": "#/definitions/NoteCreate" }
  }],
  "responses": {
    "201": {
      "description": "Created",
      "schema": { "$ref": "#/definitions/Note" }
    },
    "400": {
      "description": "Validation error",
      "schema": { "$ref": "#/definitions/ErrorResponse" }
    }
  }
}
```

---

### Метод `GetNote` — получение заметки по ID

```go
// GET /api/v1/notes/{id}
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

### Генерация документации

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
