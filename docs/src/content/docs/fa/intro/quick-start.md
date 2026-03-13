---
title: راه‌اندازی سریع
description: راه‌اندازی سریع Xray Checker
---

Xray Checker را در چند دقیقه با این مراحل ساده راه‌اندازی کنید.

## پیش‌نیازها

- لینک اشتراک برای پروکسی‌های شما
- Docker (اختیاری، برای استقرار در کانتینر)
- Prometheus (اختیاری، برای جمع‌آوری متریک‌ها)

## راه‌اندازی ۵ دقیقه‌ای

### با استفاده از Docker (توصیه شده)

1. دریافت ایمیج:

```bash
docker pull kutovoys/xray-checker
```

2. اجرا با پیکربندی پایه:

```bash
docker run -d \
  -e SUBSCRIPTION_URL=https://your-subscription-url/sub \
  -p 2112:2112 \
  kutovoys/xray-checker
```

3. بررسی وضعیت:

```bash
curl http://localhost:2112/health
```

### استفاده از فایل باینری

1. دانلود آخرین نسخه:

```bash
# Linux amd64
curl -sL -o - $(curl -s https://api.github.com/repos/kutovoys/xray-checker/releases/latest | grep "browser_download_url.*linux-amd64.tar.gz" | cut -d'"' -f4) | tar -xz
chmod +x xray-checker

# Linux arm64
curl -sL -o - $(curl -s https://api.github.com/repos/kutovoys/xray-checker/releases/latest | grep "browser_download_url.*linux-arm64.tar.gz" | cut -d'"' -f4) | tar -xz
chmod +x xray-checker
```

2. اجرا با پیکربندی پایه:

```bash
./xray-checker --subscription-url=https://your-subscription-url/sub
```

## تأیید نصب

1. باز کردن رابط وب:

   - به آدرس `http://localhost:2112` بروید
   - باید داشبورد با وضعیت پروکسی را ببینید

2. بررسی متریک‌ها:

   - به آدرس `http://localhost:2112/metrics` بروید
   - باید متریک‌های Prometheus را ببینید

3. تأیید وضعیت پروکسی:
   - روی هر لینک پروکسی در رابط وب کلیک کنید
   - پاسخ نقطه پایانی وضعیت را بررسی کنید

## مراحل بعدی

1. پیکربندی Prometheus:

```yaml
scrape_configs:
  - job_name: "xray-checker"
    static_configs:
      - targets: ["localhost:2112"]
```

2. راه‌اندازی Uptime Kuma:

   - یک مانیتور جدید اضافه کنید
   - برای هر پروکسی از endpoint مخصوص آن استفاده کنید
   - تنظیمات هشدار و نوتیفیکیشن را تکمیل کنید

3. سفارشی‌سازی پیکربندی:
   - فواصل بررسی را تنظیم کنید
   - پیکربندی احراز هویت را انجام دهید
   - تنظیمات ارسال متریک را تغییر دهید

## دستورات رایج

بررسی نسخه:

```bash
./xray-checker --version
```

اجرای بررسی تکی:

```bash
./xray-checker --subscription-url=https://your-sub-url --run-once
```

فعال‌سازی احراز هویت:

```bash
./xray-checker --subscription-url=https://your-sub-url \
  --metrics-protected=true \
  --metrics-username=user \
  --metrics-password=pass
```

## عیب‌یابی

1. بررسی وضعیت سرویس:

```bash
curl http://localhost:2112/health
```

2. مشاهده لاگ‌ها:

```bash
docker logs xray-checker
```

3. تأیید عملکرد متریک‌ها:

```bash
curl http://localhost:2112/metrics
```

## نیاز به کمک دارید؟

- مستندات را بررسی کنید
- یک issue در GitHub باز کنید
- به بحث‌های انجمن بپیوندید
