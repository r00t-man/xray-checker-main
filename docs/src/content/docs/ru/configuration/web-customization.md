---
title: Кастомизация веб-интерфейса
description: Настройка веб-интерфейса с собственными шаблонами, стилями и ассетами
---

Xray Checker позволяет полностью кастомизировать веб-интерфейс. Вы можете заменить стандартный шаблон, добавить свои стили, изменить логотип, favicon и добавить любые другие статические файлы.

## Включение кастомных ассетов

Укажите путь к директории с кастомными ассетами:

```bash
# Переменная окружения
WEB_CUSTOM_ASSETS_PATH=/path/to/custom

# Флаг CLI
xray-checker --web-custom-assets-path=/path/to/custom
```

Если путь указан и директория существует, кастомные ассеты будут загружены при запуске.

## Структура директории

Разместите файлы в плоской директории (без вложенных папок):

```
custom/
  ├── index.html       # Полная замена шаблона (опционально)
  ├── logo.svg         # Единый логотип для обеих тем (опционально)
  ├── logo.png         # Единый логотип PNG (опционально)
  ├── logo-dark.svg    # Логотип для тёмной темы (опционально)
  ├── logo-dark.png    # Логотип для тёмной темы PNG (опционально)
  ├── logo-light.svg   # Логотип для светлой темы (опционально)
  ├── logo-light.png   # Логотип для светлой темы PNG (опционально)
  ├── favicon.ico      # Замена favicon (опционально)
  ├── custom.css       # Дополнительные стили, инжектятся автоматически (опционально)
  └── any-file.ext     # Доступен по /static/any-file.ext
```

### Файлы логотипа

Вы можете кастомизировать логотип двумя способами:

1. **Единый логотип** — предоставьте `logo.svg` или `logo.png` для использования одного логотипа в обеих темах
2. **Логотипы для каждой темы** — предоставьте `logo-dark.svg`/`logo-dark.png` и `logo-light.svg`/`logo-light.png` для разных логотипов

Порядок приоритета (используется первый найденный):
1. `logo-dark.svg` / `logo-light.svg` (SVG для конкретной темы)
2. `logo-dark.png` / `logo-light.png` (PNG для конкретной темы)
3. `logo.svg` (универсальный SVG)
4. `logo.png` (универсальный PNG)

## Кастомные стили (custom.css)

Самый простой способ изменить внешний вид. Если `custom.css` существует, он автоматически подключается после стандартных стилей.

### CSS переменные

Переопределите цвета темы с помощью CSS переменных:

```css
:root {
  /* Цвета фона */
  --bg-primary: #0a0a0f;
  --bg-secondary: #12121a;
  --bg-tertiary: #1a1a24;

  /* Цвета текста */
  --text-primary: #f4f4f5;
  --text-secondary: #a1a1aa;
  --text-muted: #71717a;

  /* Акцентные цвета */
  --color-green: #22c55e;
  --color-red: #ef4444;
  --color-orange: #f97316;

  /* Границы */
  --border: #27272a;
}
```

### Пример: Светлая тема

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

### Пример: Размер логотипа

```css
.header-logo img {
  width: 48px;
  height: 48px;
}
```

## Полная замена шаблона (index.html)

Для полного контроля предоставьте свой `index.html`. Это Go-шаблон с доступом ко всем данным страницы.

:::caution
Кастомные шаблоны могут сломаться после обновлений, если изменится структура данных. Используйте на свой риск.
:::

### Синтаксис Go-шаблонов

```html
{{ .Variable }}           <!-- Вывод переменной -->
{{ if .Condition }}...{{ end }}
{{ range .Array }}...{{ end }}
{{ formatLatency .Latency }}  <!-- Форматирование как "123ms" или "n/a" -->
```

### Доступные переменные

#### PageData (корневой объект)

