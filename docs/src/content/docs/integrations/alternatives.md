---
title: Alternative monitoring
description: Alternatives to Xray Checker
---

### Node Exporter Integration

Combine with node-exporter metrics:

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

Use with run-once mode:

```bash
curl -fsS --retry 3 https://hc-ping.com/your-uuid-here && \
./xray-checker --subscription-url=... --run-once
```

### Status Page Integration

Expose status endpoints to status page providers:

- BetterStack
- UptimeRobot
- StatusCake

Example URL format:

```
https://your-server:2112/config/a1b2c3d4e5f67890
```

### Custom Monitoring

HTTP API examples for custom monitoring:

Check all proxies:

```bash
curl -s localhost:2112/metrics | grep xray_proxy_status
```

Check specific proxy:

```bash
curl -s localhost:2112/config/a1b2c3d4e5f67890
```

Parse metrics with jq:

```bash
curl -s localhost:2112/metrics | grep xray_proxy_status | \
  jq -R 'split(" ") | {name: (.[0] | split("{")[1] | split("}")[0]), value: .[1]}'
```
