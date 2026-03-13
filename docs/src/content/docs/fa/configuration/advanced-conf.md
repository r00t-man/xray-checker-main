---
title: پیکربندی پیشرفته
description: گزینه‌های پیکربندی پیشرفته
---

## پیکربندی پیشرفته

### سرویس‌های سفارشی بررسی IP

می‌توانید از سرویس‌های جایگزین بررسی IP استفاده کنید (برای جزئیات [روش‌های بررسی](/fa/configuration/check-methods) را ببینید):

- `http://ip.sb`
- `https://api64.ipify.org`
- `http://ifconfig.me`

مثال:

```bash
PROXY_IP_CHECK_URL=http://ip.sb
```

### آدرس‌های سفارشی بررسی وضعیت

آدرس‌های جایگزین برای بررسی وضعیت (برای جزئیات [روش‌های بررسی](/fa/configuration/check-methods) را ببینید):

- `http://www.gstatic.com/generate_204`
- `http://www.qualcomm.cn/generate_204`
- `http://cp.cloudflare.com/generate_204`

مثال:

```bash
PROXY_STATUS_CHECK_URL=http://www.gstatic.com/generate_204
```

### پیکربندی امنیتی

فعال‌سازی احراز هویت برای نقاط پایانی حساس:

```bash
METRICS_PROTECTED=true
METRICS_USERNAME=custom_user
METRICS_PASSWORD=secure_password
```

### برچسب‌گذاری نمونه

اضافه کردن برچسب‌های نمونه برای راه‌اندازی‌های توزیع شده:

```bash
METRICS_INSTANCE=datacenter-1
```

### فواصل به‌روزرسانی

سفارشی‌سازی فواصل بررسی و به‌روزرسانی:

```bash
# بررسی هر دقیقه
PROXY_CHECK_INTERVAL=60

# به‌روزرسانی اشتراک هر ساعت
SUBSCRIPTION_UPDATE_INTERVAL=3600
```

### پیکربندی لاگ

تنظیم لاگ Xray Core:

```bash
# فعال‌سازی لاگ اشکال‌زدایی
XRAY_LOG_LEVEL=debug

# غیرفعال‌سازی لاگ
XRAY_LOG_LEVEL=none
```

### پیکربندی پورت

سفارشی‌سازی محدوده پورت‌ها:

```bash
# شروع پورت‌های SOCKS5 از 20000
XRAY_START_PORT=20000

# تغییر پورت متریک
METRICS_PORT=9090
```

### پیکربندی برای دامنه steal-from-yourself

شما دامنه خودتان را دارید، your-domain.com، با یک وب‌سایت در حال اجرا روی آن،
و می‌خواهید نظارت را در `your-domain.com/xray/monitor` نمایش دهید.

Xray Checker را روی همان سروری که وب‌سایت شما میزبانی می‌شود اجرا کنید.
پارامتر `-p 127.0.0.1:2112:2112` تضمین می‌کند که دسترسی مستقیم
به آن فقط از خود سرور امکان‌پذیر است:

:::caution
اگر رابط وب به صورت عمومی قابل دسترسی است، توصیه می‌شود از basic auth برای محافظت استفاده کنید.
می‌توانید این را با استفاده از متغیرهای محیطی زیر فعال کنید:
`METRICS_PROTECTED`، `METRICS_USERNAME`، `METRICS_PASSWORD`.
:::

```bash
docker run -d \
  -e SUBSCRIPTION_URL=https://your-subscription-url/sub \
  -p 127.0.0.1:2112:2112 \
  -e METRICS_BASE_PATH=/xray/monitor \
  -e METRICS_PROTECTED=true \
  -e METRICS_USERNAME=custom_user \
  -e METRICS_PASSWORD=custom_pass \
  kutovoys/xray-checker
```

فایل پیکربندی nginx را باز کنید (`sudo nano /etc/nginx/your-domain.com`)، بخش اصلی را پیدا کنید:

```
server {
    root /var/www/your-domain.com/html;

    index index.html;
    server_name your-stealing-domain.com;
    ...
}
```

و ۳ location جدید را در آنجا قرار دهید:

```config

    # مدیریت /xray/monitor بدون اسلش انتهایی
    location = /xray/monitor {
        return 301 https://$host$request_uri/;
    }

    # مدیریت /xray/monitor/ - ریدایرکت به پورت داکر xray checker
    location /xray/monitor/ {
        proxy_pass http://127.0.0.1:2112/xray/monitor/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
```

سپس nginx را بررسی و بارگذاری مجدد کنید:

```bash
sudo nginx -t
sudo systemctl reload nginx
```

و در دسترس بودن نظارت را بررسی کنید:

```bash
 curl -I -L https://your-domain.com/xray/monitor
```
