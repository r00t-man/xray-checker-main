---
title: Расширенная конфигурация
description: Расширенные параметры конфигурации
---

### Альтернативные сервисы проверки IP

Вы можете использовать альтернативные сервисы проверки IP (подробнее в разделе [методы проверки](/ru/configuration/check-methods)):

- `http://ip.sb`
- `https://api64.ipify.org`
- `http://ifconfig.me`

Пример:

```bash
PROXY_IP_CHECK_URL=http://ip.sb
```

### Альтернативные URL проверки статуса

Альтернативные URL для проверки статуса (подробнее в разделе [методы проверки](/ru/configuration/check-methods)):

- `http://www.gstatic.com/generate_204`
- `http://www.qualcomm.cn/generate_204`
- `http://cp.cloudflare.com/generate_204`

Пример:

```bash
PROXY_STATUS_CHECK_URL=http://www.gstatic.com/generate_204
```

### Настройка безопасности

Включение аутентификации для чувствительных эндпоинтов:

```bash
METRICS_PROTECTED=true
METRICS_USERNAME=custom_user
METRICS_PASSWORD=secure_password
```

### Маркировка экземпляров

Добавление меток экземпляров для распределенных установок:

```bash
METRICS_INSTANCE=datacenter-1
```

### Интервалы обновления

Настройка интервалов проверки и обновления:

```bash
# Проверка каждую минуту
PROXY_CHECK_INTERVAL=60

# Обновление подписки каждый час
SUBSCRIPTION_UPDATE_INTERVAL=3600
```

### Настройка логирования

Настройка логирования Xray Core:

```bash
# Включение отладочного логирования
XRAY_LOG_LEVEL=debug

# Отключение логирования
XRAY_LOG_LEVEL=none
```

### Настройка портов

Настройка диапазонов портов:

```bash
# Начало портов SOCKS5 с 20000
XRAY_START_PORT=20000

# Изменение порта метрик
METRICS_PORT=9090
```

### Настройка на своём собственном домене

У вас есть собственный домен `your-domain.com` и сайт на нём
и вы хотите отображать мониторинг по адресу `your-domain.com/xray/monitor`.

Запустите xray checher на том же сервере, где запущен ваш сайт
(параметр `-p 127.0.0.1:2112:2112` означает, что прямой доступ
к нему будет только с самого сервера):

:::caution
Для публично доступной страницы с мониторингом настоятельно рекомендуется
включить авторизацию по логину и паролю. Включить её можно с помощью переменных окружения
`METRICS_PROTECTED`, `METRICS_USERNAME`, `METRICS_PASSWORD`.
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

Откройте файл с настройками nginx (`sudo nano /etc/nginx/your-domain.com`),
найдите там главную секцию, она выглядит так:

```nginx
server {
    root /var/www/your-domain.com/html;

    index index.html;
    server_name your-domain.com;
    ...
}
```

Добавьте в неё 2 новых location для переадресации запросов на запущенный xray-checker:

```nginx

    # Обработка адреса /xray/monitor (без слеша в конце)
    location = /xray/monitor {
        return 301 https://$host$request_uri/;
    }

    # Обработка адреса  /xray/monitor/ - редирект на xray-checker
    location /xray/monitor/ {
        proxy_pass http://127.0.0.1:2112/xray/monitor/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
```

Проверьте настройки nginx и перезапустите его:

```bash
sudo nginx -t
sudo systemctl reload nginx
```

Проверьте, что мониторинг работает:

```bash
 curl -I -L https://your-domain.com/xray/monitor
```
