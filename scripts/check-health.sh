#!/bin/bash

# Скрипт проверки работоспособности Incident System API
# Запуск: ./scripts/check-health.sh

set -e

BASE_URL=${1:-"http://localhost:8080"}
API_KEY=${2:-"operator-key-secure-change-me"}

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Функция для выполнения запросов
make_request() {
    local name=$1
    local url=$2
    local method=$3
    local data=$4
    local headers=$5
    
    echo -e "${CYAN}[$name]${NC} Отправка запроса..."
    
    local curl_cmd="curl -s -X $method -H \"Content-Type: application/json\""
    
    # Добавляем заголовки если есть
    if [[ -n "$headers" ]]; then
        curl_cmd="$curl_cmd $headers"
    fi
    
    # Добавляем данные если есть
    if [[ -n "$data" ]]; then
        curl_cmd="$curl_cmd -d '$data'"
    fi
    
    # Добавляем URL и обработку ошибок
    curl_cmd="$curl_cmd $url"
    
    echo -e "${BLUE}Команда:${NC} $curl_cmd"
    
    if eval $curl_cmd; then
        echo -e "${GREEN}[$name] ✅ Успешно${NC}"
        return 0
    else
        echo -e "${RED}[$name] ❌ Ошибка${NC}"
        return 1
    fi
}

# Функция для проверки Docker
check_docker() {
    if ! command -v docker &> /dev/null; then
        echo -e "${YELLOW}Docker не установлен${NC}"
        return 1
    fi
    
    if ! docker ps &> /dev/null; then
        echo -e "${YELLOW}Docker не запущен${NC}"
        return 1
    fi
    
    echo -e "${GREEN}Docker доступен${NC}"
    return 0
}

echo -e "${CYAN}==========================================${NC}"
echo -e "${CYAN}Incident System Health Check${NC}"
echo -e "${CYAN}==========================================${NC}"
echo ""

# Проверка 1: Health Check
echo -e "${YELLOW}1. Проверка Health Check API...${NC}"
if make_request "Health" "$BASE_URL/api/v1/system/health" "GET"; then
    echo -e "${GREEN}✅ Сервер доступен${NC}"
else
    echo -e "${RED}❌ Сервер не отвечает${NC}"
    echo -e "${YELLOW}Запустите: go run cmd/server/main.go${NC}"
    exit 1
fi

echo ""

# Проверка 2: Создание инцидента
echo -e "${YELLOW}2. Создание тестового инцидента...${NC}"
INCIDENT_DATA='{
    "user_id": "health_check_operator",
    "latitude": 55.7558,
    "longitude": 37.6173,
    "title": "Тестовый инцидент для проверки",
    "description": "Создан автоматическим скриптом проверки",
    "severity": "medium",
    "radius": 500
}'

if make_request "Create Incident" \
    "$BASE_URL/api/v1/incidents" \
    "POST" \
    "$INCIDENT_DATA" \
    "-H \"X-API-Key: $API_KEY\""; then
    echo -e "${GREEN}✅ Инцидент создан${NC}"
else
    echo -e "${RED}❌ Не удалось создать инцидент${NC}"
fi

echo ""

# Проверка 3: Проверка локации
echo -e "${YELLOW}3. Проверка локации...${NC}"
LOCATION_DATA='{
    "user_id": "health_check_user",
    "latitude": 55.7558,
    "longitude": 37.6173
}'

make_request "Location Check" \
    "$BASE_URL/api/v1/location/check" \
    "POST" \
    "$LOCATION_DATA"

echo ""

# Проверка 4: Список инцидентов
echo -e "${YELLOW}4. Получение списка инцидентов...${NC}"
make_request "List Incidents" \
    "$BASE_URL/api/v1/incidents" \
    "GET" \
    "" \
    "-H \"X-API-Key: $API_KEY\""

echo ""

# Проверка 5: Статистика
echo -e "${YELLOW}5. Получение статистики...${NC}"
make_request "Statistics" \
    "$BASE_URL/api/v1/incidents/stats?minutes=5" \
    "GET" \
    "" \
    "-H \"X-API-Key: $API_KEY\""

echo ""

# Проверка 6: Docker контейнеры
echo -e "${YELLOW}6. Проверка Docker контейнеров...${NC}"
if check_docker; then
    echo -e "${GREEN}Запущенные контейнеры:${NC}"
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    
    # Проверка необходимых контейнеров
    REQUIRED_CONTAINERS=("incident-postgres" "incident-redis" "incident-webhook-mock")
    RUNNING_CONTAINERS=$(docker ps --format "{{.Names}}")
    
    for container in "${REQUIRED_CONTAINERS[@]}"; do
        if echo "$RUNNING_CONTAINERS" | grep -q "$container"; then
            echo -e "${GREEN}✅ $container запущен${NC}"
        else
            echo -e "${YELLOW}⚠️  $container не запущен${NC}"
        fi
    done
    
    # Проверка подключения к PostgreSQL
    echo -e "${BLUE}Проверка PostgreSQL...${NC}"
    if docker exec incident-postgres pg_isready -U postgres &> /dev/null; then
        echo -e "${GREEN}✅ PostgreSQL доступен${NC}"
    else
        echo -e "${RED}❌ PostgreSQL не доступен${NC}"
    fi
    
    # Проверка подключения к Redis
    echo -e "${BLUE}Проверка Redis...${NC}"
    if docker exec incident-redis redis-cli ping | grep -q "PONG"; then
        echo -e "${GREEN}✅ Redis доступен${NC}"
    else
        echo -e "${RED}❌ Redis не доступен${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Docker не доступен${NC}"
fi

echo ""
echo -e "${CYAN}==========================================${NC}"
echo -e "${CYAN}Проверка завершена!${NC}"
echo -e "${CYAN}==========================================${NC}"

# Рекомендации
echo ""
echo -e "${YELLOW}Рекомендации:${NC}"
echo -e "${BLUE}1. Запустить сервер:${NC} go run cmd/server/main.go"
echo -e "${BLUE}2. Запустить Docker сервисы:${NC} docker-compose up -d"
echo -e "${BLUE}3. Проверить логи:${NC} docker-compose logs -f"
echo -e "${BLUE}4. Тестовые запросы:${NC} curl примеры в README.md"