# Демонстрация работы Task Service

## Как запустить (если установлен Docker)

```bash
cd C:\scripts\test-task-for-junior-backend-developer
docker compose down -v
docker compose up --build
```

После запуска откройте в браузере:
- **Swagger UI**: http://localhost:8080/swagger/
- **API**: http://localhost:8080/api/v1/tasks

---

## Примеры API запросов

### 1. Создать шаблон ежедневной задачи

**Запрос:**
```bash
POST http://localhost:8080/api/v1/tasks
Content-Type: application/json

{
  "title": "Ежедневный обзвон пациентов",
  "description": "Обзвонить всех пациентов из списка",
  "status": "new",
  "is_template": true,
  "recurrence_type": "daily",
  "recurrence_config": {
    "interval": 1
  }
}
```

**Ответ:**
```json
{
  "id": 1,
  "title": "Ежедневный обзвон пациентов",
  "description": "Обзвонить всех пациентов из списка",
  "status": "new",
  "is_template": true,
  "parent_task_id": null,
  "recurrence_type": "daily",
  "recurrence_config": {
    "interval": 1
  },
  "next_occurrence": "2026-04-07T00:00:00Z",
  "created_at": "2026-04-06T14:10:00Z",
  "updated_at": "2026-04-06T14:10:00Z"
}
```

---

### 2. Создать шаблон ежемесячной задачи

**Запрос:**
```bash
POST http://localhost:8080/api/v1/tasks
Content-Type: application/json

{
  "title": "Формирование отчетности",
  "description": "Создать месячный отчет",
  "status": "new",
  "is_template": true,
  "recurrence_type": "monthly",
  "recurrence_config": {
    "day_of_month": 1
  }
}
```

**Ответ:**
```json
{
  "id": 2,
  "title": "Формирование отчетности",
  "description": "Создать месячный отчет",
  "status": "new",
  "is_template": true,
  "parent_task_id": null,
  "recurrence_type": "monthly",
  "recurrence_config": {
    "day_of_month": 1
  },
  "next_occurrence": "2026-05-01T00:00:00Z",
  "created_at": "2026-04-06T14:10:00Z",
  "updated_at": "2026-04-06T14:10:00Z"
}
```

---

### 3. Создать шаблон для четных дней

**Запрос:**
```bash
POST http://localhost:8080/api/v1/tasks
Content-Type: application/json

{
  "title": "Инвентаризация склада",
  "description": "Провести инвентаризацию",
  "status": "new",
  "is_template": true,
  "recurrence_type": "even_odd",
  "recurrence_config": {
    "even_odd_type": "even"
  }
}
```

**Ответ:**
```json
{
  "id": 3,
  "title": "Инвентаризация склада",
  "description": "Провести инвентаризацию",
  "status": "new",
  "is_template": true,
  "parent_task_id": null,
  "recurrence_type": "even_odd",
  "recurrence_config": {
    "even_odd_type": "even"
  },
  "next_occurrence": "2026-04-08T00:00:00Z",
  "created_at": "2026-04-06T14:10:00Z",
  "updated_at": "2026-04-06T14:10:00Z"
}
```

---

### 4. Создать шаблон на конкретные даты

**Запрос:**
```bash
POST http://localhost:8080/api/v1/tasks
Content-Type: application/json

{
  "title": "Праздничные мероприятия",
  "description": "Организовать мероприятие",
  "status": "new",
  "is_template": true,
  "recurrence_type": "specific_dates",
  "recurrence_config": {
    "dates": ["2026-05-01", "2026-06-12", "2026-12-31"]
  }
}
```

**Ответ:**
```json
{
  "id": 4,
  "title": "Праздничные мероприятия",
  "description": "Организовать мероприятие",
  "status": "new",
  "is_template": true,
  "parent_task_id": null,
  "recurrence_type": "specific_dates",
  "recurrence_config": {
    "dates": ["2026-05-01", "2026-06-12", "2026-12-31"]
  },
  "next_occurrence": "2026-05-01T00:00:00Z",
  "created_at": "2026-04-06T14:10:00Z",
  "updated_at": "2026-04-06T14:10:00Z"
}
```

---

### 5. Получить все задачи

**Запрос:**
```bash
GET http://localhost:8080/api/v1/tasks
```

**Ответ:**
```json
[
  {
    "id": 5,
    "title": "Ежедневный обзвон пациентов",
    "description": "Обзвонить всех пациентов из списка",
    "status": "new",
    "is_template": false,
    "parent_task_id": 1,
    "recurrence_type": "none",
    "recurrence_config": null,
    "next_occurrence": null,
    "created_at": "2026-04-07T00:00:00Z",
    "updated_at": "2026-04-07T00:00:00Z"
  },
  {
    "id": 1,
    "title": "Ежедневный обзвон пациентов",
    "description": "Обзвонить всех пациентов из списка",
    "status": "new",
    "is_template": true,
    "parent_task_id": null,
    "recurrence_type": "daily",
    "recurrence_config": {
      "interval": 1
    },
    "next_occurrence": "2026-04-08T00:00:00Z",
    "created_at": "2026-04-06T14:10:00Z",
    "updated_at": "2026-04-07T00:00:00Z"
  }
]
```

Обратите внимание:
- Задача с `id: 5` - это автоматически созданная задача из шаблона (id: 1)
- У неё `is_template: false` и `parent_task_id: 1`
- У шаблона обновился `next_occurrence` на следующий день

---

## Как работает Scheduler

1. **Каждый час** scheduler проверяет шаблоны
2. Находит шаблоны где `next_occurrence <= текущая_дата`
3. Создает из них обычные задачи
4. Обновляет `next_occurrence` в шаблоне

**Пример логов:**
```
2026-04-06 14:00:00 Scheduler started
2026-04-06 14:00:00 Processing 1 templates for date 2026-04-06
2026-04-06 14:00:00 Created task 5 from template 1
2026-04-06 15:00:00 Processing 0 templates for date 2026-04-06
```

---

## Структура базы данных

```sql
CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    
    -- Новые поля для периодичности
    is_template BOOLEAN NOT NULL DEFAULT FALSE,
    parent_task_id BIGINT REFERENCES tasks(id) ON DELETE SET NULL,
    recurrence_type TEXT NOT NULL DEFAULT 'none',
    recurrence_config JSONB,
    next_occurrence DATE,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Альтернатива: Установка Docker Desktop

Если хотите увидеть работающее приложение:

1. Скачайте Docker Desktop: https://www.docker.com/products/docker-desktop/
2. Установите и запустите
3. Выполните команды:
```bash
cd C:\scripts\test-task-for-junior-backend-developer
docker compose up --build
```
4. Откройте http://localhost:8080/swagger/
