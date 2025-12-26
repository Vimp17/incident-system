# Incident System API

Система мониторинга опасных зон с мобильным приложением и веб-порталом на Django.

## 🚀 Быстрый старт

### Требования
- Docker & Docker Compose
- Go 1.24+
- PostgreSQL 15
- Redis

### Запуск через Docker (рекомендуется)

```bash
# Клонируйте репозиторий
git clone https://github.com/Vimp17/incident-system
cd incident-system
```

# Настройте окружение
```bash
cp .env.example .env
```
# Отредактируйте .env при необходимости

# Запустите все сервисы
```bash
docker-compose up -d
```
# Примените миграции
```bash
docker-compose exec postgres psql -U postgres -d incident_system -f /docker-entrypoint-initdb.d/001_init.sql
```
# Проверьте работоспособность
```bash
curl http://localhost:8080/api/v1/system/health
```

🛠 Технический стек
Backend: Go 1.24+ (Clean Architecture)

База данных: PostgreSQL 15

Кэш/Очередь: Redis

API Gateway: Gin Web Framework

Вебхуки: Асинхронная отправка с Retry

Контейнеризация: Docker & Docker Compose

📡 API Endpoints
Публичные эндпоинты
Health Check
```bash
GET /api/v1/system/health
```
Проверка локации
```bash
POST /api/v1/location/check
Content-Type: application/json

{
  "user_id": "user_123",
  "latitude": 55.7558,
  "longitude": 37.6173
}
```
Защищенные эндпоинты (требуют X-API-Key)
CRUD для инцидентов
Создать инцидент:

```bash
POST /api/v1/incidents
X-API-Key: operator-key-secure-change-me

{
  "user_id": "operator_1",
  "latitude": 55.7558,
  "longitude": 37.6173,
  "title": "Пожар в центре",
  "description": "Крупный пожар в бизнес-центре",
  "severity": "high",
  "radius": 1000
}
```
Получить список инцидентов:

```bash
GET /api/v1/incidents?page=1&limit=10&active_only=true
X-API-Key: operator-key-secure-change-me
```
Получить инцидент по ID:

```bash
GET /api/v1/incidents/{id}
X-API-Key: operator-key-secure-change-me
```
Обновить инцидент:

```bash
PUT /api/v1/incidents/{id}
X-API-Key: operator-key-secure-change-me

{
  "title": "Обновленное название",
  "severity": "medium",
  "active": false
}
```
Удалить инцидент:

```bash
DELETE /api/v1/incidents/{id}
X-API-Key: operator-key-secure-change-me
```
Статистика
```bash
GET /api/v1/incidents/stats?minutes=60
X-API-Key: operator-key-secure-change-me
```

## 🔍 Автоматические скрипты проверки

В папке `scripts/` находятся скрипты для автоматической проверки работоспособности системы:

### Для Windows:
```powershell
# Основной PowerShell скрипт
.\scripts\check-health.ps1
```
### Для Linux
```bash
.\scripts\check-health.sh
```
