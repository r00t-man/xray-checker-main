---
title: Badge System
description: Embed proxy status indicators anywhere
---

Xray Checker includes a powerful badge system that allows you to embed real-time proxy status indicators anywhere — status pages, dashboards, documentation, or README files.

## Features

- **Real-time status** — Badges update automatically with proxy state
- **Multiple themes** — Dark and light themes to match your design
- **Customizable styles** — Various variants, sizes, and rounded corners
- **Flexible display** — Show/hide name and latency independently
- **Easy embedding** — Works via iframe or direct URL

## Quick Start

The simplest badge URL:

```
https://your-xray-checker.com/?stableId=abc123def456
```

This displays a badge showing the proxy status, name, and latency.

## Getting the Stable ID

Each proxy has a unique `stableId` that persists across restarts. Get it from the API:

```bash
curl https://your-xray-checker.com/api/v1/public/proxies
```

Response:

```json
{
  "success": true,
  "data": [
    {
      "stableId": "a1b2c3d4e5f67890",
      "name": "US-Server-1",
      "online": true,
      "latencyMs": 150
    }
  ]
}
```

## Available Parameters

| Parameter | Values | Default | Description |
|-----------|--------|---------|-------------|
| `stableId` | `{id}` | required | Proxy stable ID from API |
| `theme` | `dark`, `light` | `dark` | Color theme |
| `variant` | `default`, `flat`, `pill`, `dot` | `default` | Badge style |
| `size` | `sm`, `md`, `lg` | `md` | Badge size |
| `rounded` | `none`, `sm`, `md`, `lg`, `full` | `md` | Corner rounding |
| `showName` | `true`, `false` | `true` | Show proxy name |
| `showLatency` | `true`, `false` | `true` | Show latency value |
| `width` | number | auto | Custom width in pixels |
| `height` | number | auto | Custom height in pixels |

## Badge Variants

### Default
Standard badge with background and border.
```
?stableId=abc123
```

### Flat
Minimal badge without background or border.
```
?stableId=abc123&variant=flat
```

### Pill
Rounded pill-shaped badge.
```
?stableId=abc123&variant=pill
```

### Dot Only
Just the status indicator dot.
```
?stableId=abc123&variant=dot
```

## Embedding Examples

### HTML iframe

```html
<iframe
  src="https://your-server.com/?stableId=abc123&theme=light"
  width="200"
  height="50"
  frameborder="0">
</iframe>
```

### Multiple Badges

Create a status dashboard by combining multiple badges:

```html
<div style="display: flex; gap: 10px;">
  <iframe src="https://your-server.com/?stableId=server1" width="200" height="50" frameborder="0"></iframe>
  <iframe src="https://your-server.com/?stableId=server2" width="200" height="50" frameborder="0"></iframe>
  <iframe src="https://your-server.com/?stableId=server3" width="200" height="50" frameborder="0"></iframe>
</div>
```

## Next Steps

- Use the [Badge Playground](/badges/playground) to interactively build your badge URL
- Learn about the [Public Status Page](/configuration/status-page) for sharing proxy status
