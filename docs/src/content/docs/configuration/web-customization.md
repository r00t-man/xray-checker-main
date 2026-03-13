---
title: Web Customization
description: Customize the web interface with your own templates, styles, and assets
---

Xray Checker allows full customization of the web interface. You can replace the default template, add custom styles, change the logo, favicon, and add any other static files.

## Enabling Custom Assets

Set the path to your custom assets directory:

```bash
# Environment variable
WEB_CUSTOM_ASSETS_PATH=/path/to/custom

# CLI flag
xray-checker --web-custom-assets-path=/path/to/custom
```

If the path is set and the directory exists, custom assets will be loaded at startup.

## Directory Structure

Place your custom files in a flat directory (no subdirectories):

```
custom/
  ├── index.html       # Full template replacement (optional)
  ├── logo.svg         # Single logo for both themes (optional)
  ├── logo.png         # Single logo PNG (optional)
  ├── logo-dark.svg    # Logo for dark theme (optional)
  ├── logo-dark.png    # Logo for dark theme PNG (optional)
  ├── logo-light.svg   # Logo for light theme (optional)
  ├── logo-light.png   # Logo for light theme PNG (optional)
  ├── favicon.ico      # Favicon override (optional)
  ├── custom.css       # Additional styles, auto-injected (optional)
  └── any-file.ext     # Available at /static/any-file.ext
```

### Logo Files

You can customize the logo in two ways:

1. **Single logo** — provide `logo.svg` or `logo.png` to use the same logo for both dark and light themes
2. **Theme-specific logos** — provide `logo-dark.svg`/`logo-dark.png` and `logo-light.svg`/`logo-light.png` for different logos per theme

Priority order (first found is used):
1. `logo-dark.svg` / `logo-light.svg` (theme-specific SVG)
2. `logo-dark.png` / `logo-light.png` (theme-specific PNG)
3. `logo.svg` (universal SVG)
4. `logo.png` (universal PNG)

## Custom Styles (custom.css)

The easiest way to customize the appearance. If `custom.css` exists, it will be automatically injected after the default styles.

### CSS Variables

Override theme colors using CSS variables:

```css
:root {
  /* Background colors */
  --bg-primary: #0a0a0f;
  --bg-secondary: #12121a;
  --bg-tertiary: #1a1a24;

  /* Text colors */
  --text-primary: #f4f4f5;
  --text-secondary: #a1a1aa;
  --text-muted: #71717a;

  /* Accent colors */
  --color-green: #22c55e;
  --color-red: #ef4444;
  --color-orange: #f97316;

  /* Border */
  --border: #27272a;
}
```

### Example: Light Theme Override

```css
:root {
  --bg-primary: #ffffff;
  --bg-secondary: #f8f9fa;
  --bg-tertiary: #e9ecef;
  --text-primary: #212529;
  --text-secondary: #495057;
  --text-muted: #6c757d;
  --border: #dee2e6;
}
```

### Example: Custom Logo Size

```css
.header-logo img {
  width: 48px;
  height: 48px;
}
```

## Full Template Replacement (index.html)

For complete control, provide your own `index.html`. This is a Go template with access to all page data.

:::caution
Custom templates may break after updates if the data structure changes. Use at your own risk.
:::

### Go Template Syntax

```html
{{ .Variable }}           <!-- Output variable -->
{{ if .Condition }}...{{ end }}
{{ range .Array }}...{{ end }}
{{ formatLatency .Latency }}  <!-- Format duration as "123ms" or "n/a" -->
```

### Available Variables

#### PageData (root object)

