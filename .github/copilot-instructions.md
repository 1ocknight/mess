# Copilot / AI agent instructions for this repository

Кратко — что важно знать, чтобы быстро работать с репозиторием и писать исправления/фичи.

- **Большая картина:** проект — набор микросервисов на Go: `chat`, `profile`, `websocket`, плюс `front` (React). Общие утилиты и интеграции лежат в `shared`.
- **Основные границы ответственности:**
  - `chat`: CRUD по чатам и сообщениям, подсчёт непрочитанных, outbox-паттерн для фоновых задач. См. [chat/README.md](chat/README.md).
  - `profile`: CRUD профилей, работа с S3 (presigned URL) для аватаров, outbox для удаления из S3. См. [profile/README.md](profile/README.md).
  - `websocket`: поддерживает подключения, hub и доставку событий клиентам через общий `chan`. См. [websocket/README.md](websocket/README.md).

- **Интеграции и внешние зависимости:**
  - PostgreSQL (миграции в `*/migrations`, compose в [compose-local/docker-compose.yml](compose-local/docker-compose.yml) использует порт 5430).
  - Kafka (используется для сообщений/событий, см. `shared/kafkav2` и `compose-local/docker-compose.yml`).
  - MinIO/S3 (аватары), presigned URL паттерн (см. `profile/internal/adapter/avatar`).
  - Keycloak для аутентификации (локально на порту 7070 в compose).

- **Запуск и отладка локально:**
  - Быстрый стек (локально):
    ```bash
    docker compose -f compose-local/docker-compose.yml up --build
    ```
  - Запуск отдельного сервиса для разработки (пример из `chat`):
    ```bash
    cd chat
    go run ./cmd
    ```
  - Фронтенд:
    ```bash
    cd front
    npm install
    npm run dev
    ```

- **Тесты и e2e:** e2e находится в папке `e2e` (есть `cmd/main_test.go`) — запускать из контейнера/compose или локально, настраивая переменные окружения и `CONFIG`.

- **Ключевые реализационные паттерны, которые нужно соблюдать/учитывать:**
  - Outbox + фоновые воркеры: схемы удаления/репликации реализованы через таблицы outbox и воркеры, использующие транзакции + `FOR UPDATE SKIP LOCKED` (см. `internal/wokers` в `chat`/`profile`).
  - Версионирование сущностей: обновления данных реализованы через версионирование (см. `lastread` и миграции в `chat/migrations`).
  - Батчи и пагинация на уровне SQL: при больших объёмах данных код использует батчи и SQL-пагинацию (см. `storage` пакеты).
  - Hub для WebSocket: локальный репликационный hub хранит соединения и читает из общего `chan` (см. `websocket/internal/transport`).

- **Где смотреть примеры кода:**
  - Сервисы: [chat/cmd/main.go](chat/cmd/main.go), [profile/cmd/main.go](profile/cmd/main.go), [websocket/cmd/main.go](websocket/cmd/main.go).
  - Kafka: [shared/kafkav2/producer.go](shared/kafkav2/producer.go) (producer/consumer паттерн).
  - Общие адаптеры: `shared/postgres`, `shared/s3client`.

- **Стиль и конвенции проекта:**
  - Структура сервисов повторяющаяся: `cmd`, `config`, `internal/{ctxkey,domain,loglables,model,storage,transport,wokers}`.
  - Логирование с заранее определёнными лейблами (`loglables`).
  - Небольшие отклонения от общих названий: папка воркеров иногда называется `wokers` — обращайте внимание на опечатки при поиске.

- **PR/изменения — практические рекомендации для AI-ассистента:**
  - Минимизируйте изменения в API и сигнатурах без обсуждения; предпочтительны внутренние изменения и добавление тестов.
  - При правках, касающихся интеграций (Kafka, S3, Keycloak, Postgres), указывайте ожидаемые переменные окружения и зависимости `compose-local`.
  - Для изменений в фоновых воркерах соблюдайте транзакционный outbox-паттерн и `skip locked` логику — это критично для согласованности.

- **Что проверить при изменениях:**
  - Миграции: обновились ли SQL миграции в соответствующей папке `migrations`.
  - Конфиги: обновлены ли шаблоны в `compose-local/configs` (chat.yaml, profile.yaml, ws.yaml) при изменении переменных окружения.
  - Производительность: используйте батчи/пагинацию, если меняете код чтения большого объёма данных.

Если нужно, могу сократить инструкции, добавить шаблоны PR checklist или включить конкретные команды для интерграционного локального запуска (например, пошаговый сценарий для поднятия Keycloak+Postgres+Kafka+MinIO).
