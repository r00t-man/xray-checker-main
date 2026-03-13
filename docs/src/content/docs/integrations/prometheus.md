---
title: Prometheus
description: Prometheus integration options and examples
---

### Direct Scraping

Basic prometheus.yml configuration:

```yaml
scrape_configs:
  - job_name: "xray-checker"
    metrics_path: "/metrics"
    static_configs:
      - targets: ["localhost:2112"]
    scrape_interval: 1m
```

With authentication:

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

### Pushgateway Integration

Prometheus configuration for Pushgateway:

```yaml
scrape_configs:
  - job_name: "pushgateway"
    honor_labels: true
    static_configs:
      - targets: ["pushgateway:9091"]
```

Xray Checker configuration:

```bash
METRICS_PUSH_URL="http://user:password@pushgateway:9091"
```
