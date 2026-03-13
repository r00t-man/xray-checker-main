---
title: Быстрый старт
description: Быстрый старт с Xray Checker
---

Запустите Xray Checker за считанные минуты, выполнив эти простые шаги.

## Предварительные требования

- URL подписки для ваших прокси
- Docker (опционально, для развертывания в контейнере)
- Prometheus (опционально, для сбора метрик)

## Настройка за 5 минут

### Использование Docker (Рекомендуется)

1. Загрузите образ:

```bash
docker pull kutovoys/xray-checker
```

2. Запустите с базовой конфигурацией:

```bash
docker run -d \
  -e SUBSCRIPTION_URL=https://your-subscription-url/sub \
  -p 2112:2112 \
  kutovoys/xray-checker
```

3. Проверьте статус:

```bash
curl http://localhost:2112/health
```

### Использование исполняемого файла

1. Загрузите последний релиз:

```bash
# Linux amd64
curl -sL -o - $(curl -s https://api.github.com/repos/kutovoys/xray-checker/releases/latest | grep "browser_download_url.*linux-amd64.tar.gz" | cut -d'"' -f4) | tar -xz
chmod +x xray-checker

# Linux arm64
curl -sL -o - $(curl -s https://api.github.com/repos/kutovoys/xray-checker/releases/latest | grep "browser_download_url.*linux-arm64.tar.gz" | cut -d'"' -f4) | tar -xz
chmod +x xray-checker
```

2. Запустите с базовой конфигурацией:

```bash
./xray-checker --subscription-url=https://your-subscription-url/sub
```

## Проверка установки

1. Откройте веб-интерфейс:

   - Перейдите по адресу `http://localhost:2112`
   - Вы должны увидеть панель управления со статусом прокси

2. Проверьте метрики:

   - Перейдите по адресу `http://localhost:2112/metrics`
   - Вы должны увидеть метрики Prometheus

3. Проверьте статус прокси:
   - Нажмите на любую ссылку прокси в веб-интерфейсе
   - Проверьте ответ конечной точки статуса

## Следующие шаги

1. Настройте Prometheus:

```yaml
scrape_configs:
  - job_name: "xray-checker"
    static_configs:
      - targets: ["localhost:2112"]
```

2. Настройте Uptime Kuma:

   - Добавьте новый монитор
   - Используйте специфичные для прокси конечные точки
   - Настройте оповещения

3. Настройте конфигурацию:
   - Настройте интервалы проверки
   - Настройте аутентификацию
   - Настройте отправку метрик

## Основные команды

Проверка версии:

```bash
./xray-checker --version
```

Запуск одиночной проверки:

```bash
./xray-checker --subscription-url=https://your-sub-url --run-once
```

Включение аутентификации:

```bash
./xray-checker --subscription-url=https://your-sub-url \
  --metrics-protected=true \
  --metrics-username=user \
  --metrics-password=pass
```

## Устранение неполадок

1. Проверка статуса сервиса:

```bash
curl http://localhost:2112/health
```

2. Просмотр логов:

```bash
docker logs xray-checker
```

3. Проверка метрик:

```bash
curl http://localhost:2112/metrics
```

## Нужна помощь?

- Ознакомьтесь с полной документацией
- Создайте issue на GitHub
- Присоединяйтесь к обсуждениям сообщества
