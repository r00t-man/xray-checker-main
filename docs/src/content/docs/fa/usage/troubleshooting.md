---
title: عیب‌یابی
description: عیب‌یابی
tableOfContents:
  minHeadingLevel: 2
  maxHeadingLevel: 4
---

## مشکلات رایج

### مشکلات اشتراک

#### پاسخ نامعتبر لینک اشتراک

```
error parsing subscription: error getting subscription: unexpected status code: 403
```

**علل احتمالی:**

- URL نادرست است
- URL دیگر معتبر نیست
- سرور user agent Xray Checker را مسدود می‌کند

**راه‌حل‌ها:**

1. URL اشتراک را مجددا بررسی کنید
2. بررسی کنید URL هنوز فعال باشد
3. با ارائه‌دهنده اشتراک تماس بگیرید
4. به جای URL مستقیماً از فرمت Base64 استفاده کنید

#### ناموفقیت دیکد Base64

```
error decoding Base64: illegal base64 data...
```

**علل احتمالی:**

- کدگذاری Base64 نامعتبر
- Base64 URL-safe در مقابل استاندارد
- فضای خالی یا کاراکتر newline اضافی

**راه‌حل‌ها:**

1. مطمئن شوید رشته Base64 تمیز و بدون فضای خالی است
2. اگر کد Base64 استاندارد ناموفق بود دیکد URL-safe Base64 را امتحان کنید
3. بررسی کنید که محتوا نیاز به دیکد چندباره دارد یا نه

### مشکلات بررسی پروکسی

#### اجرا روی سرور پروکسی

وقتی Xray Checker را روی همان سروری که پروکسی‌های شما میزبانی می‌شوند اجرا می‌کنید، **باید** از روش بررسی `status` به جای روش پیش‌فرض `ip` استفاده کنید.

**چرا:**

- روش بررسی `ip` IP شما را با و بدون پروکسی مقایسه می‌کند
- وقتی روی سرور پروکسی اجرا می‌شود، هر دو IP یکسان خواهند بود
- این باعث منفی‌های کاذب می‌شود - پروکسی‌های کارا به عنوان ناموفق گزارش می‌شوند

**راه‌حل:**

```bash
# در محیط
PROXY_CHECK_METHOD=status
PROXY_STATUS_CHECK_URL=http://cp.cloudflare.com/generate_204

# یا از طریق CLI
--proxy-check-method=status --proxy-status-check-url="http://cp.cloudflare.com/generate_204"
```

#### ناموفقیت همه پروکسی‌ها

```
Warning: error parsing proxy URL: connection refused
```

**علل احتمالی:**

- مشکلات اتصال شبکه
- فایروال اتصالات را مسدود می‌کند
- سرویس بررسی IP در دسترس نیست

**راه‌حل‌ها:**

1. اتصال شبکه را بررسی کنید
2. قوانین فایروال را چک کنید
3. روش بررسی جایگزین را امتحان کنید:
   ```bash
   PROXY_CHECK_METHOD=status
   ```
4. از سرویس بررسی IP جایگزین استفاده کنید:
   ```bash
   PROXY_IP_CHECK_URL=http://ip.sb
   ```

#### تأخیر بالا یا timeout

```
Warning: error getting current IP: timeout
```

**علل احتمالی:**

- اتصال شبکه کند
- سرویس بررسی IP کند
- timeout پروکسی خیلی کم

**راه‌حل‌ها:**

1. timeout را افزایش دهید:
   ```bash
   PROXY_TIMEOUT=60
   ```
2. از سرویس بررسی IP سریع‌تر استفاده کنید
3. شبیه‌سازی تأخیر را غیرفعال کنید:
   ```bash
   SIMULATE_LATENCY=false
   ```

### مشکلات متریک

#### عدم دسترسی به متریک‌ها

```
Error: Unauthorized
```

**علل احتمالی:**

- احراز هویت فعال است
- مشخصات ورود نادرست
- پورت اشتباه

**راه‌حل‌ها:**

1. بررسی کنید احراز هویت فعال است یا خیر:
   ```bash
   METRICS_PROTECTED=false
   ```
2. مشخصات ورود را بررسی کنید:
   ```bash
   METRICS_USERNAME=user
   METRICS_PASSWORD=pass
   ```
3. پورت صحیح را بررسی کنید:
   ```bash
   METRICS_PORT=2112
   ```

