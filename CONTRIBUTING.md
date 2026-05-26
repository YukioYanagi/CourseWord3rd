# Участие в проекте

## Коммиты

Используйте [Conventional Commits](https://www.conventionalcommits.org/) (см. [`docs/GIT_WORKFLOW.md`](docs/GIT_WORKFLOW.md)).

Рекомендуется шаблон: `git config commit.template .gitmessage`.

## Ветки

Git Flow: фичи в `feature/*` от `develop`, стабильность в `main`. Подробности — в [`docs/GIT_WORKFLOW.md`](docs/GIT_WORKFLOW.md).

## Перед отправкой PR

1. `go fmt ./...` и `go test ./...`
2. `cd python && pytest`
3. При наличии Rust toolchain: `cargo test --manifest-path tools/ci-skeleton/Cargo.toml`
4. Локально при необходимости: `scripts/security-check.ps1`

## CI

Проверки в [`.github/workflows/ci.yml`](.github/workflows/ci.yml) должны быть зелёными.
