---
title: جایگزین‌ها
description: جایگزین‌های Xray Checker
---

### یکپارچه‌سازی Node Exporter

ترکیب با متریک‌های node-exporter:

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

استفاده با حالت run-once:

```bash
curl -fsS --retry 3 https://hc-ping.com/your-uuid-here && \
./xray-checker --subscription-url=... --run-once
```

### یکپارچه‌سازی صفحه وضعیت

ارائه نقاط پایانی وضعیت به ارائه‌دهندگان صفحه وضعیت:

- BetterStack
- UptimeRobot
- StatusCake

فرمت URL مثال:

```
https://your-server:2112/config/a1b2c3d4e5f67890
```

### نظارت سفارشی

مثال‌های HTTP API برای نظارت سفارشی:

بررسی همه پروکسی‌ها:

```bash
curl -s localhost:2112/metrics | grep xray_proxy_status
```

بررسی پروکسی خاص:

```bash
curl -s localhost:2112/config/a1b2c3d4e5f67890
```

تجزیه متریک‌ها با jq:

```bash
curl -s localhost:2112/metrics | grep xray_proxy_status | \
  jq -R 'split(" ") | {name: (.[0] | split("{")[1] | split("}")[0]), value: .[1]}'
```
