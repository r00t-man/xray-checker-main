---
title: Альтернативный мониторинг
description: Альтернативные системы мониторинга
---

### Интеграция с Node Exporter

Комбинирование с метриками node-exporter:

```yaml
scrape_configs:
  - job_name: "xray-checker"
    static_configs:
      - targets: ["localhost:2112"]
  - job_name: "node"
    static_configs:
      - targets: ["localhost:9100"]
```

### Healthchecks.io

Использование с режимом однократного запуска:

```bash
curl -fsS --retry 3 https://hc-ping.com/your-uuid-here && \
./xray-checker --subscription-url=... --run-once
```

### Интеграция со страницей статуса

Предоставление эндпоинтов статуса для провайдеров страниц статуса:

- BetterStack
- UptimeRobot
- StatusCake

Пример формата URL:

```
https://your-server:2112/config/a1b2c3d4e5f67890
```

### Пользовательский мониторинг

Примеры использования HTTP API для пользовательского мониторинга:

Проверка всех прокси:

```bash
curl -s localhost:2112/metrics | grep xray_proxy_status
```

Проверка конкретного прокси:

```bash
curl -s localhost:2112/config/a1b2c3d4e5f67890
```

Разбор метрик с помощью jq:

```bash
curl -s localhost:2112/metrics | grep xray_proxy_status | \
  jq -R 'split(" ") | {name: (.[0] | split("{")[1] | split("}")[0]), value: .[1]}'
```
