# HTTP API (версия v1)

Базовый префикс: **`/api/v1`**. Для маршрутов API в ответ добавляется заголовок **`X-API-Version: v1`**.

## Общие сведения

| Метод | Путь | Описание |
|--------|------|-----------|
| GET | `/api/v1/health` | Проверка доступности шлюза |
| POST | `/api/v1/send` | Приём файла (`multipart/form-data`) |
| GET | `/api/v1/received` | Список сохранённых записей (JSON) |
| GET | `/api/v1/received/{id}/download` | Скачивание файла по `id` |

Статическая раздача веб-интерфейса: **`GET /`** (каталог `web/`).

---

## GET /api/v1/health

**Ответ 200** (`application/json`):

```json
{"status":"ok","version":"v1"}
```

---

## POST /api/v1/send

**Content-Type:** `multipart/form-data`

| Поле | Обязательно | Значение |
|------|-------------|----------|
| `file` | да | тело файла |
| `format` | нет, по умолчанию `json` | `json` \| `xml` \| `soap` |
| `transform_to` | нет | целевой формат для сервиса Python: `json` \| `xml` \| `soap` |

**Успех:** `201 Created`, тело — объект записи: `id`, `format`, `filename`, `size_bytes`, `created_at`.

**Ошибки:** `400` (нет файла / неверный формат), `422` (валидация содержимого), `502` (ошибка Python при преобразовании).

---

## GET /api/v1/received

**Ответ 200:** JSON-массив записей (новые первыми).

---

## GET /api/v1/received/{id}/download

**Ответ 200:** бинарное тело файла, `Content-Disposition: attachment`.

**404:** запись не найдена или файл отсутствует на диске.

---

## Сервис Python (трансформация)

По умолчанию: `http://127.0.0.1:5000`. Задаётся переменной **`PYTHON_TRANSFORM_URL`** для шлюза Go.

### POST /transform

**Content-Type:** `application/json`

```json
{
  "source_format": "json",
  "target_format": "xml",
  "payload": "{\"a\": 1}"
}
```

**Ответ 200:**

```json
{"result": "<...>", "error": null}
```

**Ошибка валидации / преобразования:** `400`, тело FastAPI `{"detail": "..."}`.

### GET /health

Проверка сервиса Python.
