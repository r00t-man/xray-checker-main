---
title: راه‌اندازی Prometheus
description: گزینه‌ها و مثال‌های یکپارچه‌سازی Prometheus
---

### Scraping مستقیم

پیکربندی پایه prometheus.yml:

```yaml
scrape_configs:
  - job_name: "xray-checker"
    metrics_path: "/metrics"
    static_configs:
      - targets: ["localhost:2112"]
    scrape_interval: 1m
```

با احراز هویت:

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

### یکپارچه‌سازی Pushgateway

پیکربندی Prometheus برای Pushgateway:

```yaml
scrape_configs:
  - job_name: "pushgateway"
    honor_labels: true
    static_configs:
      - targets: ["pushgateway:9091"]
```

پیکربندی Xray Checker:

```bash
METRICS_PUSH_URL="http://user:password@pushgateway:9091"
```
