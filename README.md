# Workout Tracker API

REST API для учёта выполненных упражнений и получения агрегированной статистики.

## Требования

- Docker + Docker Compose

## Запуск

```bash
cp .env.example .env
docker compose up -d --build
```

Миграции применяются автоматически перед стартом сервера.

## Миграции вручную

Для управления миграциями вне Docker (локальная разработка, откат):

```bash
# применить
DATABASE_URL=postgres://postgres:postgres@localhost:5432/workout_tracker?sslmode=disable go run ./cmd/migrate up

# откатить последнюю
DATABASE_URL=postgres://postgres:postgres@localhost:5432/workout_tracker?sslmode=disable go run ./cmd/migrate down

# статус
DATABASE_URL=postgres://postgres:postgres@localhost:5432/workout_tracker?sslmode=disable go run ./cmd/migrate status
```

> Требует Go 1.25+ и запущенного PostgreSQL на localhost:5432.

## API

### Упражнения

**Создать упражнение**
```
POST /api/v1/exercises
Content-Type: application/json

{"name": "Жим лёжа", "description": "Базовое упражнение на грудь"}
```

**Список упражнений**
```
GET /api/v1/exercises
```

### Пользователи

**Создать пользователя**
```
POST /api/v1/users
Content-Type: application/json

{"name": "Никита"}
```

### Логи упражнений

**Зафиксировать выполнение**
```
POST /api/v1/exercise-logs
Content-Type: application/json

{
  "user_id": "uuid",
  "exercise_id": "uuid",
  "performed_at": "2026-06-06T10:00:00Z",
  "sets": 3,
  "reps": 10,
  "weight_kg": 80.5,
  "duration_seconds": 120,
  "notes": "хорошо пошло"
}
```
Поля `performed_at`, `sets`, `reps`, `weight_kg`, `duration_seconds`, `notes` — опциональны.

**История выполнений пользователя**
```
GET /api/v1/users/{user_id}/exercise-logs
```

### Статистика

**Статистика пользователя**
```
GET /api/v1/users/{user_id}/stats
```

Ответ:
```json
{
  "total": 42,
  "today": 3,
  "last_7_days": [
    {"date": "2026-05-31", "count": 5},
    {"date": "2026-06-01", "count": 3}
  ]
}
```

## Структура проекта

```
cmd/
  server/       — точка входа HTTP-сервера
  migrate/      — CLI для управления миграциями
internal/
  config/       — загрузка конфигурации из env
  httputil/     — общие хелперы JSON-ответов и валидации
  router/       — сборка HTTP-роутера и middleware
  exercise/     — упражнения: handler, service, repository
  user/         — пользователи: handler, service, repository
  exerciselog/  — логи выполнений и статистика: handler, service, repository
migrations/     — SQL-миграции (goose)
```

## Технические решения

**go-chi/chi** — минималистичный роутер, стандартные middleware (логирование, recover). Выбран за идиоматичность и отсутствие лишних абстракций.

**jackc/pgx/v5** — современный PostgreSQL-драйвер с нативной поддержкой типов и пулом соединений. Производительнее `database/sql` + `lib/pq`.

**pressly/goose** — миграции в одном SQL-файле с аннотациями `-- +goose Up / Down`. Удобнее `golang-migrate` (два файла на миграцию), поддерживает Go-миграции при необходимости.

**Группировка по сущности** — каждый пакет в `internal/` содержит handler, service и repository одной сущности. Проще навигация, нет лишних слоёв абстракции для проекта такого размера.

**Агрегации в PostgreSQL** — статистика считается одним CTE-запросом с `json_agg`. Никакой агрегации на стороне приложения.

**Статистика — часть `exerciselog`, а не отдельный пакет** — она лишь агрегированное представление тех же `exercise_logs`, без своей предметной логики. Отдельный пакет добавил бы лишний слой косвенности без пользы.

**Миграции отдельным бинарником** — `cmd/migrate` запускается независимо от сервера, что соответствует принципу разделения ответственности и позволяет применять миграции до деплоя.
