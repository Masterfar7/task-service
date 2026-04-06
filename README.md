# Task Service - Трекер задач с поддержкой периодичности

Сервис для управления задачами с HTTP API на Go. Реализована поддержка периодических задач (ежедневные, ежемесячные, на конкретные даты, четные/нечетные дни).

## Описание решения

Это решение тестового задания для позиции junior backend разработчика. Подробное описание реализации, архитектурных решений и принятых предположений находится в файле [SOLUTION.md](SOLUTION.md).


## Быстрый запуск

### Через Docker Compose (рекомендуется)

```bash
# Клонируйте репозиторий
git clone <your-repo-url>
cd test-task-for-junior-backend-developer

# Запустите сервис
docker compose up --build
```

Если база данных уже запускалась ранее, пересоздайте volumes:

```bash
docker compose down -v
docker compose up --build
```

### Локальный запуск (без Docker)

```bash
# Установите зависимости
go mod download

# Запустите PostgreSQL (например, через Docker)
docker run -d \
  --name taskservice-postgres \
  -e POSTGRES_DB=taskservice \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine

# Примените миграции
psql -h localhost -U postgres -d taskservice -f migrations/0001_create_tasks.up.sql
psql -h localhost -U postgres -d taskservice -f migrations/0002_add_recurrence.up.sql

# Запустите приложение
export DATABASE_DSN="postgres://postgres:postgres@localhost:5432/taskservice?sslmode=disable"
export HTTP_ADDR=":8080"
go run cmd/api/main.go
```

## Доступ к сервису

После запуска сервис будет доступен:

- **Главная страница**: http://localhost:8080/
- **Swagger UI**: http://localhost:8080/swagger/
- **OpenAPI JSON**: http://localhost:8080/swagger/openapi.json
- **API**: http://localhost:8080/api/v1/tasks

## API Endpoints

Базовый префикс: `/api/v1`

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/tasks` | Создать задачу или шаблон |
| GET | `/tasks` | Получить список всех задач |
| GET | `/tasks/{id}` | Получить задачу по ID |
| PUT | `/tasks/{id}` | Обновить задачу |
| DELETE | `/tasks/{id}` | Удалить задачу |

## Примеры использования

### Создание обычной задачи

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Провести операцию",
    "description": "Операция пациента Иванова",
    "status": "new"
  }'
```

### Создание шаблона ежедневной задачи

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Ежедневный обзвон пациентов",
    "description": "Обзвонить всех пациентов из списка",
    "status": "new",
    "is_template": true,
    "recurrence_type": "daily",
    "recurrence_config": {
      "interval": 1
    }
  }'
```

### Создание шаблона для четных дней месяца

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Формирование отчетности",
    "description": "Создать отчет за четные дни",
    "status": "new",
    "is_template": true,
    "recurrence_type": "even_odd",
    "recurrence_config": {
      "even_odd_type": "even"
    }
  }'
```

### Получение всех задач

```bash
curl http://localhost:8080/api/v1/tasks
```

Больше примеров в [SOLUTION.md](SOLUTION.md) и в Swagger UI.

## Типы периодичности

1. **daily** - ежедневные задачи (каждые N дней)
2. **monthly** - ежемесячные задачи (на определенное число месяца)
3. **specific_dates** - задачи на конкретные даты
4. **even_odd** - задачи на четные или нечетные дни месяца

Подробное описание форматов конфигурации в [SOLUTION.md](SOLUTION.md).

## Тестирование

Используйте Swagger UI для интерактивного тестирования API: http://localhost:8080/swagger/

