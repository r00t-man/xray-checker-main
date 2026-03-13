---
title: Advanced Configuration
description: Advanced configuration options
---

## Advanced Configuration

### Custom IP Check Services

You can use alternative IP check services (see [check methods](/configuration/check-methods) for details):

- `http://ip.sb`
- `https://api64.ipify.org`
- `http://ifconfig.me`

Example:

```bash
PROXY_IP_CHECK_URL=http://ip.sb
```

### Custom Status Check URLs

Alternative URLs for status checking (see [check methods](/configuration/check-methods) for details):

- `http://www.gstatic.com/generate_204`
- `http://www.qualcomm.cn/generate_204`
- `http://cp.cloudflare.com/generate_204`

Example:

```bash
PROXY_STATUS_CHECK_URL=http://www.gstatic.com/generate_204
```

### Security Configuration

Enable authentication for sensitive endpoints:

```bash
METRICS_PROTECTED=true
METRICS_USERNAME=custom_user
METRICS_PASSWORD=secure_password
```

### Instance Labeling

Add instance labels for distributed setups:

```bash
METRICS_INSTANCE=datacenter-1
```

### Update Intervals

Customize check and update intervals:

```bash
# Check every minute
PROXY_CHECK_INTERVAL=60

# Update subscription every hour
SUBSCRIPTION_UPDATE_INTERVAL=3600
```

### Logging Configuration

Adjust Xray Core logging:

```bash
# Enable debug logging
XRAY_LOG_LEVEL=debug

# Disable logging
XRAY_LOG_LEVEL=none
```

### Port Configuration

Customize port ranges:

```bash
# Start SOCKS5 ports from 20000
XRAY_START_PORT=20000

# Change metrics port
METRICS_PORT=9090
```

### Configuration for steal-from-yourself domain

You have your own domain, your-domain.com, with a website running on it,
and you want to display monitoring at `your-domain.com/xray/monitor`.

Run Xray Checker on the same server where your website is hosted.
The parameter `-p 127.0.0.1:2112:2112` ensures that direct access
to it is only possible from the server itself:

:::caution
If the web interface is publicly accessible, it's recommended to use basic auth for protection.
You can enable this using the following environment variables:
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

Open nginx configuration file (`sudo nano /etc/nginx/your-domain.com`), find main section:

```nginx
server {
    root /var/www/your-domain.com/html;

    index index.html;
    server_name your-stealing-domain.com;
    ...
}
```

and paste there 3 new locations:

```nginx

    # Handle /xray/monitor without a trailing slash
    location = /xray/monitor {
        return 301 https://$host$request_uri/;
    }

    # Handle /xray/monitor/ - redirect to xray checker's docker port
    location /xray/monitor/ {
        proxy_pass http://127.0.0.1:2112/xray/monitor/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
```

then check and reload nginx:

```bash
sudo nginx -t
sudo systemctl reload nginx
```

and check availability of monitoring:

```bash
 curl -I -L https://your-domain.com/xray/monitor
```
