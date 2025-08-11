# Cистемы обработки заказов с использованием Temporal.io

## Описание

Проект моделирует процесс обработки заказа с помощью Workflow и Activities, включает обработку ошибок, повторные попытки, логику отмены и простой web-интерфейс для взаимодействия.
---
images\GitHub_README.jpg
## Архитектура

- **API** (`cmd/api/main.go`):  HTTP сервер для взаимодействия с системой.
- **Worker** (`cmd/worker/main.go`):  Запускает Temporal worker, регистрирует workflow и activities.
- **Workflow** (`internal/workflows/order_workflow.go`):  Управляет жизненным циклом заказа, вызывает Activities для каждого шага, обрабатывает ошибки, поддерживает Query и Signal.
- **Activities** (`internal/activities/`):
  - `CheckInventoryActivity` - проверка наличия товара
  - `ProcessPaymentActivity` - обработка платежа
  - `NotifyCustomerActivity` - уведомление клиента
- **Модели** (`internal/models/`):  Описывают структуру заказа и его позиций.
- **Handlers** (`internal/handlers/`):  HTTP-обработчики для создания заказа, проверки статуса, отмены.
- **Web-интерфейс** (`web/index.html`):  Простой UI для создания заказа, проверки статуса и отмены.

---

## Запуск проекта

### 1. Установите зависимости

```bash
go mod tidy
```

### 2. Запустите Temporal 

```bash
docker-compose -f docker-compose/docker-compose.yml up -d
```


### 3. Запустите worker

```bash
go run ./cmd/worker/main.go
```

### 4. Запустите API сервер

```bash
go run ./cmd/api/main.go
```

### 5. Откройте web-интерфейс

Перейдите в браузере по адресу:  
```
http://localhost:8081/web
```
Temporal UI:
```
http://localhost:8080
```

---

## Основные фичи

- **Создание заказа** через web-интерфейс или API.
- **Пошаговая обработка заказа**: проверка склада, оплата, уведомление.
- **Обработка ошибок** и автоматические повторные попытки Activities.
- **Отмена заказа** через Signal (кнопка Cancel в UI).
- **Получение статуса заказа** через Query (отображается в UI).
- **Логирование** всех этапов процесса.

---

## Пример API

### Order Management
- `POST /api/orders/create` - создать заказ
- `GET /api/orders/{orderId}` - получить статус заказа
- `POST /api/orders/{orderId}/cancel` - отменить заказ

### Workflow Management
- `GET /api/workflows/{workflowId}/status` - получить статус workflow

### Health Checks
- `GET /health` - проверка здоровья API сервера

### Web Interface
- `GET /web/` - веб-интерфейс для управления заказами

---

## Требования

- Go 1.18+
- Docker (для Temporal)
- Temporal SDK Go
