---
title: Quick Start
description: Quick Start of Xray Checker
---

Get Xray Checker up and running in minutes with these simple steps.

## Prerequisites

- Subscription URL for your proxies
- Docker (optional, for container deployment)
- Prometheus (optional, for metrics collection)

## 5-Minute Setup

### Using Docker (Recommended)

1. Pull the image:

```bash
docker pull kutovoys/xray-checker
```

2. Run with basic configuration:

```bash
docker run -d \
  -e SUBSCRIPTION_URL=https://your-subscription-url/sub \
  -p 2112:2112 \
  kutovoys/xray-checker
```

3. Check the status:

```bash
curl http://localhost:2112/health
```

### Using Binary

1. Download the latest release:

```bash
# Linux amd64
curl -sL -o - $(curl -s https://api.github.com/repos/kutovoys/xray-checker/releases/latest | grep "browser_download_url.*linux-amd64.tar.gz" | cut -d'"' -f4) | tar -xz
chmod +x xray-checker

# Linux arm64
curl -sL -o - $(curl -s https://api.github.com/repos/kutovoys/xray-checker/releases/latest | grep "browser_download_url.*linux-arm64.tar.gz" | cut -d'"' -f4) | tar -xz
chmod +x xray-checker
```

2. Run with basic configuration:

```bash
./xray-checker --subscription-url=https://your-subscription-url/sub
```

## Verify Installation

1. Open web interface:

   - Navigate to `http://localhost:2112`
   - You should see the dashboard with proxy status

2. Check metrics:

   - Navigate to `http://localhost:2112/metrics`
   - You should see Prometheus metrics

3. Verify proxy status:
   - Click on any proxy link in the web interface
   - Check the status endpoint response

## Next Steps

1. Configure Prometheus:

```yaml
scrape_configs:
  - job_name: "xray-checker"
    static_configs:
      - targets: ["localhost:2112"]
```

2. Set up Uptime Kuma:

   - Add new monitor
   - Use proxy-specific endpoints
   - Configure alerts

3. Customize configuration:
   - Adjust check intervals
   - Configure authentication
   - Set up metric pushing

## Common Commands

Check version:

```bash
./xray-checker --version
```

Run single check:

```bash
./xray-checker --subscription-url=https://your-sub-url --run-once
```

Enable authentication:

```bash
./xray-checker --subscription-url=https://your-sub-url \
  --metrics-protected=true \
  --metrics-username=user \
  --metrics-password=pass
```

## Troubleshooting

1. Check service status:

```bash
curl http://localhost:2112/health
```

2. View logs:

```bash
docker logs xray-checker
```

3. Verify metrics:

```bash
curl http://localhost:2112/metrics
```

## Need Help?

- Check the full documentation
- Open an issue on GitHub
- Join the community discussions
