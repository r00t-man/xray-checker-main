---
title: Uptime Kuma
description: Deployment options and examples
---

Xray Checker provides individual status endpoints for each proxy, perfect for Uptime Kuma monitoring.

### Basic Setup

1. Open Uptime Kuma
2. Click "Add Monitor"
3. Select "HTTP(s)"
4. Configure monitor:
   - Name: Proxy name
   - URL: `http://localhost:2112/config/a1b2c3d4e5f67890`
   - Interval: Your preferred check interval
   - Retry: Recommended 3 times

### With Authentication

If metrics protection is enabled:

URL format:

```
http://username:password@localhost:2112/config/a1b2c3d4e5f67890
```

Configuration:

- Authentication Method: Basic Auth
- Username: Your METRICS_USERNAME
- Password: Your METRICS_PASSWORD

### Status Codes

- 200: Proxy working
- 503: Proxy failed
- 401: Authentication required
- 403: Authentication failed
