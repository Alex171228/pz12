package httpx

const swaggerJSON = `{
  "swagger": "2.0",
  "info": {
    "title": "Notes API",
    "version": "1.0",
    "description": "Учебный REST API для заметок (CRUD) для практического занятия №12."
  },
  "basePath": "/api/v1",
  "schemes": [
    "http"
  ],
  "paths": {
    "/notes": {
      "get": {
        "summary": "Список заметок",
        "tags": [
          "notes"
        ],
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/Note"
              }
            }
          },
          "500": {
            "description": "Internal error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      },
      "post": {
        "summary": "Создать заметку",
        "tags": [
          "notes"
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "parameters": [
          {
            "in": "body",
            "name": "input",
            "required": true,
            "schema": {
              "$ref": "#/definitions/NoteCreate"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Created",
            "schema": {
              "$ref": "#/definitions/Note"
            }
          },
          "400": {
            "description": "Validation error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "500": {
            "description": "Internal error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/notes/{id}": {
      "parameters": [
        {
          "name": "id",
          "in": "path",
          "required": true,
          "type": "integer",
          "format": "int64"
        }
      ],
      "get": {
        "summary": "Получить заметку по ID",
        "tags": [
          "notes"
        ],
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/Note"
            }
          },
          "400": {
            "description": "Invalid id",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "404": {
            "description": "Not found",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "500": {
            "description": "Internal error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      },
      "patch": {
        "summary": "Частично обновить заметку",
        "tags": [
          "notes"
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "parameters": [
          {
            "in": "body",
            "name": "input",
            "required": true,
            "schema": {
              "$ref": "#/definitions/NoteUpdate"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Updated",
            "schema": {
              "$ref": "#/definitions/Note"
            }
          },
          "400": {
            "description": "Invalid data",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "404": {
            "description": "Not found",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "500": {
            "description": "Internal error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      },
      "delete": {
        "summary": "Удалить заметку",
        "tags": [
          "notes"
        ],
        "produces": [
          "application/json"
        ],
        "responses": {
          "204": {
            "description": "No Content"
          },
          "400": {
            "description": "Invalid id",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "404": {
            "description": "Not found",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "500": {
            "description": "Internal error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Note": {
      "type": "object",
      "properties": {
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "title": {
          "type": "string"
        },
        "content": {
          "type": "string"
        },
        "createdAt": {
          "type": "string",
          "format": "date-time"
        },
        "updatedAt": {
          "type": "string",
          "format": "date-time"
        }
      }
    },
    "NoteCreate": {
      "type": "object",
      "required": [
        "title"
      ],
      "properties": {
        "title": {
          "type": "string",
          "example": "Новая заметка"
        },
        "content": {
          "type": "string",
          "example": "Текст заметки"
        }
      }
    },
    "NoteUpdate": {
      "type": "object",
      "properties": {
        "title": {
          "type": "string",
          "example": "Обновлённый заголовок"
        },
        "content": {
          "type": "string",
          "example": "Обновлённый текст"
        }
      }
    },
    "ErrorResponse": {
      "type": "object",
      "properties": {
        "error": {
          "type": "string",
          "example": "something went wrong"
        }
      }
    }
  }
}`
