---
title: GitHub Actions
description: اجرای Xray Checker با GitHub Actions
---

# یکپارچه‌سازی با GitHub Actions

می‌توانید Xray Checker را با استفاده از GitHub Actions اجرا کنید. این روش زمانی مفید است که نیاز دارید بررسی‌ها را از مکان‌های مختلف اجرا کنید یا سرور اختصاصی ندارید.

## راه‌اندازی سریع

1. مخزن [xray-checker-in-actions](https://github.com/kutovoys/xray-checker-in-actions) را fork کنید
2. secrets زیر را در مخزن fork شده خود پیکربندی کنید:
   - `SUBSCRIPTION_URL`: آدرس اشتراک شما
   - `PUSH_URL`: آدرس Prometheus pushgateway برای جمع‌آوری متریک‌ها
   - `INSTANCE`: (اختیاری) نام نمونه برای شناسایی متریک‌ها

این Action:

- هر ۵ دقیقه اجرا می‌شود
- از آخرین نسخه Xray Checker استفاده می‌کند
- متریک‌ها را به Prometheus pushgateway شما ارسال می‌کند
- با فلگ `--run-once` اجرا می‌شود تا اجرای تمیز تضمین شود

این روش به یک Prometheus pushgateway نیاز دارد زیرا نمی‌تواند متریک‌ها را مستقیماً ارائه دهد. متریک‌ها به `PUSH_URL` مشخص شده با برچسب instance از پیکربندی شما ارسال می‌شوند.

## پیکربندی‌های پیشرفته

اگر به کنترل بیشتر بر راه‌اندازی GitHub Actions نیاز دارید، در اینجا چند پیکربندی پیشرفته آورده شده است.

### راه‌اندازی چند منطقه‌ای

اجرای بررسی‌ها از مناطق مختلف به طور همزمان:

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

### اعلان‌های خطا

اضافه کردن اعلان‌های Slack یا Email برای بررسی‌های ناموفق:

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
          body: 'Check failed in workflow run: ' + context.runId
        })
```

### فواصل بررسی سفارشی

الگوهای زمان‌بندی مختلف بر اساس نیازهای شما:

```yaml
on:
  schedule:
    - cron: "*/5 * * * *" # هر ۵ دقیقه (پیش‌فرض)
    - cron: "0 * * * *" # هر ساعت
    - cron: "0 */2 * * *" # هر ۲ ساعت
```

### بهینه‌سازی منابع

بهینه‌سازی استفاده از GitHub Actions با کنترل همزمانی:

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

## راه‌اندازی نظارت

### پیکربندی مورد نیاز Prometheus

برای جمع‌آوری متریک‌ها از pushgateway، این را به پیکربندی Prometheus خود اضافه کنید:

```yaml
scrape_configs:
  - job_name: "pushgateway"
    honor_labels: true
    static_configs:
      - targets: ["pushgateway:9091"]
```

متریک‌ها با برچسب instance که در پیکربندی GitHub Actions مشخص کرده‌اید ظاهر می‌شوند، که به شما امکان می‌دهد بررسی‌ها را از مکان‌های مختلف پیگیری کنید.
