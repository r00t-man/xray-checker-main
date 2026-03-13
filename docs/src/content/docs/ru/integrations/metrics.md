---
title: Метрики
description: Параметры и примеры метрик
---

Xray Checker предоставляет две метрики Prometheus для мониторинга состояния и производительности прокси. Подробнее о настройке интеграции читайте в разделе [Prometheus](/ru/integrations/prometheus).

Для визуализации метрик рекомендуется использовать [Grafana](/ru/integrations/grafana).

### xray_proxy_status

Метрика состояния, показывающая доступность прокси:

- Тип: Gauge
- Значения: 1 (работает) или 0 (не работает)
- Метки:
  - `protocol`: Протокол прокси (vless/vmess/trojan/shadowsocks)
  - `address`: Адрес и порт сервера
  - `name`: Имя конфигурации прокси
  - `sub_name`: Имя подписки (из фрагмента URL или заголовка profile-title)
  - `instance`: Имя экземпляра (если настроено)

:::tip
Загляните в [расширенную конфигурацию](/ru/configuration/advanced-conf#маркировка-экземпляров) для настройки меток экземпляров.
:::

Пример:

```text
# HELP xray_proxy_status Статус прокси-соединения (1: успешно, 0: неудача)
# TYPE xray_proxy_status gauge
xray_proxy_status{protocol="vless",address="example.com:443",name="proxy1",sub_name="Premium VPN",instance="dc1"} 1
```

### xray_proxy_latency_ms

Метрика задержки, показывающая время отклика соединения:

- Тип: Gauge
- Значения: Миллисекунды (0 при неудаче)
- Метки: Те же, что и у xray_proxy_status

Пример:

```text
# HELP xray_proxy_latency_ms Задержка прокси-соединения в миллисекундах
# TYPE xray_proxy_latency_ms gauge
xray_proxy_latency_ms{protocol="vless",address="example.com:443",name="proxy1",sub_name="Premium VPN",instance="dc1"} 156
```
