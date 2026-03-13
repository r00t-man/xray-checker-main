---
title: Troubleshooting
description: Troubleshooting
tableOfContents:
  minHeadingLevel: 2
  maxHeadingLevel: 4
---

## Common Issues

### Subscription Problems

#### Invalid Subscription URL Response

```
error parsing subscription: error getting subscription: unexpected status code: 403
```

**Possible causes:**

- URL is incorrect
- URL is no longer valid
- Server blocks Xray Checker user agent

**Solutions:**

1. Verify subscription URL
2. Check if URL is still active
3. Contact subscription provider
4. Try using Base64 format directly instead of URL

#### Base64 Decode Failed

```
error decoding Base64: illegal base64 data...
```

**Possible causes:**

- Invalid Base64 encoding
- URL-safe vs standard Base64
- Additional whitespace or newlines

**Solutions:**

1. Verify Base64 string is clean without whitespace
2. Try URL-safe Base64 decode if standard fails
3. Check if content needs to be decoded multiple times

### Proxy Check Issues

#### Running on Proxy Server

When running Xray Checker on the same server where your proxies are hosted, you **must** use the `status` check method instead of the default `ip` method.

**Why:**

- The `ip` check method compares your IP with and without proxy
- When running on the proxy server, both IPs will be the same
- This causes false negatives - working proxies reported as failed

**Solution:**

```bash
# In environment
PROXY_CHECK_METHOD=status
PROXY_STATUS_CHECK_URL=http://cp.cloudflare.com/generate_204

# Or via CLI
--proxy-check-method=status --proxy-status-check-url="http://cp.cloudflare.com/generate_204"
```

#### All Proxies Failing

```
Warning: error parsing proxy URL: connection refused
```

**Possible causes:**

- Network connectivity issues
- Firewall blocking connections
- IP check service unavailable

**Solutions:**

1. Check network connectivity
2. Verify firewall rules
3. Try alternative check method:
   ```bash
   PROXY_CHECK_METHOD=status
   ```
4. Use alternative IP check service:
   ```bash
   PROXY_IP_CHECK_URL=http://ip.sb
   ```

#### High Latency or Timeouts

```
Warning: error getting current IP: timeout
```

**Possible causes:**

- Slow network connection
- IP check service slow
- Proxy timeout too low

**Solutions:**

1. Increase timeout:
   ```bash
   PROXY_TIMEOUT=60
   ```
2. Use faster IP check service
3. Disable latency simulation:
   ```bash
   SIMULATE_LATENCY=false
   ```

### Metrics Issues

#### Cannot Access Metrics

```
Error: Unauthorized
```

**Possible causes:**

- Authentication enabled
- Incorrect credentials
- Wrong port

**Solutions:**

1. Check if authentication is enabled:
   ```bash
   METRICS_PROTECTED=false
   ```
2. Verify credentials:
   ```bash
   METRICS_USERNAME=user
   METRICS_PASSWORD=pass
   ```
3. Verify correct port:
   ```bash
   METRICS_PORT=2112
   ```

#### Pushgateway Errors

```
Error pushing metrics: unexpected status code 401
```

**Possible causes:**

- Invalid pushgateway URL
- Authentication required
- Network issues

**Solutions:**

1. Check URL format:
   ```bash
   METRICS_PUSH_URL="http://user:pass@host:9091"
   ```
2. Verify network connectivity
3. Check pushgateway logs

### Port Conflicts

#### Port Already in Use

```
error starting server: listen tcp :2112: bind: address already in use
```

**Possible causes:**

- Another service using the port
- Previous instance still running
- System port restrictions

**Solutions:**

1. Change metrics port:
   ```bash
   METRICS_PORT=2113
   ```
2. Check for running processes:
   ```bash
   lsof -i :2112
   ```
3. Stop conflicting services

#### SOCKS Port Range Issues

```
error starting Xray: port already in use
```

**Possible causes:**

- Port range conflict
- Too many proxies
- System port limits

**Solutions:**

1. Change start port:
   ```bash
   XRAY_START_PORT=20000
   ```
2. Check system limits:
   ```bash
   ulimit -n
   ```
3. Free up port range

## Debugging Techniques

### Enable Debug Logging

There are two log level settings:

**Application Log Level** - controls Xray Checker's own logging:

```bash
LOG_LEVEL=debug
```

Available levels: `debug`, `info`, `warn`, `error`

Debug log will show:
- Subscription parsing details
- Proxy check results
- Configuration loading
- HTTP request/response info

**Xray Core Log Level** - controls the embedded Xray proxy engine:

```bash
XRAY_LOG_LEVEL=debug
```

Available levels: `debug`, `info`, `warning`, `error`, `none`

Xray debug log will show:
- Connection attempts
- Protocol negotiation
- Transport layer details
- Timing information

:::tip[Recommended Debug Setup]
For troubleshooting, enable both:
```bash
LOG_LEVEL=debug
XRAY_LOG_LEVEL=warning
```
This gives detailed app logs while keeping Xray output manageable.
:::

### Check Process Status

```bash
# Check if process is running
ps aux | grep xray-checker

# Check open ports
netstat -tulpn | grep xray-checker
```

### Verify Network Connectivity

```bash
# Test IP check service
curl -v https://api.ipify.org?format=text

# Test proxy connection
curl --socks5 localhost:10000 -v https://api.ipify.org?format=text
```

### Docker Debugging

```bash
# Check container logs
docker logs xray-checker

# Access container shell
docker exec -it xray-checker sh

# Check container network
docker inspect xray-checker
```

## Getting Help

If you're still experiencing issues:

1. Check GitHub Issues for similar problems
2. Create new issue with:

   - Full error message
   - Configuration used
   - Debug logs
   - Steps to reproduce

3. Include environment details:
   - OS version
   - Docker version (if using)
   - Xray Checker version
