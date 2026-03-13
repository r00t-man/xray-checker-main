---
title: سفارشی‌سازی وب
description: سفارشی‌سازی رابط وب با قالب‌ها، استایل‌ها و فایل‌های خودتان
---

Xray Checker امکان سفارشی‌سازی کامل رابط وب را فراهم می‌کند. می‌توانید قالب پیش‌فرض را جایگزین کنید، استایل‌های سفارشی اضافه کنید، لوگو و favicon را تغییر دهید و هر فایل استاتیک دیگری اضافه کنید.

## فعال‌سازی فایل‌های سفارشی

مسیر پوشه فایل‌های سفارشی را مشخص کنید:

```bash
# متغیر محیطی
WEB_CUSTOM_ASSETS_PATH=/path/to/custom

# فلگ CLI
xray-checker --web-custom-assets-path=/path/to/custom
```

اگر مسیر تنظیم شده و پوشه وجود داشته باشد، فایل‌های سفارشی در زمان راه‌اندازی بارگذاری می‌شوند.

## ساختار پوشه

فایل‌های خود را در یک پوشه مسطح قرار دهید (بدون زیرپوشه):

```
custom/
  ├── index.html       # جایگزینی کامل قالب (اختیاری)
  ├── logo.svg         # لوگوی واحد برای هر دو تم (اختیاری)
  ├── logo.png         # لوگوی واحد PNG (اختیاری)
  ├── logo-dark.svg    # لوگو برای تم تاریک (اختیاری)
  ├── logo-dark.png    # لوگو برای تم تاریک PNG (اختیاری)
  ├── logo-light.svg   # لوگو برای تم روشن (اختیاری)
  ├── logo-light.png   # لوگو برای تم روشن PNG (اختیاری)
  ├── favicon.ico      # جایگزینی favicon (اختیاری)
  ├── custom.css       # استایل‌های اضافی، به‌صورت خودکار تزریق می‌شوند (اختیاری)
  └── any-file.ext     # در /static/any-file.ext قابل دسترسی است
```

### فایل‌های لوگو

شما می‌توانید لوگو را به دو روش سفارشی کنید:

1. **لوگوی واحد** — فایل `logo.svg` یا `logo.png` را برای استفاده از یک لوگو در هر دو تم ارائه دهید
2. **لوگوهای خاص تم** — فایل‌های `logo-dark.svg`/`logo-dark.png` و `logo-light.svg`/`logo-light.png` را برای لوگوهای متفاوت ارائه دهید

ترتیب اولویت (اولین یافت شده استفاده می‌شود):
1. `logo-dark.svg` / `logo-light.svg` (SVG خاص تم)
2. `logo-dark.png` / `logo-light.png` (PNG خاص تم)
3. `logo.svg` (SVG جهانی)
4. `logo.png` (PNG جهانی)

## استایل‌های سفارشی (custom.css)

ساده‌ترین راه برای تغییر ظاهر. اگر `custom.css` وجود داشته باشد، به‌صورت خودکار بعد از استایل‌های پیش‌فرض بارگذاری می‌شود.

### متغیرهای CSS

رنگ‌های تم را با استفاده از متغیرهای CSS بازنویسی کنید:

```css
:root {
  /* رنگ‌های پس‌زمینه */
  --bg-primary: #0a0a0f;
  --bg-secondary: #12121a;
  --bg-tertiary: #1a1a24;

  /* رنگ‌های متن */
  --text-primary: #f4f4f5;
  --text-secondary: #a1a1aa;
  --text-muted: #71717a;

  /* رنگ‌های تأکیدی */
  --color-green: #22c55e;
  --color-red: #ef4444;
  --color-orange: #f97316;

  /* حاشیه */
  --border: #27272a;
}
```

### مثال: تم روشن

```css
:root {
  --bg-primary: #ffffff;
  --bg-secondary: #f8f9fa;
  --bg-tertiary: #e9ecef;
  --text-primary: #212529;
  --text-secondary: #495057;
  --text-muted: #6c757d;
  --border: #dee2e6;
}
```

### مثال: اندازه لوگو

```css
.header-logo img {
  width: 48px;
  height: 48px;
}
```

## جایگزینی کامل قالب (index.html)

برای کنترل کامل، فایل `index.html` خودتان را ارائه دهید. این یک قالب Go با دسترسی به تمام داده‌های صفحه است.

:::caution
قالب‌های سفارشی ممکن است پس از به‌روزرسانی‌ها در صورت تغییر ساختار داده خراب شوند. با مسئولیت خودتان استفاده کنید.
:::

### سینتکس قالب Go

```html
{{ .Variable }}           <!-- خروجی متغیر -->
{{ if .Condition }}...{{ end }}
{{ range .Array }}...{{ end }}
{{ formatLatency .Latency }}  <!-- فرمت به "123ms" یا "n/a" -->
```

### متغیرهای موجود

#### PageData (شیء ریشه)

