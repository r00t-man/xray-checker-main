---
title: Prometheus
description: Варианты и примеры интеграции с Prometheus
---

### Прямой сбор метрик

Базовая конфигурация prometheus.yml:

```yaml
scrape_configs:
  - job_name: "xray-checker"
    metrics_path: "/metrics"
    static_configs:
      - targets: ["localhost:2112"]
    scrape_interval: 1m
```

С аутентификацией:

```yaml
scrape_configs:
  - job_name: "xray-checker"
    metrics_path: "/metrics"
    basic_auth:
      username: "metricsUser"
      password: "MetricsVeryHardPassword"
    static_configs:
      - targets: ["localhost:2112"]
```

### Интеграция с Pushgateway

Конфигурация Prometheus для Pushgateway:

```yaml
scrape_configs:
  - job_name: "pushgateway"
    honor_labels: true
    static_configs:
      - targets: ["pushgateway:9091"]
```

Конфигурация Xray Checker:

```bash
METRICS_PUSH_URL="http://user:password@pushgateway:9091/metrics/job/xray-checker"
```
`user:password@`: Используются только в случае когда включен basic auth
`/metrics/job/{yourname}`: Эндпоинт для метрик, подробнее читатйте в документации pushgateway
