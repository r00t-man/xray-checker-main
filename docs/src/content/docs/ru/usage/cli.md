---
title: CLI
description: Использование CLI в Xray Checker
---

### Базовое использование командной строки

CLI интерфейс предоставляет полный контроль над функциональностью Xray Checker через аргументы командной строки.

### Установка

Скачайте последнюю версию бинарного файла из релизов:

```bash
# Для Linux amd64
curl -sL -o - $(curl -s https://api.github.com/repos/kutovoys/xray-checker/releases/latest | grep "browser_download_url.*linux-amd64.tar.gz" | cut -d'"' -f4) | tar -xz
chmod +x xray-checker

# Для Linux arm64
curl -sL -o - $(curl -s https://api.github.com/repos/kutovoys/xray-checker/releases/latest | grep "browser_download_url.*linux-arm64.tar.gz" | cut -d'"' -f4) | tar -xz
chmod +x xray-checker

# Для macOS (Intel)
curl -sL -o - $(curl -s https://api.github.com/repos/kutovoys/xray-checker/releases/latest | grep "browser_download_url.*darwin-amd64.tar.gz" | cut -d'"' -f4) | tar -xz
chmod +x xray-checker

# Для macOS (Apple Silicon)
curl -sL -o - $(curl -s https://api.github.com/repos/kutovoys/xray-checker/releases/latest | grep "browser_download_url.*darwin-arm64.tar.gz" | cut -d'"' -f4) | tar -xz
chmod +x xray-checker
```

### Базовое использование

Минимально необходимая конфигурация:

```bash
./xray-checker --subscription-url=https://your-subscription-url/sub
```

### Несколько подписок

Вы можете указать несколько URL подписок, используя флаг `--subscription-url` несколько раз:

```bash
./xray-checker \
  --subscription-url=https://provider1.com/sub \
  --subscription-url=https://provider2.com/sub \
  --subscription-url=file:///path/to/local/config.json
```

Все прокси из всех подписок будут объединены и мониторятся вместе.

### Пример полной конфигурации

```bash
./xray-checker \
  --subscription-url=https://your-subscription-url/sub \
  --subscription-update=true \
  --subscription-update-interval=300 \
  --proxy-check-interval=300 \
  --proxy-timeout=30 \
  --proxy-check-method=ip \
  --proxy-ip-check-url="https://api.ipify.org?format=text" \
  --proxy-status-check-url="http://cp.cloudflare.com/generate_204" \
  --proxy-download-url="https://proof.ovh.net/files/1Mb.dat" \
  --proxy-download-timeout=60 \
  --proxy-download-min-size=51200 \
  --proxy-resolve-domains=false \
  --simulate-latency=true \
  --xray-start-port=10000 \
  --xray-log-level=none \
  --metrics-host=0.0.0.0 \
  --metrics-port=2112 \
  --metrics-protected=true \
  --metrics-username=custom_user \
  --metrics-password=custom_pass \
  --metrics-instance=node-1 \
  --metrics-push-url="https://push.example.com" \
  --metrics-base-path="/xray/monitor" \
  --web-show-details=false \
  --web-public=false \
  --log-level=info \
  --run-once=false
```

### Основные операции CLI

Проверка версии:

```bash
./xray-checker --version
```

Запуск одного цикла проверки:

```bash
./xray-checker --subscription-url=https://your-sub-url --run-once
```

Включение аутентификации метрик:

```bash
./xray-checker \
  --subscription-url=https://your-sub-url \
  --metrics-protected=true \
  --metrics-username=user \
  --metrics-password=pass
```

Изменение портов:

```bash
./xray-checker \
  --subscription-url=https://your-sub-url \
  --metrics-host=127.0.0.1 \
  --metrics-port=3000 \
  --xray-start-port=20000
```