| متغیر | نوع | توضیحات |
|-------|-----|---------|
| `.Version` | string | نسخه Xray Checker |
| `.Host` | string | هاست سرور |
| `.Port` | string | پورت سرور |
| `.CheckInterval` | int | فاصله بررسی پروکسی به ثانیه |
| `.Timeout` | int | مهلت بررسی پروکسی به ثانیه |
| `.CheckMethod` | string | روش بررسی: `ip`، `status` یا `download` |
| `.IPCheckUrl` | string | URL برای بررسی IP |
| `.StatusCheckUrl` | string | URL برای بررسی وضعیت |
| `.DownloadUrl` | string | URL برای بررسی دانلود |
| `.SimulateLatency` | bool | آیا شبیه‌سازی تأخیر فعال است |
| `.SubscriptionUpdate` | bool | آیا به‌روزرسانی خودکار اشتراک فعال است |
| `.SubscriptionUpdateInterval` | int | فاصله به‌روزرسانی اشتراک به ثانیه |
| `.StartPort` | int | شماره اولین پورت پروکسی |
| `.Instance` | string | برچسب instance برای متریک‌ها |
| `.PushUrl` | string | URL پوش‌گیت‌وی Prometheus |
| `.ShowServerDetails` | bool | آیا IP و پورت سرورها نمایش داده شود |
| `.IsPublic` | bool | آیا حالت عمومی فعال است |
| `.SubscriptionName` | string | نام اشتراک برای نمایش |
| `.Endpoints` | []EndpointInfo | آرایه پروکسی‌ها |

#### EndpointInfo (هر آیتم در `.Endpoints`)

| متغیر | نوع | دسترسی | توضیحات |
|-------|-----|--------|---------|
| `.Name` | string | همیشه | نام پروکسی |
| `.StableID` | string | همیشه | شناسه یکتای پروکسی |
| `.Index` | int | همیشه | ایندکس پروکسی (از 0) |
| `.Status` | bool | همیشه | `true` اگر آنلاین |
| `.Latency` | time.Duration | همیشه | تأخیر پاسخ |
| `.ServerInfo` | string | وقتی `ShowServerDetails && !IsPublic` | آدرس و پورت سرور |
| `.ProxyPort` | int | وقتی `ShowServerDetails && !IsPublic` | پورت محلی پروکسی |
| `.URL` | string | وقتی `!IsPublic` | URL اندپوینت وضعیت |

### توابع قالب

| تابع | توضیحات | مثال |
|------|---------|------|
| `formatLatency` | فرمت مدت به میلی‌ثانیه | `{{ formatLatency .Latency }}` → `"123ms"` یا `"n/a"` |

### رندر شرطی

```html
<!-- فقط در حالت غیرعمومی نمایش بده -->
{{ if not .IsPublic }}
  <a href="{{ .URL }}">Config</a>
{{ end }}

<!-- نمایش جزئیات سرور در صورت فعال بودن -->
{{ if .ShowServerDetails }}
  <span>{{ .ServerInfo }}</span>
{{ end }}

<!-- حلقه روی پروکسی‌ها -->
{{ range .Endpoints }}
  <div class="{{ if .Status }}online{{ else }}offline{{ end }}">
    {{ .Name }} - {{ formatLatency .Latency }}
  </div>
{{ end }}
```

### مثال قالب حداقلی

```html
<!DOCTYPE html>
<html>
<head>
  <title>{{ if .SubscriptionName }}{{ .SubscriptionName }}{{ else }}وضعیت{{ end }}</title>
</head>
<body>
  <h1>وضعیت پروکسی</h1>
  {{ range .Endpoints }}
  <div>
    {{ if .Status }}✓{{ else }}✗{{ end }}
    {{ .Name }}
    ({{ formatLatency .Latency }})
  </div>
  {{ end }}
</body>
</html>
```

## مثال Docker

```yaml
services:
  xray-checker:
    image: kutovoys/xray-checker:latest
    volumes:
      - ./custom:/app/custom:ro
    environment:
      - WEB_CUSTOM_ASSETS_PATH=/app/custom
      - SUBSCRIPTION_URL=https://...
```

ساختار پوشه‌ها:
```
my-project/
  ├── docker-compose.yml
  └── custom/
      ├── logo.svg        # یا logo-dark.svg + logo-light.svg
      ├── favicon.ico
      └── custom.css
```

## لاگ‌های راه‌اندازی

هنگام بارگذاری فایل‌های سفارشی، خواهید دید:

```
INFO  Custom assets enabled: /app/custom
INFO  Custom assets loaded:
INFO    ✓ logo.svg
INFO    ✓ custom.css
INFO  Using default template
```

یا با قالب سفارشی:
```
INFO  Custom assets loaded:
INFO    ✓ index.html
INFO    ✓ custom.css
INFO  Using custom template: index.html
```

## خطاها

| خطا | دلیل |
|-----|------|
| `custom assets directory does not exist` | مسیر تنظیم شده اما پوشه پیدا نشد |
| `failed to parse custom template` | سینتکس نامعتبر قالب Go در index.html |

برنامه **راه‌اندازی نمی‌شود** اگر مسیر فایل‌های سفارشی تنظیم شده اما نامعتبر باشد.
