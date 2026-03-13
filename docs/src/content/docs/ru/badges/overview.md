---
title: Система бейджей
description: Встраивайте индикаторы статуса прокси куда угодно
---

Xray Checker включает мощную систему бейджей, которая позволяет встраивать индикаторы статуса прокси в реальном времени куда угодно — на страницы статуса, дашборды, документацию или README файлы.

## Возможности

- **Статус в реальном времени** — Бейджи автоматически обновляются с состоянием прокси
- **Несколько тем** — Тёмная и светлая темы для соответствия вашему дизайну
- **Настраиваемые стили** — Различные варианты, размеры и скругления углов
- **Гибкое отображение** — Показ/скрытие имени и задержки независимо друг от друга
- **Простое встраивание** — Работает через iframe или прямой URL

## Быстрый старт

Простейший URL бейджа:

```
https://your-xray-checker.com/?stableId=abc123def456
```

Отображает бейдж со статусом прокси, именем и задержкой.

## Получение Stable ID

Каждый прокси имеет уникальный `stableId`, который сохраняется между перезапусками. Получите его через API:

```bash
curl https://your-xray-checker.com/api/v1/public/proxies
```

Ответ:

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

## Доступные параметры

| Параметр | Значения | По умолчанию | Описание |
|----------|----------|--------------|----------|
| `stableId` | `{id}` | обязательно | Stable ID прокси из API |
| `theme` | `dark`, `light` | `dark` | Цветовая тема |
| `variant` | `default`, `flat`, `pill`, `dot` | `default` | Стиль бейджа |
| `size` | `sm`, `md`, `lg` | `md` | Размер бейджа |
| `rounded` | `none`, `sm`, `md`, `lg`, `full` | `md` | Скругление углов |
| `showName` | `true`, `false` | `true` | Показывать имя прокси |
| `showLatency` | `true`, `false` | `true` | Показывать задержку |
| `width` | число | auto | Пользовательская ширина в пикселях |
| `height` | число | auto | Пользовательская высота в пикселях |

## Варианты бейджей

### Default (По умолчанию)
Стандартный бейдж с фоном и рамкой.
```
?stableId=abc123
```

### Flat (Плоский)
Минималистичный бейдж без фона и рамки.
```
?stableId=abc123&variant=flat
```

### Pill (Таблетка)
Бейдж в форме таблетки со скруглёнными краями.
```
?stableId=abc123&variant=pill
```

### Dot Only (Только точка)
Только индикатор статуса в виде точки.
```
?stableId=abc123&variant=dot
```

## Примеры встраивания

### HTML iframe

```html
<iframe
  src="https://your-server.com/?stableId=abc123&theme=light"
  width="200"
  height="50"
  frameborder="0">
</iframe>
```

### Несколько бейджей

Создайте дашборд статуса, комбинируя несколько бейджей:

```html
<div style="display: flex; gap: 10px;">
  <iframe src="https://your-server.com/?stableId=server1" width="200" height="50" frameborder="0"></iframe>
  <iframe src="https://your-server.com/?stableId=server2" width="200" height="50" frameborder="0"></iframe>
  <iframe src="https://your-server.com/?stableId=server3" width="200" height="50" frameborder="0"></iframe>
</div>
```

## Следующие шаги

- Используйте [Конструктор бейджей](/ru/badges/playground) для интерактивного создания URL бейджа
- Узнайте о [Публичной странице статуса](/ru/configuration/status-page) для публикации статуса прокси
