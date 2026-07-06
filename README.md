# AnalogYouTube

AnalogYouTube - это backend-проект на Go для сервиса, похожего на YouTube. В проекте используется разделение по слоям: HTTP delivery, usecase, repository, domain models, migrations и infrastructure.

## Стек

- Go
- Gorilla Mux
- PostgreSQL
- pgx
- zerolog
- JWT access и refresh tokens
- FFmpeg
- WebSocket
- SQL migrations

## Основной функционал

- Регистрация, авторизация и обновление токенов
- Профиль пользователя и загрузка аватарки
- Категории
- Реальная загрузка видео как файла
- Загрузка обложки видео
- Обработка видео через FFmpeg в качествах: 1080p, 720p, 480p, 360p
- Скорости воспроизведения: 0.25, 1.0, 1.25, 1.5, 2.0
- Лайки, комментарии, просмотры и рекомендации видео
- Подписки
- Донаты
- Публичные плейлисты
- Приватные чаты через WebSocket

## Переменные окружения

Создайте файл `.env` на основе `.env.example`:

```env
SERVER_URL=localhost
SERVER_NAME=GlobalServer
PORT=6666
MODE=debug

JWT_SECRET=Scibidi_toilet
ACCESS_TOKEN_TTL_MINUTES=15
REFRESH_TOKEN_TTL_DAYS=30

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=change_me
POSTGRES_DATABASE=analogyoutube

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

Важно: название базы данных должно оставаться `analogyoutube`.
Файл `configs.json` больше не используется: реальные настройки берутся из `.env` или системных переменных окружения.

## База данных

Создайте базу данных:

```sql
CREATE DATABASE analogyoutube;
```

После этого примените миграции из папки `migrations` по порядку. Если используется `golang-migrate`, команда может выглядеть так:

```bash
migrate -path migrations -database "postgres://postgres:your_password@localhost:5432/analogyoutube?sslmode=disable" up
```

## FFmpeg

Установите FFmpeg и проверьте, что он доступен из терминала:

```bash
ffmpeg -version
```

Backend вызывает бинарник как `ffmpeg`, поэтому он должен быть доступен в `PATH`.

## Запуск

```bash
go run cmd/main.go
```

По умолчанию сервер запускается на порту `6666`.

Проверка, что сервер работает:

```http
GET http://localhost:6666/ping
```

## Основные endpoints

Auth:

- `POST /auth/sign-up`
- `POST /auth/sign-in`
- `GET /auth/refresh`
- `POST /api/register`
- `POST /api/login`

Профиль:

- `GET /api/me`
- `PUT /api/me`
- `GET /api/users/{id}`
- `GET /api/users/{id}/videos`
- `GET /api/users/{id}/subscribers/count`
- `GET /api/users/{id}/subscriptions/count`

Видео:

- `GET /api/videos`
- `GET /api/videos/search?title=golang`
- `GET /api/videos/playback-speeds`
- `GET /api/videos/{id}`
- `POST /api/videos`
- `PUT /api/videos/{id}`
- `DELETE /api/videos/{id}`
- `DELETE /api/admin/videos/{id}`

Категории:

- `GET /api/categories`
- `GET /api/categories/{name}`
- `POST /api/categories`
- `PUT /api/categories/{id}`
- `DELETE /api/categories/{id}`

Лайки и комментарии:

- `POST /api/videos/{id}/like`
- `DELETE /api/videos/{id}/like`
- `GET /api/videos/{id}/liked`
- `GET /api/videos/{id}/likes/count`
- `GET /api/videos/{id}/comments`
- `POST /api/videos/{id}/comments`
- `PUT /api/comments/{id}`
- `DELETE /api/comments/{id}`

Для создания и обновления комментария используется поле `text`:

```json
{
  "text": "Great video!"
}
```

Подписки:

- `POST /api/users/{id}/subscribe`
- `DELETE /api/users/{id}/subscribe`
- `GET /api/users/{id}/subscribed`

Донаты:

- `POST /api/donations`
- `GET /api/donations/sent`
- `GET /api/donations/received`
- `GET /api/users/{id}/donations`

Для доната пользователю без привязки к видео нужен `receiver_id`:

```json
{
  "receiver_id": 2,
  "amount": 10.5,
  "message": "Great content!"
}
```

Для доната под конкретным видео можно передать только `video_id`; получателем станет автор видео:

```json
{
  "video_id": 1,
  "amount": 10.5,
  "message": "Great video!"
}
```

Плейлисты:

- `GET /api/users/{id}/playlists`
- `GET /api/playlists/{id}`
- `POST /api/playlists`
- `PUT /api/playlists/{id}`
- `DELETE /api/playlists/{id}`
- `POST /api/playlists/{id}/videos`
- `DELETE /api/playlists/{id}/videos/{video_id}`

Чаты:

- `POST /api/chats` - отправить заявку на личный чат пользователю.
  ```json
  {
    "user_id": 2
  }
  ```
- `GET /api/chats/requests/incoming` - получить входящие pending-заявки.
- `GET /api/chats/requests/outgoing` - получить исходящие pending-заявки.
- `POST /api/chats/requests/{id}/accept` - принять заявку. После этого создается чат.
- `POST /api/chats/requests/{id}/reject` - отклонить заявку. Чат не создается.
- `GET /api/chats` - получить мои уже созданные чаты.
- `GET /api/chats/{id}/messages` - получить историю сообщений чата.
- `POST /api/chats/{id}/messages` - отправить сообщение через обычный HTTP-запрос.
  ```json
  {
    "text": "Hello! How are you?"
  }
  ```
- `GET /ws/chats/{id}?access_token=YOUR_ACCESS_TOKEN` - подключиться к чату через WebSocket.

Сообщение в WebSocket отправляется так:

```json
{
  "text": "Hello from WebSocket!"
}
```

Читать сообщения, отправлять сообщения и подключаться по WebSocket могут только участники конкретного чата.
Пока заявка в статусе `pending`, второй такой pending между этими же пользователями не создается.

## Проверка загрузки видео

Рекомендуемый сценарий перед демонстрацией:

1. Зарегистрировать пользователя.
2. Войти в аккаунт и скопировать `access_token`.
3. Создать категорию.
4. Загрузить настоящее видео через `POST /api/videos` с `form-data`.
5. Проверить, что оригинальный файл сохранился в папке `uploads`.
6. Проверить, что FFmpeg создал файлы 1080p, 720p, 480p и 360p.
7. Открыть `GET /api/videos/{id}` и проверить массив `qualities`.
8. Проверить, что обложка сохранилась сразу при создании видео.

Поля для создания видео:

- `title` - название
- `description` - описание
- `category_name` - имя категории, необязательное поле
- `video` - файл видео
- `thumbnail` - файл обложки

При обновлении видео через `PUT /api/videos/{id}` можно менять только `title`, `description`, `category_name` и файл `thumbnail`. Сам видеофайл после создания не меняется.

## Проверка перед демо

```bash
gofmt -w .
go build ./...
```