| Variable | Type | Description |
|----------|------|-------------|
| `.Version` | string | Xray Checker version |
| `.Host` | string | Server host |
| `.Port` | string | Server port |
| `.CheckInterval` | int | Proxy check interval in seconds |
| `.Timeout` | int | Proxy check timeout in seconds |
| `.CheckMethod` | string | Check method: `ip`, `status`, or `download` |
| `.IPCheckUrl` | string | URL for IP checking |
| `.StatusCheckUrl` | string | URL for status checking |
| `.DownloadUrl` | string | URL for download checking |
| `.SimulateLatency` | bool | Whether latency simulation is enabled |
| `.SubscriptionUpdate` | bool | Whether subscription auto-update is enabled |
| `.SubscriptionUpdateInterval` | int | Subscription update interval in seconds |
| `.StartPort` | int | First proxy port number |
| `.Instance` | string | Metrics instance label |
| `.PushUrl` | string | Prometheus pushgateway URL |
| `.ShowServerDetails` | bool | Whether to show server IPs and ports |
| `.IsPublic` | bool | Whether public mode is enabled |
| `.SubscriptionName` | string | Subscription name for display |
| `.Endpoints` | []EndpointInfo | Array of proxy endpoints |

#### EndpointInfo (each item in `.Endpoints`)

| Variable | Type | Availability | Description |
|----------|------|--------------|-------------|
| `.Name` | string | Always | Proxy name |
| `.StableID` | string | Always | Unique proxy identifier |
| `.Index` | int | Always | Proxy index (0-based) |
| `.Status` | bool | Always | `true` if online |
| `.Latency` | time.Duration | Always | Response latency |
| `.ServerInfo` | string | When `ShowServerDetails && !IsPublic` | Server address and port |
| `.ProxyPort` | int | When `ShowServerDetails && !IsPublic` | Local proxy port |
| `.URL` | string | When `!IsPublic` | Config status endpoint URL |

### Template Functions

| Function | Description | Example |
|----------|-------------|---------|
| `formatLatency` | Formats duration as milliseconds | `{{ formatLatency .Latency }}` → `"123ms"` or `"n/a"` |

### Conditional Rendering

```html
<!-- Show only in non-public mode -->
{{ if not .IsPublic }}
  <a href="{{ .URL }}">Config</a>
{{ end }}

<!-- Show server details when enabled -->
{{ if .ShowServerDetails }}
  <span>{{ .ServerInfo }}</span>
{{ end }}

<!-- Loop through proxies -->
{{ range .Endpoints }}
  <div class="{{ if .Status }}online{{ else }}offline{{ end }}">
    {{ .Name }} - {{ formatLatency .Latency }}
  </div>
{{ end }}
```

### Minimal Template Example

```html
<!DOCTYPE html>
<html>
<head>
  <title>{{ if .SubscriptionName }}{{ .SubscriptionName }}{{ else }}Status{{ end }}</title>
</head>
<body>
  <h1>Proxy Status</h1>
  {{ range .Endpoints }}
  <div>
    {{ if .Status }}✓{{ else }}✗{{ end }}
    {{ .Name }}
    ({{ formatLatency .Latency }})
  </div>
  {{ end }}
</body>
</html>
```

## Docker Example

```yaml
services:
  xray-checker:
    image: kutovoys/xray-checker:latest
    volumes:
      - ./custom:/app/custom:ro
    environment:
      - WEB_CUSTOM_ASSETS_PATH=/app/custom
      - SUBSCRIPTION_URL=https://...
```

Directory structure:
```
my-project/
  ├── docker-compose.yml
  └── custom/
      ├── logo.svg        # or logo-dark.svg + logo-light.svg
      ├── favicon.ico
      └── custom.css
```

## Startup Logs

When custom assets are loaded, you'll see:

```
INFO  Custom assets enabled: /app/custom
INFO  Custom assets loaded:
INFO    ✓ logo.svg
INFO    ✓ custom.css
INFO  Using default template
```

Or with custom template:
```
INFO  Custom assets loaded:
INFO    ✓ index.html
INFO    ✓ custom.css
INFO  Using custom template: index.html
```

## Errors

| Error | Cause |
|-------|-------|
| `custom assets directory does not exist` | Path set but directory not found |
| `failed to parse custom template` | Invalid Go template syntax in index.html |

The application will **not start** if custom assets path is set but invalid.
