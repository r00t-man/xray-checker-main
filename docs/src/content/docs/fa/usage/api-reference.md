---
title: مرجع API
description: مرجع API Xray Checker
---

## نمای کلی

Xray Checker هم endpointهای HTTP عمومی و هم محافظت شده ارائه می‌دهد. endpointهای محافظت شده وقتی `METRICS_PROTECTED=true` تنظیم شده باشد نیاز به احراز هویت خواهند داشت.

## نقاط پایانی عمومی

این نقاط پایانی همیشه بدون احراز هویت قابل دسترسی هستند.

### بررسی سلامت

```http
GET /health
```

endpoint ساده بررسی سلامت.

**پاسخ:** `200 OK` با بدنه `OK`

### وضعیت عمومی پروکسی

```http
GET /api/v1/public/proxies
```

وضعیت پروکسی را بدون داده‌های حساس (بدون IP/پورت سرور) برمی‌گرداند. توسط رابط وب برای به‌روزرسانی خودکار استفاده می‌شود.

**پاسخ:**
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

## نقاط پایانی محافظت شده

وقتی `METRICS_PROTECTED=true`، این نقاط پایانی نیاز به احراز هویت پایه دارند.

### رابط وب

```http
GET /
```

داشبورد HTML با نمای کلی وضعیت پروکسی، جستجو، فیلتر، مرتب‌سازی و به‌روزرسانی خودکار.

### متریک‌های Prometheus

```http
GET /metrics
```

نقطه پایانی متریک‌های Prometheus.

**مثال متریک‌ها:**
```text
# HELP xray_proxy_status Status of proxy connection (1: success, 0: failure)
# TYPE xray_proxy_status gauge
xray_proxy_status{protocol="vless",address="example.com:443",name="proxy1",sub_name="Premium VPN"} 1

# HELP xray_proxy_latency_ms Latency of proxy connection in milliseconds
# TYPE xray_proxy_latency_ms gauge
xray_proxy_latency_ms{protocol="vless",address="example.com:443",name="proxy1",sub_name="Premium VPN"} 156
```

### وضعیت پروکسی تکی

```http
GET /config/{stableId}
```

نقطه پایانی وضعیت برای پروکسی فردی، ایده‌آل برای نظارت بر آپتایم.

**پارامترها:**
- `stableId`: شناسه هش پایدار ۱۶ کاراکتری برای پروکسی

**پاسخ:**
- `200 OK` با بدنه `OK` اگر پروکسی کار می‌کند
- `503 Service Unavailable` با بدنه `Failed` اگر پروکسی کار نمی‌کند

:::tip[پیدا کردن Stable ID‌ها]
Stable ID‌ها در URL رابط وب هنگام کلیک روی نام پروکسی قابل مشاهده هستند، یا از طریق نقطه پایانی `/api/v1/proxies`.
:::

### لیست همه پروکسی‌ها

```http
GET /api/v1/proxies
```

اطلاعات کامل همه پروکسی‌ها را برمی‌گرداند.

**پاسخ:**
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

### دریافت پروکسی با ID

```http
GET /api/v1/proxies/{stableId}
```

اطلاعات یک پروکسی خاص را برمی‌گرداند.

**پاسخ:** همان ساختار آیتم تکی از `/api/v1/proxies`

### وضعیت سیستم

```http
GET /api/v1/status
```

آمار خلاصه را برمی‌گرداند.

**پاسخ:**
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

### پیکربندی

```http
GET /api/v1/config
```

پیکربندی فعلی checker را برمی‌گرداند.

**پاسخ:**
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

### اطلاعات سیستم

```http
GET /api/v1/system/info
```

اطلاعات نسخه و uptime را برمی‌گرداند.

**پاسخ:**
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

### IP فعلی

```http
GET /api/v1/system/ip
```

آدرس IP فعلی شناسایی شده سرور را برمی‌گرداند.

**پاسخ:**
```json
{
  "success": true,
  "data": {
    "ip": "203.0.113.1"
  }
}
```

### مستندات API

```http
GET /api/v1/docs
```

رابط کاربری Swagger برای مستندات تعاملی API.

```http
GET /api/v1/openapi.yaml
```

فایل مشخصات OpenAPI.

## احراز هویت

وقتی فعال باشد (`METRICS_PROTECTED=true`)، نقاط پایانی محافظت شده نیاز به احراز هویت پایه (نام کاربری و رمز عبور) دارند:

```bash
curl -u username:password http://localhost:2112/metrics
```

**توجه:** نقاط پایانی عمومی (`/health`، `/api/v1/public/proxies`) هرگز نیاز به احراز هویت ندارند.

## مثال‌های یکپارچه‌سازی

### Uptime Kuma

```bash
# آدرس نظارت (از stableId از رابط وب یا API استفاده کنید)
http://localhost:2112/config/a1b2c3d4e5f67890

# با احراز هویت
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

## پاسخ‌های خطا

تمام نقاط پایانی API فرمت خطای یکسانی برمی‌گردانند:

```json
{
  "success": false,
  "error": "Error message"
}
```

کدهای وضعیت HTTP:
- `200 OK`: درخواست موفق
- `400 Bad Request`: پارامترهای نامعتبر
- `401 Unauthorized`: احراز هویت مورد نیاز است
- `404 Not Found`: منبع پیدا نشد
- `500 Internal Server Error`: خطای سرور
- `503 Service Unavailable`: بررسی پروکسی ناموفق
