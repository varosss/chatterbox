# Chatterbox

**Chatterbox** — backend-часть мессенджера, реализованная на Go с использованием микросервисной архитектуры.

## Описание

Проект представляет собой серверную часть мессенджера с поддержкой:

* регистрации и авторизации пользователей
* управления чатами
* отправки и получения сообщений в реальном времени
* доставки уведомлений через WebSocket

## Архитектурные особенности

* Микросервисная архитектура
* Clean Architecture / DDD / Hexagonal подход
* Event-driven взаимодействие между сервисами
* JWT-аутентификация с RSA-подписью
* Асинхронная обработка событий через брокер сообщений

## Используемые технологии

* **Go**
* **PostgreSQL**
* **RabbitMQ**
* **Docker / Docker Compose**
* **WebSocket**
* **JWT (RSA)**
* **Prometheus + Grafana**

## Сервисы

* **User Service** — регистрация, авторизация, управление пользователями
* **Chat Service** — работа с чатами и сообщениями
* **Notification Service** — доставка уведомлений и WebSocket-соединения

## Развертывание

### 1. Клонирование репозитория

```bash
git clone https://github.com/varosss/chatterbox.git
cd chatterbox
```

### 2. Настройка переменных окружения

Каждый сервис использует собственный .env файл:

```
chat/.env
user/.env
notification/.env
```

Создайте их на основе .env.dist, если он присутствует в репозитории:

```bash
cp chat/.env.dist chat/.env
cp user/.env.dist user/.env
cp notification/.env.dist notification/.env
```

Установите переменные окружения в файлах .env

### 3. Сборка образов

```bash
make build.all
```

### 4. Запуск

```bash
docker compose up -d
```

Поздравляю, приложение запущено!
