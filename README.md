# Планировщик задач

REST API для управления пользователями, задачами и тегами.

Деплой - TODO

## Возможности

- регистрация и логин пользователя
- получение текущего пользователя
- создание, получение, обновление и удаление задач
- создание, получение списка и удаление тегов
- привязка тегов к задачам

## Запуск

Перед запуском нужно задать переменные окружения:

```bash
export DB_URL='postgres://postgres:postgres@localhost:5432/taskscheduler?sslmode=disable'
export JWT_SECRET='super-secret'
```

После этого приложение можно запустить так:

```bash
docker-composr up --build
```
Порт HTTP-сервера берется из [internal/config/config.yaml](/home/smag/Study/Pet/TaskScheduler/internal/config/config.yaml).

## API

### Auth

- `POST /auth/register` - регистрация пользователя
- `POST /auth/login` - вход пользователя
- `GET /auth/me` - получить текущего пользователя

### Tasks

- `POST /tasks` - создать задачу
- `GET /tasks/{id}` - получить задачу пользователя по `id`
- `PUT /tasks/{id}` - обновить задачу пользователя по `id`
- `DELETE /tasks/{id}` - удалить задачу пользователя по `id`

### Tags

- `POST /tags` - создать тег
- `GET /tags` - получить список тегов пользователя
- `DELETE /tags/{id}` - удалить тег пользователя по `id`



