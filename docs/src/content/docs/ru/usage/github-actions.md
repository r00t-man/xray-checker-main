---
title: GitHub Actions
description: Запуск Xray Checker с помощью GitHub Actions
---

# Интеграция с GitHub Actions

Вы можете запускать Xray Checker с помощью GitHub Actions. Этот подход полезен, когда вам нужно выполнять проверки из разных локаций или у вас нет собственного сервера.

## Быстрая настройка

1. Сделайте форк репозитория [xray-checker-in-actions](https://github.com/kutovoys/xray-checker-in-actions)
2. Настройте следующие секреты в вашем форке репозитория:
   - `SUBSCRIPTION_URL`: URL вашей подписки
   - `PUSH_URL`: URL Prometheus pushgateway для сбора метрик
   - `INSTANCE`: (Опционально) Имя экземпляра для идентификации метрик

Action будет:

- Запускаться каждые 5 минут
- Использовать последнюю версию Xray Checker
- Отправлять метрики в ваш Prometheus pushgateway
- Запускаться с флагом `--run-once` для обеспечения чистого выполнения

Этот метод требует правильно настроенного Prometheus pushgateway, так как не может напрямую экспортировать метрики. Метрики будут отправляться на указанный вами `PUSH_URL` с меткой экземпляра из вашей конфигурации.

## Расширенные конфигурации

Если вам нужен больший контроль над настройками GitHub Actions, вот несколько расширенных конфигураций.

### Настройка нескольких регионов

Запуск проверок из разных регионов одновременно:

```yaml
name: Xray Checker
on:
  schedule:
    - cron: "*/5 * * * *"
  workflow_dispatch:

jobs:
  check:
    strategy:
      matrix:
        include:
          - location: us
            runs-on: us-east-1
          - location: eu
            runs-on: eu-west-1
          - location: asia
            runs-on: ap-east-1

    runs-on: ${{ matrix.runs-on }}

    steps:
      - name: Run Xray Checker
        uses: docker://kutovoys/xray-checker
        env:
          SUBSCRIPTION_URL: ${{ secrets.SUBSCRIPTION_URL }}
          METRICS_PUSH_URL: ${{ secrets.PUSH_URL }}
          METRICS_INSTANCE: ${{ matrix.location }}
          RUN_ONCE: true
```

### Уведомления об ошибках

Добавление уведомлений в Slack или по электронной почте при неудачных проверках:

```yaml
steps:
  - name: Run Xray Checker
    id: checker
    uses: docker://kutovoys/xray-checker
    env:
      SUBSCRIPTION_URL: ${{ secrets.SUBSCRIPTION_URL }}
      METRICS_PUSH_URL: ${{ secrets.PUSH_URL }}
    continue-on-error: true

  - name: Notify on Failure
    if: steps.checker.outcome == 'failure'
    uses: actions/github-script@v6
    with:
      script: |
        github.rest.issues.create({
          owner: context.repo.owner,
          repo: context.repo.repo,
          title: 'Xray Checker Failed',
          body: 'Проверка не удалась в запуске workflow: ' + context.runId
        })
```

### Пользовательские интервалы проверки

Различные шаблоны расписания в зависимости от ваших потребностей:

```yaml
on:
  schedule:
    - cron: "*/5 * * * *" # Каждые 5 минут (по умолчанию)
    - cron: "0 * * * *" # Каждый час
    - cron: "0 */2 * * *" # Каждые 2 часа
```

### Оптимизация ресурсов

Оптимизация использования GitHub Actions с помощью управления параллельным выполнением:

```yaml
jobs:
  check:
    runs-on: ubuntu-latest
    timeout-minutes: 5

    concurrency:
      group: ${{ github.workflow }}-${{ github.ref }}
      cancel-in-progress: true

    steps:
      - name: Run Xray Checker
        uses: docker://kutovoys/xray-checker
        env:
          SUBSCRIPTION_URL: ${{ secrets.SUBSCRIPTION_URL }}
          METRICS_PUSH_URL: ${{ secrets.PUSH_URL }}
          RUN_ONCE: true
```

## Настройка мониторинга

### Необходимая конфигурация Prometheus

Чтобы собирать метрики из pushgateway, добавьте это в вашу конфигурацию Prometheus:

```yaml
scrape_configs:
  - job_name: "pushgateway"
    honor_labels: true
    static_configs:
      - targets: ["pushgateway:9091"]
```

Метрики будут появляться с меткой экземпляра, которую вы указали в конфигурации GitHub Actions, что позволяет отслеживать проверки из разных локаций.
