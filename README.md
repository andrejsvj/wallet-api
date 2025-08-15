## Быстрый старт

Перед началом убедитесь, что у вас есть:

- Go
- Docker

Создайте файл `.env` с переменными окружения:

```env
# Конфигурация HTTP сервера
APP_ENV=local
HTTP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=wallet_user
DB_PASSWORD=wallet_password
DB_NAME=wallet_db
MAX_CONNECTIONS=100

POSTGRES_DB=wallet_db
POSTGRES_USER=wallet_user
POSTGRES_PASSWORD=wallet_password
```

## Основные команды

### ▶️ Локальный запуск

```bash
make local
```

Запускает сервер из директории `cmd`.

### 🛠 Сборка бинарника

```bash
make build
```

Собирает бинарный файл из кода в `cmd`.

### 🗄 Миграции базы данных

#### Применить миграции

```bash
make migrate
```

Применяет миграции из `./migrations` с помощью `goose`.

#### Создать новую миграцию

```bash
make new-migration name=название_миграции
```

Создает новую миграцию с указанным именем.

### ✅ Тестирование

```bash
make test
```

Запускает юнит-тесты с флагами `-v`, `-short` и `-race`.

## API Endpoints

### POST `/api/v1/wallet`
Операции с кошельком (пополнение/снятие)

### GET `/api/v1/wallets/{walletId}`
Получить информацию о кошельке


### 🚀 Docker

```bash
make start
```

Запускает приложение через Docker Compose.

```bash
make stop
```

Останавливает Docker Compose.

```bash
make docker-build
```

Собирает Docker образ.

---