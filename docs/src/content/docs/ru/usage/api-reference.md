---
title: Справочник по API
description: Справочник по API Xray Checker
---

## Обзор

Xray Checker предоставляет публичные и защищённые HTTP-эндпоинты. Защищённые эндпоинты требуют аутентификации при `METRICS_PROTECTED=true`.

## Публичные эндпоинты

Эти эндпоинты всегда доступны без аутентификации.

### Проверка работоспособности

```http
GET /health
```

Простой эндпоинт проверки работоспособности.

**Ответ:** `200 OK` с телом `OK`

### Публичный статус прокси

```http
GET /api/v1/public/proxies
```

Возвращает статус прокси без конфиденциальных данных (без IP/портов серверов). Используется веб-интерфейсом для автообновления.

**Ответ:**
```json
{
  "success": true,
  "data": [
    {
      "stableId": "a1b2c3d4e5f67890",
      "name": "US-Server-1",
      "online": true,
      "latencyMs": 150
    }
  ]
}
```

## Защищённые эндпоинты

При `METRICS_PROTECTED=true` эти эндпоинты требуют Basic Authentication.

### Веб-интерфейс

```http
GET /
```

HTML-панель с обзором статуса прокси, поиском, фильтрацией, сортировкой и автообновлением.

### Метрики Prometheus

```http
GET /metrics
```

Эндпоинт метрик Prometheus.

**Пример метрик:**
```text
# HELP xray_proxy_status Статус прокси-соединения (1: успешно, 0: неудача)
# TYPE xray_proxy_status gauge
xray_proxy_status{protocol="vless",address="example.com:443",name="proxy1",sub_name="Premium VPN"} 1

# HELP xray_proxy_latency_ms Задержка прокси-соединения в миллисекундах
# TYPE xray_proxy_latency_ms gauge
xray_proxy_latency_ms{protocol="vless",address="example.com:443",name="proxy1",sub_name="Premium VPN"} 156
```

### Статус отдельного прокси

```http
GET /config/{stableId}
```

Эндпоинт статуса отдельного прокси, идеально подходит для мониторинга доступности.

**Параметры:**
- `stableId`: 16-символьный стабильный идентификатор-хеш прокси

**Ответ:**
- `200 OK` с телом `OK` если прокси работает
- `503 Service Unavailable` с телом `Failed` если прокси не работает

:::tip[Поиск Stable ID]
Stable ID видны в URL веб-интерфейса при клике на имя прокси, или через эндпоинт `/api/v1/proxies`.
:::

### Список всех прокси

```http
GET /api/v1/proxies
```

Возвращает полную информацию о всех прокси.

**Ответ:**
```json
{
  "success": true,
  "data": [
    {
      "index": 0,
      "stableId": "a1b2c3d4e5f67890",
      "name": "US-Server-1",
      "subName": "Premium VPN",
      "server": "192.168.1.1",
      "port": 443,
      "protocol": "vless",
      "proxyPort": 10000,
      "online": true,
      "latencyMs": 150
    }
  ]
}
```

### Получить прокси по ID

```http
GET /api/v1/proxies/{stableId}
```

Возвращает информацию о конкретном прокси.

**Ответ:** Та же структура, что и для одного элемента из `/api/v1/proxies`

### Статус системы

```http
GET /api/v1/status
```

Возвращает сводную статистику.

**Ответ:**
```json
{
  "success": true,
  "data": {
    "total": 10,
    "online": 8,
    "offline": 2,
    "avgLatencyMs": 200
  }
}
```

### Конфигурация

```http
GET /api/v1/config
```

Возвращает текущую конфигурацию чекера.

**Ответ:**
```json
{
  "success": true,
  "data": {
    "checkInterval": 300,
    "checkMethod": "ip",
    "timeout": 30,
    "startPort": 10000,
    "subscriptionUpdate": true,
    "subscriptionUpdateInterval": 300,
    "simulateLatency": true,
    "subscriptionNames": ["Premium VPN", "Basic VPN"]
  }
}
```

### Информация о системе

```http
GET /api/v1/system/info
```

Возвращает информацию о версии и времени работы.

**Ответ:**
```json
{
  "success": true,
  "data": {
    "version": "1.0.0",
    "uptime": "1h 30m 45s",
    "uptimeSec": 5445,
    "instance": "prod-1"
  }
}
```

### Текущий IP

```http
GET /api/v1/system/ip
```

Возвращает текущий определённый IP-адрес сервера.

**Ответ:**
```json
{
  "success": true,
  "data": {
    "ip": "203.0.113.1"
  }
}
```

### Документация API

```http
GET /api/v1/docs
```

Swagger UI для интерактивной документации API.

```http
GET /api/v1/openapi.yaml
```

Файл спецификации OpenAPI.

## Аутентификация

При включении (`METRICS_PROTECTED=true`) защищённые эндпоинты требуют Basic Authentication:

```bash
curl -u username:password http://localhost:2112/metrics
```

**Примечание:** Публичные эндпоинты (`/health`, `/api/v1/public/proxies`) никогда не требуют аутентификации.

## Примеры интеграции

### Uptime Kuma

```bash
# URL монитора (используйте stableId из веб-интерфейса или API)
http://localhost:2112/config/a1b2c3d4e5f67890

# С аутентификацией
http://username:password@localhost:2112/config/a1b2c3d4e5f67890
```

### Prometheus

```yaml
scrape_configs:
  - job_name: "xray-checker"
    metrics_path: "/metrics"
    basic_auth:
      username: "username"
      password: "password"
    static_configs:
      - targets: ["localhost:2112"]
```

## Ответы об ошибках

Все API-эндпоинты возвращают единообразный формат ошибок:

```json
{
  "success": false,
  "error": "Сообщение об ошибке"
}
```

HTTP-коды статуса:
- `200 OK`: Запрос успешен
- `400 Bad Request`: Неверные параметры
- `401 Unauthorized`: Требуется аутентификация
- `404 Not Found`: Ресурс не найден
- `500 Internal Server Error`: Ошибка сервера
- `503 Service Unavailable`: Проверка прокси не удалась
