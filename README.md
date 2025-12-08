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