#### خطاهای Pushgateway

```
Error pushing metrics: unexpected status code 401
```

**علل احتمالی:**

- URL pushgateway نامعتبر
- احراز هویت مورد نیاز است
- مشکلات شبکه

**راه‌حل‌ها:**

1. فرمت URL را بررسی کنید:
   ```bash
   METRICS_PUSH_URL="http://user:pass@host:9091"
   ```
2. اتصال شبکه را تأیید کنید
3. لاگ‌های pushgateway را بررسی کنید

### تداخل پورت

#### پورت در حال استفاده

```
error starting server: listen tcp :2112: bind: address already in use
```

**علل احتمالی:**

- سرویس دیگری از پورت استفاده می‌کند
- نمونه قبلی هنوز در حال اجرا است
- محدودیت‌های پورت از سمت سیستم

**راه‌حل‌ها:**

1. پورت متریک را تغییر دهید:
   ```bash
   METRICS_PORT=2113
   ```
2. پروسه‌های در حال اجرا را بررسی کنید:
   ```bash
   lsof -i :2112
   ```
3. سرویس‌های متداخل را متوقف کنید

#### مشکلات محدوده پورت SOCKS

```
error starting Xray: port already in use
```

**علل احتمالی:**

- تداخل محدوده پورت
- تعداد زیاد پروکسی
- محدودیت‌های پورت از سمت سیستم

**راه‌حل‌ها:**

1. پورت شروع را تغییر دهید:
   ```bash
   XRAY_START_PORT=20000
   ```
2. محدودیت‌های سیستم را بررسی کنید:
   ```bash
   ulimit -n
   ```
3. محدوده پورت را آزاد کنید

## تکنیک‌های اشکال‌زدایی

### فعال‌سازی لاگ اشکال‌زدایی

دو تنظیم سطح لاگ وجود دارد:

**سطح لاگ برنامه** - لاگ خود Xray Checker را کنترل می‌کند:

```bash
LOG_LEVEL=debug
```

سطوح موجود: `debug`، `info`، `warn`، `error`

لاگ اشکال‌زدایی موارد زیر را نشان می‌دهد:
- جزئیات تجزیه اشتراک
- نتایج بررسی پروکسی
- بارگذاری پیکربندی
- اطلاعات درخواست/پاسخ HTTP

**سطح لاگ Xray Core** - موتور پروکسی Xray تعبیه شده را کنترل می‌کند:

```bash
XRAY_LOG_LEVEL=debug
```

سطوح موجود: `debug`، `info`، `warning`، `error`، `none`

لاگ اشکال‌زدایی Xray نشان می‌دهد:
- تلاش‌های اتصال
- مذاکره پروتکل
- جزئیات لایه انتقال
- اطلاعات زمان‌بندی

:::tip[راه‌اندازی توصیه شده اشکال‌زدایی]
برای عیب‌یابی، هر دو را فعال کنید:
```bash
LOG_LEVEL=debug
XRAY_LOG_LEVEL=warning
```
این لاگ‌های دقیق برنامه را در حالی که خروجی Xray را کم حجم و قابل مدیریت نگه می‌دارد ارائه می‌دهد.
:::

### بررسی وضعیت فرآیند

```bash
# بررسی اجرای فرآیند
ps aux | grep xray-checker

# بررسی پورت‌های باز
netstat -tulpn | grep xray-checker
```

### تأیید اتصال شبکه

```bash
# تست سرویس بررسی IP
curl -v https://api.ipify.org?format=text

# تست اتصال پروکسی
curl --socks5 localhost:10000 -v https://api.ipify.org?format=text
```

### اشکال‌زدایی Docker

```bash
# بررسی لاگ‌های کانتینر
docker logs xray-checker

# دسترسی به shell کانتینر
docker exec -it xray-checker sh

# بررسی شبکه کانتینر
docker inspect xray-checker
```

## دریافت کمک

اگر هنوز مشکل دارید:

1. GitHub Issues را برای مشکلات مشابه بررسی کنید
2. issue جدید ایجاد کنید با:

   - پیام خطای کامل
   - پیکربندی استفاده شده
   - لاگ‌های اشکال‌زدایی
   - مراحل بازتولید

3. جزئیات محیط را ذکر کنید:
   - نسخه سیستم عامل
   - نسخه Docker (اگر استفاده می‌کنید)
   - نسخه Xray Checker
