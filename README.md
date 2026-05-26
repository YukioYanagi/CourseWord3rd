# Шлюз обмена данными (Go + Python)

HTTP-шлюз с версией API **`/api/v1`**, логированием, веб-интерфейсом (две вкладки) и сервисом трансформаций **JSON / XML / SOAP** на Python.

## Быстрый старт

Подробные шаги: [`docs/ЗАПУСК_И_ПРОВЕРКА.md`](docs/ЗАПУСК_И_ПРОВЕРКА.md).

Кратко:

1. `cd python && pip install -r requirements.txt && python -m uvicorn app:app --host 127.0.0.1 --port 5000`
2. В корне: `go build -o server.exe ./cmd/server && .\server.exe`
3. Браузер: `http://127.0.0.1:8080/`

## Docker Compose

```bash
docker compose up --build
```

Шлюз: `http://localhost:8080/`, Python: порт `5000`. Данные шлюза: том `gateway-data` (`DATA_DIR=/data`).

## Документация

| Документ | Содержание |
|----------|------------|
| [`docs/API.md`](docs/API.md) | Описание HTTP API |
| [`docs/DIAGRAMS.md`](docs/DIAGRAMS.md) | Диаграммы (Mermaid) |
| [`docs/GIT_WORKFLOW.md`](docs/GIT_WORKFLOW.md) | Git Flow и семантические коммиты |
| [`docs/ЗАПУСК_И_ПРОВЕРКА.md`](docs/ЗАПУСК_И_ПРОВЕРКА.md) | Запуск и проверка вручную |

## Тесты

```bash
go test ./...
cd python && pip install -r requirements.txt -r requirements-dev.txt && pytest
cargo test --manifest-path tools/ci-skeleton/Cargo.toml
```

Скрипт проверки API: [`examples/check.ps1`](examples/check.ps1).

## CI и безопасность

- GitHub Actions: [`.github/workflows/ci.yml`](.github/workflows/ci.yml) — `go test`, `go vet`, `govulncheck`, `staticcheck`, `pytest`, `pip-audit`, `bandit`, `cargo test`, сборка Docker.
- [`.github/workflows/security.yml`](.github/workflows/security.yml) — по расписанию: `gosec`, `go mod verify`, `pip-audit`.
- Локально (Windows): `powershell -File scripts/security-check.ps1`

## Участие в разработке

См. [`CONTRIBUTING.md`](CONTRIBUTING.md) (Conventional Commits, ветвление).