| Переменная | Тип | Описание |
|------------|-----|----------|
| `.Version` | string | Версия Xray Checker |
| `.Host` | string | Хост сервера |
| `.Port` | string | Порт сервера |
| `.CheckInterval` | int | Интервал проверки прокси в секундах |
| `.Timeout` | int | Таймаут проверки прокси в секундах |
| `.CheckMethod` | string | Метод проверки: `ip`, `status` или `download` |
| `.IPCheckUrl` | string | URL для проверки IP |
| `.StatusCheckUrl` | string | URL для проверки статуса |
| `.DownloadUrl` | string | URL для проверки загрузки |
| `.SimulateLatency` | bool | Включена ли симуляция задержки |
| `.SubscriptionUpdate` | bool | Включено ли автообновление подписки |
| `.SubscriptionUpdateInterval` | int | Интервал обновления подписки в секундах |
| `.StartPort` | int | Начальный порт прокси |
| `.Instance` | string | Метка instance для метрик |
| `.PushUrl` | string | URL Prometheus pushgateway |
| `.ShowServerDetails` | bool | Показывать ли IP и порты серверов |
| `.IsPublic` | bool | Включён ли публичный режим |
| `.SubscriptionName` | string | Имя подписки для отображения |
| `.Endpoints` | []EndpointInfo | Массив прокси |

#### EndpointInfo (каждый элемент в `.Endpoints`)

| Переменная | Тип | Доступность | Описание |
|------------|-----|-------------|----------|
| `.Name` | string | Всегда | Имя прокси |
| `.StableID` | string | Всегда | Уникальный идентификатор прокси |
| `.Index` | int | Всегда | Индекс прокси (с 0) |
| `.Status` | bool | Всегда | `true` если онлайн |
| `.Latency` | time.Duration | Всегда | Задержка ответа |
| `.ServerInfo` | string | При `ShowServerDetails && !IsPublic` | Адрес и порт сервера |
| `.ProxyPort` | int | При `ShowServerDetails && !IsPublic` | Локальный порт прокси |
| `.URL` | string | При `!IsPublic` | URL эндпоинта статуса |

### Функции шаблона

| Функция | Описание | Пример |
|---------|----------|--------|
| `formatLatency` | Форматирует длительность в миллисекунды | `{{ formatLatency .Latency }}` → `"123ms"` или `"n/a"` |

### Условный рендеринг

```html
<!-- Показать только в непубличном режиме -->
{{ if not .IsPublic }}
  <a href="{{ .URL }}">Config</a>
{{ end }}

<!-- Показать детали сервера если включено -->
{{ if .ShowServerDetails }}
  <span>{{ .ServerInfo }}</span>
{{ end }}

<!-- Цикл по прокси -->
{{ range .Endpoints }}
  <div class="{{ if .Status }}online{{ else }}offline{{ end }}">
    {{ .Name }} - {{ formatLatency .Latency }}
  </div>
{{ end }}
```

### Минимальный пример шаблона

```html
<!DOCTYPE html>
<html>
<head>
  <title>{{ if .SubscriptionName }}{{ .SubscriptionName }}{{ else }}Статус{{ end }}</title>
</head>
<body>
  <h1>Статус прокси</h1>
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

## Пример для Docker

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

Структура директорий:
```
my-project/
  ├── docker-compose.yml
  └── custom/
      ├── logo.svg        # или logo-dark.svg + logo-light.svg
      ├── favicon.ico
      └── custom.css
```

## Логи при запуске

При загрузке кастомных ассетов вы увидите:

```
INFO  Custom assets enabled: /app/custom
INFO  Custom assets loaded:
INFO    ✓ logo.svg
INFO    ✓ custom.css
INFO  Using default template
```

Или с кастомным шаблоном:
```
INFO  Custom assets loaded:
INFO    ✓ index.html
INFO    ✓ custom.css
INFO  Using custom template: index.html
```

## Ошибки

| Ошибка | Причина |
|--------|---------|
| `custom assets directory does not exist` | Путь указан, но директория не найдена |
| `failed to parse custom template` | Неверный синтаксис Go-шаблона в index.html |

Приложение **не запустится**, если путь к кастомным ассетам указан, но невалиден.
