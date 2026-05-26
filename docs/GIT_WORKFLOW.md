# Git, семантические коммиты и Git Flow

## Ветки (Git Flow)

| Ветка | Назначение |
|--------|------------|
| **`main`** | стабильные релизы, только слияния из `release/*` или hotfix |
| **`develop`** | интеграция фич, «почти готово к релизу» |
| **`feature/<кратко>`** | одна задача / фича, ответвление от `develop` |
| **`release/<версия>`** | заморозка перед релизом, только правки версии и багфиксы |
| **`hotfix/<кратко>`** | срочный фикс в проде, от `main`, затем merge в `main` и `develop` |

Типичный цикл фичи:

1. `git checkout develop && git pull`
2. `git checkout -b feature/add-docker`
3. коммиты в ветке фичи
4. Pull Request в `develop`, ревью, merge
5. перед релизом: из `develop` — ветка `release/1.1.0` → правки → PR в `main` и обратно в `develop`

## Семантические коммиты (Conventional Commits)

Формат заголовка:

```text
<тип>(<область>): <краткое описание>
```

**Типы:** `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `ci`, `perf`, `build`.

**Примеры:**

- `feat(gateway): add docker compose for gateway and transform`
- `fix(storage): handle missing index file on first run`
- `docs: update API description for send endpoint`
- `ci: add gosec and pip-audit workflows`
- `test(python): add pytest for transform endpoint`

Тело коммита (по желанию): что сделано, зачем, breaking changes (`BREAKING CHANGE:` в подвале).

## Шаблон сообщения Git

В корне репозитория: [`.gitmessage`](../.gitmessage).

Подключение один раз:

```bash
git config commit.template .gitmessage
```

## Теги релизов

Семвер: `v1.0.0`, `v1.1.0`. Тег ставится на коммит в **`main`** после merge релиза:

```bash
git tag -a v1.0.0 -m "Release 1.0.0"
git push origin v1.0.0
```
