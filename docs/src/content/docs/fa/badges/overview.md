---
title: سیستم نشان
description: جاسازی نشانگرهای وضعیت پروکسی در هر جایی
---

Xray Checker شامل یک سیستم نشان قدرتمند است که به شما امکان می‌دهد نشانگرهای وضعیت پروکسی در لحظه را در هر جایی جاسازی کنید — صفحات وضعیت، داشبوردها، مستندات یا فایل‌های README.

## ویژگی‌ها

- **وضعیت در لحظه** — نشان‌ها به طور خودکار با وضعیت پروکسی به‌روزرسانی می‌شوند
- **چندین تم** — تم‌های روشن و تاریک برای تطابق با طراحی شما
- **سبک‌های قابل تنظیم** — انواع مختلف، اندازه‌ها و گوشه‌های گرد
- **نمایش انعطاف‌پذیر** — نمایش/پنهان کردن نام و تأخیر به طور مستقل
- **جاسازی آسان** — کار از طریق iframe یا URL مستقیم

## شروع سریع

ساده‌ترین URL نشان:

```
https://your-xray-checker.com/?stableId=abc123def456
```

این یک نشان با وضعیت پروکسی، نام و تأخیر نمایش می‌دهد.

## دریافت Stable ID

هر پروکسی یک `stableId` منحصر به فرد دارد که بین راه‌اندازی‌های مجدد حفظ می‌شود. آن را از API دریافت کنید:

```bash
curl https://your-xray-checker.com/api/v1/public/proxies
```

پاسخ:

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

## پارامترهای موجود

| پارامتر | مقادیر | پیش‌فرض | توضیحات |
|---------|--------|---------|---------|
| `stableId` | `{id}` | اجباری | شناسه پایدار پروکسی از API |
| `theme` | `dark`, `light` | `dark` | تم رنگی |
| `variant` | `default`, `flat`, `pill`, `dot` | `default` | سبک نشان |
| `size` | `sm`, `md`, `lg` | `md` | اندازه نشان |
| `rounded` | `none`, `sm`, `md`, `lg`, `full` | `md` | گردی گوشه‌ها |
| `showName` | `true`, `false` | `true` | نمایش نام پروکسی |
| `showLatency` | `true`, `false` | `true` | نمایش مقدار تأخیر |
| `width` | عدد | auto | عرض سفارشی به پیکسل |
| `height` | عدد | auto | ارتفاع سفارشی به پیکسل |

## انواع نشان

### Default (پیش‌فرض)
نشان استاندارد با پس‌زمینه و حاشیه.
```
?stableId=abc123
```

### Flat (مسطح)
نشان حداقلی بدون پس‌زمینه یا حاشیه.
```
?stableId=abc123&variant=flat
```

### Pill (قرص)
نشان قرص‌شکل با گوشه‌های گرد.
```
?stableId=abc123&variant=pill
```

### Dot Only (فقط نقطه)
فقط نشانگر وضعیت نقطه‌ای.
```
?stableId=abc123&variant=dot
```

## مثال‌های جاسازی

### HTML iframe

```html
<iframe
  src="https://your-server.com/?stableId=abc123&theme=light"
  width="200"
  height="50"
  frameborder="0">
</iframe>
```

### چند نشان

یک داشبورد وضعیت با ترکیب چند نشان ایجاد کنید:

```html
<div style="display: flex; gap: 10px;">
  <iframe src="https://your-server.com/?stableId=server1" width="200" height="50" frameborder="0"></iframe>
  <iframe src="https://your-server.com/?stableId=server2" width="200" height="50" frameborder="0"></iframe>
  <iframe src="https://your-server.com/?stableId=server3" width="200" height="50" frameborder="0"></iframe>
</div>
```

## مراحل بعدی

- از [سازنده نشان](/fa/badges/playground) برای ساخت تعاملی URL نشان استفاده کنید
- درباره [صفحه وضعیت عمومی](/fa/configuration/status-page) برای اشتراک‌گذاری وضعیت پروکسی بیاموزید
