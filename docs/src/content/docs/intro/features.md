---
title: Features
description: Xray Checker Features
tableOfContents: false
---

### 🚀 Core Features

- 🔍 Monitor the health of Xray proxy servers with support for various protocols (VLESS, VMess, Trojan, Shadowsocks)

- 🔄 Automatic proxy configuration updates from subscription URLs with [configurable intervals](/configuration/envs#subscription_update_interval)

- 📊 [Export metrics](/integrations/metrics) in Prometheus format with proxy status and latency information

- 🌓 Web interface with dark/light theme for monitoring all proxy endpoints status

  - 🔍 Search and filter proxies by name or status
  - 📊 Sort by name, latency, or status
  - 🔄 Auto-refresh without page reload
  - 🎨 [Full customization](/configuration/web-customization) — custom logo, styles, or entire template

- 🌐 [REST API](/usage/api-reference) with OpenAPI/Swagger documentation

### 📝 Formats and Configuration

- 📋 [Support for various configuration formats](/configuration/subscription):

  - 🔗 URL subscriptions (with multiple URL support)
  - 🔐 Base64-encoded strings
  - 📄 V2Ray/Xray JSON files
  - 📦 Xray JSON array (multi-config)
  - 📁 Configuration folders

### 🔌 Integrations

- 🌐 [REST API](/usage/api-reference) with OpenAPI/Swagger documentation for custom integrations

- 📄 [Public status page](/configuration/status-page) for VPN services — display proxy status without authentication, customizable title from subscription name

- 📥 [Automatic endpoint generation](/integrations/uptime-kuma) for integration with monitoring systems (e.g., Uptime-Kuma)

- ⏱️ [Latency simulation](/configuration/advanced-conf) for endpoints to ensure accurate monitoring system testing

- 📡 [Integration with Prometheus Pushgateway](/integrations/prometheus#pushgateway-integration) for sending metrics to external monitoring systems

### ⚡ Check Methods

- 🔧 [Support for three proxy verification methods](/configuration/check-methods):

  - 🌐 Via IP address comparison
  - ✅ Via HTTP status checks
  - 📥 Via file download verification

- ⏱️ Accurate latency measurement using TTFB (Time To First Byte)

### 🔒 Security

- 🛡️ [Protect metrics and web interface](/configuration/advanced-conf#security-settings) using Basic Authentication

### 🚀 Deployment

- 🐳 Can be run both in a [Docker container](/usage/docker) (including Docker Compose) and as a [standalone CLI application](/usage/cli)

:::tip[💡 Quick Start]
To start using Xray Checker right now, go to the [Quick Start](/intro/quick-start) section
:::
