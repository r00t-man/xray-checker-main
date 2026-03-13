---
title: Public Status Page
description: Share proxy status with your users
---

Xray Checker includes a built-in public status page that allows you to share proxy availability information with your users without exposing sensitive configuration details.

## Overview

The public status page provides:

- Real-time proxy status (online/offline)
- Latency information for each proxy
- Clean, responsive interface
- Automatic updates

When enabled, sensitive information like server addresses, ports, and configuration details are hidden from view.

## Enabling Public Mode

Set the environment variable:

```bash
WEB_PUBLIC=true
```

Or use the CLI flag:

```bash
xray-checker --web-public
```

## What's Hidden in Public Mode

| Information           | Public Mode | Normal Mode |
| --------------------- | ----------- | ----------- |
| Proxy Name            | ✓ Visible   | ✓ Visible   |
| Status                | ✓ Visible   | ✓ Visible   |
| Latency               | ✓ Visible   | ✓ Visible   |
| Server Address        | ✗ Hidden    | ✓ Visible   |
| Port                  | ✗ Hidden    | ✓ Visible   |
| Local Proxy Port      | ✗ Hidden    | ✓ Visible   |
| Configuration Details | ✗ Hidden    | ✓ Visible   |

## Custom Subscription Name

You can set a custom name for your status page using the subscription info from your provider, or set it manually:

```bash
# Will be displayed as the page title
SUBSCRIPTION_URL=https://provider.com/sub#MyVPN
```

## URL Customization

The status page supports URL parameters to customize the view:

| Parameter                | Description                        |
| ------------------------ | ---------------------------------- |
| `hideHeader=true`        | Hide header with logo and controls |
| `hideStats=true`         | Hide statistics cards              |
| `hideServersHeader=true` | Hide "Servers" heading             |
| `hideControls=true`      | Hide search and filters            |
| `hideStatusInfo=true`    | Hide technical info footer         |
| `hideFooter=true`        | Hide page footer                   |
| `hideBackground=true`    | Transparent background             |

### Example: Minimal Embed

```
https://your-server.com/?hideHeader=true&hideFooter=true&hideControls=true
```

## Use Cases

### Sharing with Users

Provide your users with a link to check proxy availability:

```
https://status.yourdomain.com/
```

### Embedding in Documentation

Use iframe to embed the status page:

```html
<iframe
  src="https://status.yourdomain.com/?hideHeader=true&hideFooter=true"
  width="100%"
  height="400"
  frameborder="0"
>
</iframe>
```

### Integration with Status Services

The public page works well behind reverse proxies like Nginx or Cloudflare for additional security and caching.
