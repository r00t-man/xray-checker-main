---
title: GitHub Actions
description: Run Xray Checker using GitHub Actions
---

# GitHub Actions Integration

You can run Xray Checker using GitHub Actions. This approach is useful when you need to run checks from different locations or don't have your own server.

## Quick Setup

1. Fork the [xray-checker-in-actions](https://github.com/kutovoys/xray-checker-in-actions) repository
2. Configure the following secrets in your forked repository:
   - `SUBSCRIPTION_URL`: Your subscription URL
   - `PUSH_URL`: Prometheus pushgateway URL for metrics collection
   - `INSTANCE`: (Optional) Instance name for metrics identification

The Action will:

- Run every 5 minutes
- Use the latest version of Xray Checker
- Push metrics to your Prometheus pushgateway
- Run with `--run-once` flag to ensure clean execution

This method requires a properly configured Prometheus pushgateway as it can't expose metrics directly. The metrics will be pushed to your specified `PUSH_URL` with the instance label from your configuration.

## Advanced Configurations

If you need more control over your GitHub Actions setup, here are some advanced configurations you can use.

### Multiple Region Setup

Run checks from different regions simultaneously:

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

### Error Notifications

Add Slack or Email notifications for failed checks:

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

### Custom Check Intervals

Different schedule patterns based on your needs:

```yaml
on:
  schedule:
    - cron: "*/5 * * * *" # Every 5 minutes (default)
    - cron: "0 * * * *" # Every hour
    - cron: "0 */2 * * *" # Every 2 hours
```

### Resource Optimization

Optimize GitHub Actions usage with concurrency controls:

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

## Monitoring Setup

### Required Prometheus Configuration

To collect metrics from the pushgateway, add this to your Prometheus configuration:

```yaml
scrape_configs:
  - job_name: "pushgateway"
    honor_labels: true
    static_configs:
      - targets: ["pushgateway:9091"]
```

The metrics will appear with the instance label you specified in your GitHub Actions configuration, allowing you to track checks from different locations.
