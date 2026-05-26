# Диаграммы

Ниже используется [Mermaid](https://mermaid.js.org/) (рендер в GitHub, VS Code с расширением Mermaid и т.п.).

## Компоненты и поток запроса

```mermaid
flowchart LR
  subgraph Client
    B[Браузер / 1С / curl]
  end
  subgraph Docker["Docker Compose"]
    G[Gateway Go :8080]
    P[Transform Python :5000]
    V[(Том gateway-data)]
  end
  B -->|HTTP /api/v1| G
  G -->|POST /transform при transform_to| P
  G --> V
```

## Последовательность: отправка с преобразованием

```mermaid
sequenceDiagram
  participant C as Клиент
  participant GW as Go шлюз
  participant PY as Python
  participant FS as Файловое хранилище
  C->>GW: POST /api/v1/send multipart
  GW->>GW: Валидация формата
  alt transform_to задан
    GW->>PY: POST /transform
    PY-->>GW: result
  end
  GW->>FS: Сохранить файл + index.json
  GW-->>C: 201 + метаданные записи
```

## Git Flow (ветвление)

```mermaid
gitGraph
  commit id: "init"
  branch develop
  checkout develop
  commit id: "feat-1"
  branch feature/x
  checkout feature/x
  commit id: "feat"
  checkout develop
  merge feature/x
  checkout main
  merge develop tag: "v1.0.0"
```
