---
title: فرمت اشتراک
description: گزینه‌ها و مثال‌های فرمت اشتراک
---

Xray Checker از پنج فرمت مختلف برای پیکربندی پروکسی پشتیبانی می‌کند. از [متغیر محیطی](/fa/configuration/envs#subscription_url) `SUBSCRIPTION_URL` برای تنظیم روش بررسی استفاده کنید.

برای اطلاعات درباره نحوه تأیید پروکسی‌ها، [روش‌های بررسی](/fa/configuration/check-methods) را ببینید.

### ۱. آدرس اشتراک (پیش‌فرض)

آدرس اشتراک استاندارد که لیست کدگذاری شده Base64 از لینک‌های پروکسی را برمی‌گرداند.

مثال:

```bash
SUBSCRIPTION_URL=https://example.com/subscription
```

الزامات:

- آدرس HTTPS
- محتوای کدگذاری شده Base64 برگرداند
- محتوا آدرس‌های پروکسی با خط جدید (کاراکتر newline) از همدیگر جدا شده باشند
- از هدرهای استاندارد User-Agent پشتیبانی کند

هدرهای ارسالی:

```
Accept: */*
User-Agent: Xray-Checker
```

### ۲. رشته Base64

رشته مستقیم کدگذاری شده Base64 حاوی لینک‌های پیکربندی پروکسی.

مثال:

```bash
SUBSCRIPTION_URL=dmxlc3M6Ly91dWlkQGV4YW1wbGUuY29tOjQ0MyVlbmNyeXB0aW9uPW5vbmUmc2VjdXJpdHk9dGxzI3Byb3h5MQ==
```

فرمت محتوا (قبل از کدگذاری):

```
vless://uuid@example.com:443?encryption=none&security=tls#proxy1
trojan://password@example.com:443?security=tls#proxy2
vmess://base64encodedconfig
ss://base64encodedconfig
```

### ۳. فایل JSON V2Ray

فایل پیکربندی JSON تکی در فرمت V2Ray/Xray.

مثال:

```bash
SUBSCRIPTION_URL=file:///path/to/config.json
```

فرمت فایل:

```json
{
  "outbounds": [
    {
      "protocol": "vless",
      "settings": {
        "vnext": [
          {
            "address": "example.com",
            "port": 443,
            "users": [
              {
                "id": "uuid",
                "encryption": "none"
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "tcp",
        "security": "tls"
      }
    }
  ]
}
```

### ۴. آرایه JSON Xray (چند پیکربندی)

آرایه JSON حاوی چندین پیکربندی Xray با remarks. این فرمت برای صادر کردن پیکربندی‌ها از کلاینت‌های GUI یا مدیریت چندین پیکربندی نام‌گذاری شده در یک فایل مفید است.

مثال:

```bash
SUBSCRIPTION_URL=file:///path/to/configs.json
```

فرمت فایل:

```json
[
  {
    "remarks": "US Server 1",
    "outbounds": [
      {
        "protocol": "vless",
        "settings": {
          "vnext": [
            {
              "address": "us1.example.com",
              "port": 443,
              "users": [{ "id": "uuid-1", "encryption": "none" }]
            }
          ]
        },
        "streamSettings": { "network": "tcp", "security": "tls" }
      }
    ]
  },
  {
    "remarks": "EU Server 1",
    "outbounds": [
      {
        "protocol": "trojan",
        "settings": {
          "servers": [
            {
              "address": "eu1.example.com",
              "port": 443,
              "password": "password123"
            }
          ]
        },
        "streamSettings": { "network": "tcp", "security": "tls" }
      }
    ]
  }
]
```

فیلد `remarks` از هر پیکربندی به عنوان نام پروکسی در داشبورد استفاده خواهد شد.

### ۵. پوشه پیکربندی

دایرکتوری حاوی چندین فایل پیکربندی JSON V2Ray/Xray.

مثال:

```bash
SUBSCRIPTION_URL=folder:///path/to/configs
```

الزامات:

- دایرکتوری باید حاوی فایل‌های .json باشد
- هر فایل از فرمت JSON V2Ray پیروی می‌کند
- فایل‌ها به ترتیب الفبایی پردازش می‌شوند
- فایل‌های نامعتبر با هشدار نادیده گرفته می‌شوند
