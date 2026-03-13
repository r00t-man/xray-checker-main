---
title: متریک‌ها
description: گزینه‌ها و مثال‌های متریک‌ها
---

Xray Checker دو متریک Prometheus برای نظارت بر وضعیت و عملکرد پروکسی ارائه می‌دهد. برای دستورالعمل‌های راه‌اندازی دقیق، [یکپارچه‌سازی Prometheus](/fa/integrations/prometheus) را ببینید.

برای به تصویر کشیدن متریک‌ها، توصیه می‌کنیم از [Grafana](/fa/integrations/grafana) استفاده کنید.

### xray_proxy_status

متریک وضعیت که در دسترس بودن پروکسی را نشان می‌دهد:

- نوع: Gauge
- مقادیر: ۱ (کار می‌کند) یا ۰ (ناموفق)
- برچسب‌ها:
  - `protocol`: پروتکل پروکسی (vless/vmess/trojan/shadowsocks)
  - `address`: آدرس و پورت سرور
  - `name`: نام پیکربندی پروکسی
  - `sub_name`: نام اشتراک (از فرگمنت URL یا هدر profile-title)
  - `instance`: نام نمونه (اگر پیکربندی شده باشد)

:::tip
برای راه‌اندازی برچسب‌گذاری نمونه [پیکربندی پیشرفته](/fa/configuration/advanced-conf#instance-labeling) را ببینید.
:::

مثال:

```text
# HELP xray_proxy_status Status of proxy connection (1: success, 0: failure)
# TYPE xray_proxy_status gauge
xray_proxy_status{protocol="vless",address="example.com:443",name="proxy1",sub_name="Premium VPN",instance="dc1"} 1
```

### xray_proxy_latency_ms

متریک تأخیر که زمان پاسخ اتصال را نشان می‌دهد:

- نوع: Gauge
- مقادیر: میلی‌ثانیه (۰ اگر ناموفق)
- برچسب‌ها: همانند xray_proxy_status

مثال:

```text
# HELP xray_proxy_latency_ms Latency of proxy connection in milliseconds
# TYPE xray_proxy_latency_ms gauge
xray_proxy_latency_ms{protocol="vless",address="example.com:443",name="proxy1",sub_name="Premium VPN",instance="dc1"} 156
```
