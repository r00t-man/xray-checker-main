---
title: API Reference
description: API Reference of Xray Checker
---

## Overview

Xray Checker provides both public and protected HTTP endpoints. Protected endpoints require authentication when `METRICS_PROTECTED=true`.

## Public Endpoints

These endpoints are always accessible without authentication.

### Health Check

```http
GET /health
```

Simple health check endpoint.

**Response:** `200 OK` with body `OK`

### Public Proxy Status

```http
GET /api/v1/public/proxies
```

Returns proxy status without sensitive data (no server IPs/ports). Used by the web UI for auto-refresh.

**Response:**
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

## Protected Endpoints

When `METRICS_PROTECTED=true`, these endpoints require Basic Authentication.

### Web Interface

```http
GET /
```

HTML dashboard with proxy status overview, search, filtering, sorting, and auto-refresh.

### Prometheus Metrics

```http
GET /metrics
```

Prometheus metrics endpoint.

**Example metrics:**
```text
# HELP xray_proxy_status Status of proxy connection (1: success, 0: failure)
# TYPE xray_proxy_status gauge
xray_proxy_status{protocol="vless",address="example.com:443",name="proxy1",sub_name="Premium VPN"} 1

# HELP xray_proxy_latency_ms Latency of proxy connection in milliseconds
# TYPE xray_proxy_latency_ms gauge
xray_proxy_latency_ms{protocol="vless",address="example.com:443",name="proxy1",sub_name="Premium VPN"} 156
```

### Individual Proxy Status

```http
GET /config/{stableId}
```

Status endpoint for individual proxy, perfect for uptime monitoring.

**Parameters:**
- `stableId`: 16-character stable identifier hash for the proxy

**Response:**
- `200 OK` with body `OK` if proxy is working
- `503 Service Unavailable` with body `Failed` if proxy is not working

:::tip[Finding Stable IDs]
Stable IDs are visible in the web UI URL when clicking on a proxy name, or via the `/api/v1/proxies` endpoint.
:::

### List All Proxies

```http
GET /api/v1/proxies
```

Returns full information for all proxies.

**Response:**
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

### Get Proxy by ID

```http
GET /api/v1/proxies/{stableId}
```

Returns information for a specific proxy.

**Response:** Same structure as single item from `/api/v1/proxies`

### System Status

```http
GET /api/v1/status
```

Returns summary statistics.

**Response:**
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

### Configuration

```http
GET /api/v1/config
```

Returns current checker configuration.

**Response:**
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

### System Info

```http
GET /api/v1/system/info
```

Returns version and uptime information.

**Response:**
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

### Current IP

```http
GET /api/v1/system/ip
```

Returns the server's current detected IP address.

**Response:**
```json
{
  "success": true,
  "data": {
    "ip": "203.0.113.1"
  }
}
```

### API Documentation

```http
GET /api/v1/docs
```

Swagger UI for interactive API documentation.

```http
GET /api/v1/openapi.yaml
```

OpenAPI specification file.

## Authentication

When enabled (`METRICS_PROTECTED=true`), protected endpoints require Basic Authentication:

```bash
curl -u username:password http://localhost:2112/metrics
```

**Note:** Public endpoints (`/health`, `/api/v1/public/proxies`) never require authentication.

## Integration Examples

### Uptime Kuma

```bash
# Monitor URL (use stableId from web UI or API)
http://localhost:2112/config/a1b2c3d4e5f67890

# With authentication
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

## Error Responses

All API endpoints return consistent error format:

```json
{
  "success": false,
  "error": "Error message"
}
```

HTTP Status Codes:
- `200 OK`: Request successful
- `400 Bad Request`: Invalid parameters
- `401 Unauthorized`: Authentication required
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error
- `503 Service Unavailable`: Proxy check failed
