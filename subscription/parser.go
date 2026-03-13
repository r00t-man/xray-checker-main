package subscription

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"xray-checker/config"
	"xray-checker/logger"
	"xray-checker/models"
)

type Parser struct{}

type fetchResult struct {
	Content []byte
	Name    string
}

func NewParser() *Parser {
	return &Parser{}
}

type libXrayOutbound struct {
	Protocol       string                 `json:"protocol"`
	SendThrough    string                 `json:"sendThrough"`
	Tag            string                 `json:"tag"`
	Settings       *libXraySettings       `json:"settings"`
	StreamSettings *libXrayStreamSettings `json:"streamSettings"`
}

type libXraySettings struct {
	Address    string `json:"address"`
	Port       int    `json:"port"`
	Level      int    `json:"level"`
	ID         string `json:"id"`
	Flow       string `json:"flow"`
	Encryption string `json:"encryption"`
	AlterId    int    `json:"alterId"`
	Security   string `json:"security"`
	Password   string `json:"password"`
	Method     string `json:"method"`
}

type libXrayStreamSettings struct {
	Network             string                      `json:"network"`
	Security            string                      `json:"security"`
	TlsSettings         *libXrayTlsSettings         `json:"tlsSettings"`
	RealitySettings     *libXrayRealitySettings     `json:"realitySettings"`
	RawSettings         *libXrayRawSettings         `json:"rawSettings"`
	WsSettings          *libXrayWsSettings          `json:"wsSettings"`
	GrpcSettings        *libXrayGrpcSettings        `json:"grpcSettings"`
	HttpSettings        *libXrayHttpSettings        `json:"httpSettings"`
	HttpupgradeSettings *libXrayHttpupgradeSettings `json:"httpupgradeSettings"`
	XhttpSettings       json.RawMessage             `json:"xhttpSettings"`
	SplithttpSettings   json.RawMessage             `json:"splithttpSettings"`
}

type libXrayTlsSettings struct {
	ServerName    string   `json:"serverName"`
	AllowInsecure bool     `json:"allowInsecure"`
	Fingerprint   string   `json:"fingerprint"`
	Alpn          []string `json:"alpn"`
}

type libXrayRealitySettings struct {
	ServerName  string `json:"serverName"`
	Fingerprint string `json:"fingerprint"`
	PublicKey   string `json:"publicKey"`
	ShortId     string `json:"shortId"`
}

type libXrayRawSettings struct {
	Header *struct {
		Type    string `json:"type"`
		Request *struct {
			Path    []string `json:"path"`
			Headers *struct {
				Host []string `json:"Host"`
			} `json:"headers"`
		} `json:"request"`
	} `json:"header"`
}

type libXrayWsSettings struct {
	Path    string `json:"path"`
	Headers *struct {
		Host string `json:"Host"`
	} `json:"headers"`
	Host string `json:"host"`
}

type libXrayGrpcSettings struct {
	ServiceName string `json:"serviceName"`
	MultiMode   bool   `json:"multiMode"`
}

type libXrayHttpSettings struct {
	Path string   `json:"path"`
	Host []string `json:"host"`
}

type libXrayHttpupgradeSettings struct {
	Path string `json:"path"`
	Host string `json:"host"`
}

type libXrayXhttpSettings struct {
	Path string `json:"path"`
	Host string `json:"host"`
	Mode string `json:"mode"`
}

type originalLinkData struct {
	Name          string
	Encryption    string
	Type          string
	Path          string
	Host          string
	AllowInsecure bool
}

type parsedLink struct {
	Server        string
	Port          int
	Name          string
	Encryption    string
	Type          string
	Path          string
	Host          string
	AllowInsecure bool
}

type xrayStandardSettings struct {
	Vnext []struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
		Users   []struct {
			ID         string `json:"id"`
			Flow       string `json:"flow"`
			Encryption string `json:"encryption"`
			AlterId    int    `json:"alterId"`
			Security   string `json:"security"`
			Level      int    `json:"level"`
		} `json:"users"`
	} `json:"vnext"`
	Servers []struct {
		Address  string `json:"address"`
		Port     int    `json:"port"`
		Password string `json:"password"`
		Method   string `json:"method"`
		Flow     string `json:"flow"`
	} `json:"servers"`
}

type ParseResult struct {
	Configs []*models.ProxyConfig
	Name    string
}

func (p *Parser) Parse(subscriptionData string) (*ParseResult, error) {
	sourceType := p.detectSourceType(subscriptionData)
	logger.Debug("Detected source type: %s", sourceType)

	var rawData []byte
	var subName string
	var err error

	switch sourceType {
	case "url":
		result, fetchErr := p.fetchURLContent(subscriptionData)
		if fetchErr != nil {
			return nil, fmt.Errorf("failed to fetch URL content: %v", fetchErr)
		}
		rawData = result.Content
		subName = result.Name
	case "folder":
		folderPath := strings.TrimPrefix(subscriptionData, "folder://")
		configs, folderErr := p.parseFolder(folderPath)
		if folderErr != nil {
			return nil, folderErr
		}
		return &ParseResult{Configs: configs, Name: ""}, nil
	case "file":
		filePath := strings.TrimPrefix(subscriptionData, "file://")
		rawData, err = os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %v", err)
		}
	case "base64":
		rawData = []byte(strings.TrimPrefix(subscriptionData, "base64://"))
		rawData = []byte(strings.TrimPrefix(string(rawData), "data:text/plain;base64,"))
	default:
		rawData = []byte(subscriptionData)
	}

	trimmedData := strings.TrimSpace(string(rawData))
	if strings.HasPrefix(trimmedData, "[") {
		logger.Debug("Detected JSON array format")
		configs, jsonErr := p.parseJSONConfigs(rawData)
		if jsonErr != nil {
			return nil, jsonErr
		}
		return &ParseResult{Configs: configs, Name: subName}, nil
	}

	if strings.HasPrefix(trimmedData, "{") {
		logger.Debug("Detected single JSON object format")
		configs, jsonErr := p.parseSingleJSONConfig(rawData)
		if jsonErr != nil {
			return nil, jsonErr
		}
		return &ParseResult{Configs: configs, Name: subName}, nil
	}

	originalData := p.parseOriginalLinks(rawData)

	cleanedData := p.cleanEmptyLines(rawData)

	outbounds := p.parseShareLinksOutbounds(cleanedData)
	if len(outbounds) == 0 {
		return nil, fmt.Errorf("failed to parse subscription links: no supported share links found")
	}

	logger.Debug("Parsed %d outbounds", len(outbounds))

	var proxyConfigs []*models.ProxyConfig
	configIndex := 0
	for _, outboundRaw := range outbounds {
		proxyConfig, err := p.convertOutbound(outboundRaw, configIndex, originalData)
		if err != nil {
			logger.Debug("Skipping outbound: %v", err)
			continue
		}
		if proxyConfig != nil {
			proxyConfigs = append(proxyConfigs, proxyConfig)
			configIndex++
		}
	}

	if len(proxyConfigs) == 0 {
		return nil, fmt.Errorf("no valid proxy configurations found")
	}

	return &ParseResult{Configs: proxyConfigs, Name: subName}, nil
}

func (p *Parser) parseShareLinksOutbounds(rawData []byte) []json.RawMessage {
	decoded := p.tryDecodeBase64(rawData)
	lines := strings.Split(string(decoded), "\n")
	outbounds := make([]json.RawMessage, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		outbound := p.parseShareLinkToOutbound(line)
		if outbound == nil {
			continue
		}

		raw, err := json.Marshal(outbound)
		if err != nil {
			logger.Debug("Skipping outbound due to marshal error: %v", err)
			continue
		}
		outbounds = append(outbounds, raw)
	}

	return outbounds
}

func (p *Parser) parseShareLinkToOutbound(link string) map[string]interface{} {
	if strings.HasPrefix(link, "vmess://") {
		return p.parseVMessOutbound(link)
	}

	if strings.HasPrefix(link, "vless://") {
		return p.parseVLESSOutbound(link)
	}

	if strings.HasPrefix(link, "trojan://") {
		return p.parseTrojanOutbound(link)
	}

	if strings.HasPrefix(link, "ss://") {
		return p.parseShadowsocksOutbound(link)
	}

	return nil
}

func (p *Parser) parseVLESSOutbound(link string) map[string]interface{} {
	u, err := url.Parse(link)
	if err != nil || u.Hostname() == "" || u.Port() == "" || u.User == nil {
		return nil
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil || port == 0 {
		return nil
	}
	id := u.User.Username()
	if id == "" {
		return nil
	}

	q := u.Query()
	encryption := q.Get("encryption")
	if encryption == "" {
		encryption = "none"
	}

	outbound := map[string]interface{}{
		"protocol": "vless",
		"settings": map[string]interface{}{
			"vnext": []map[string]interface{}{{
				"address": u.Hostname(),
				"port":    port,
				"users": []map[string]interface{}{{
					"id":         id,
					"encryption": encryption,
					"flow":       q.Get("flow"),
				}},
			}},
		},
	}
	p.attachStreamSettings(outbound, q)
	return outbound
}

func (p *Parser) parseTrojanOutbound(link string) map[string]interface{} {
	u, err := url.Parse(link)
	if err != nil || u.Hostname() == "" || u.Port() == "" || u.User == nil {
		return nil
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil || port == 0 {
		return nil
	}
	password := u.User.Username()
	if password == "" {
		return nil
	}

	q := u.Query()
	outbound := map[string]interface{}{
		"protocol": "trojan",
		"settings": map[string]interface{}{
			"servers": []map[string]interface{}{{
				"address":  u.Hostname(),
				"port":     port,
				"password": password,
				"flow":     q.Get("flow"),
			}},
		},
	}
	p.attachStreamSettings(outbound, q)
	return outbound
}

func (p *Parser) parseShadowsocksOutbound(link string) map[string]interface{} {
	method, password, server, port, err := p.parseShadowsocksCreds(link)
	if err != nil {
		return nil
	}

	outbound := map[string]interface{}{
		"protocol": "shadowsocks",
		"settings": map[string]interface{}{
			"servers": []map[string]interface{}{{
				"address":  server,
				"port":     port,
				"method":   method,
				"password": password,
			}},
		},
	}
	return outbound
}

func (p *Parser) parseShadowsocksCreds(link string) (string, string, string, int, error) {
	u, err := url.Parse(link)
	if err == nil && u.Hostname() != "" && u.Port() != "" {
		port, pErr := strconv.Atoi(u.Port())
		if pErr == nil && u.User != nil {
			method := u.User.Username()
			password, _ := u.User.Password()
			if method != "" && password != "" {
				return method, password, u.Hostname(), port, nil
			}
		}
	}

	raw := strings.TrimPrefix(link, "ss://")
	if i := strings.Index(raw, "#"); i >= 0 {
		raw = raw[:i]
	}
	if i := strings.Index(raw, "?"); i >= 0 {
		raw = raw[:i]
	}

	parts := strings.SplitN(raw, "@", 2)
	if len(parts) != 2 {
		return "", "", "", 0, fmt.Errorf("invalid ss link")
	}

	left := parts[0]
	right := parts[1]
	if decoded, err := p.decodeBase64(left); err == nil {
		left = string(decoded)
	}

	cred := strings.SplitN(left, ":", 2)
	if len(cred) != 2 {
		return "", "", "", 0, fmt.Errorf("invalid ss credentials")
	}

	host, portStr, ok := strings.Cut(right, ":")
	if !ok {
		return "", "", "", 0, fmt.Errorf("invalid ss host")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port == 0 {
		return "", "", "", 0, fmt.Errorf("invalid ss port")
	}

	return cred[0], cred[1], host, port, nil
}

func (p *Parser) parseVMessOutbound(link string) map[string]interface{} {
	encoded := strings.TrimPrefix(link, "vmess://")
	decoded, err := p.decodeBase64(encoded)
	if err != nil {
		return nil
	}

	var vmess map[string]interface{}
	if err := json.Unmarshal(decoded, &vmess); err != nil {
		return nil
	}

	server, _ := vmess["add"].(string)
	if server == "" {
		return nil
	}

	port := 0
	switch v := vmess["port"].(type) {
	case float64:
		port = int(v)
	case string:
		port, _ = strconv.Atoi(v)
	}
	if port == 0 {
		return nil
	}

	id, _ := vmess["id"].(string)
	if id == "" {
		return nil
	}

	aid := 0
	switch v := vmess["aid"].(type) {
	case float64:
		aid = int(v)
	case string:
		aid, _ = strconv.Atoi(v)
	}

	security := "auto"
	if v, ok := vmess["scy"].(string); ok && v != "" {
		security = v
	}

	outbound := map[string]interface{}{
		"protocol": "vmess",
		"settings": map[string]interface{}{
			"vnext": []map[string]interface{}{{
				"address": server,
				"port":    port,
				"users": []map[string]interface{}{{
					"id":       id,
					"alterId":  aid,
					"security": security,
				}},
			}},
		},
	}

	q := url.Values{}
	if v, ok := vmess["net"].(string); ok {
		q.Set("type", v)
	}
	if v, ok := vmess["host"].(string); ok {
		q.Set("host", v)
	}
	if v, ok := vmess["path"].(string); ok {
		q.Set("path", v)
	}
	if v, ok := vmess["tls"].(string); ok && v != "" && v != "none" {
		q.Set("security", "tls")
	}
	if v, ok := vmess["sni"].(string); ok {
		q.Set("sni", v)
	}

	p.attachStreamSettings(outbound, q)
	return outbound
}

func (p *Parser) attachStreamSettings(outbound map[string]interface{}, q url.Values) {
	network := q.Get("type")
	if network == "" {
		network = "tcp"
	}
	security := q.Get("security")
	if security == "" {
		security = "none"
	}

	stream := map[string]interface{}{
		"network":  network,
		"security": security,
	}

	host := q.Get("host")
	path := q.Get("path")

	if security == "tls" {
		tlsSettings := map[string]interface{}{}
		if sni := q.Get("sni"); sni != "" {
			tlsSettings["serverName"] = sni
		}
		if allow := q.Get("allowInsecure"); allow == "1" || strings.EqualFold(allow, "true") {
			tlsSettings["allowInsecure"] = true
		}
		stream["tlsSettings"] = tlsSettings
	}

	if security == "reality" {
		reality := map[string]interface{}{}
		if sni := q.Get("sni"); sni != "" {
			reality["serverName"] = sni
		}
		if fp := q.Get("fp"); fp != "" {
			reality["fingerprint"] = fp
		}
		if pbk := q.Get("pbk"); pbk != "" {
			reality["publicKey"] = pbk
		}
		if sid := q.Get("sid"); sid != "" {
			reality["shortId"] = sid
		}
		if spx := q.Get("spx"); spx != "" {
			reality["spiderX"] = spx
		}
		stream["realitySettings"] = reality
	}

	switch network {
	case "ws":
		ws := map[string]interface{}{}
		if path != "" {
			ws["path"] = path
		}
		if host != "" {
			ws["headers"] = map[string]interface{}{"Host": host}
		}
		stream["wsSettings"] = ws
	case "grpc":
		grpc := map[string]interface{}{}
		if serviceName := q.Get("serviceName"); serviceName != "" {
			grpc["serviceName"] = serviceName
		}
		if strings.EqualFold(q.Get("mode"), "multi") {
			grpc["multiMode"] = true
		}
		stream["grpcSettings"] = grpc
	case "h2", "http":
		httpSettings := map[string]interface{}{}
		if path != "" {
			httpSettings["path"] = path
		}
		if host != "" {
			httpSettings["host"] = strings.Split(host, ",")
		}
		stream["httpSettings"] = httpSettings
	case "httpupgrade":
		httpUp := map[string]interface{}{}
		if path != "" {
			httpUp["path"] = path
		}
		if host != "" {
			httpUp["host"] = host
		}
		stream["httpupgradeSettings"] = httpUp
	case "xhttp", "splithttp":
		xhttp := map[string]interface{}{}
		if path != "" {
			xhttp["path"] = path
		}
		if host != "" {
			xhttp["host"] = host
		}
		if mode := q.Get("mode"); mode != "" {
			xhttp["mode"] = mode
		}
		stream["xhttpSettings"] = xhttp
	case "tcp":
		if strings.EqualFold(q.Get("headerType"), "http") {
			stream["tcpSettings"] = map[string]interface{}{
				"header": map[string]interface{}{
					"type": "http",
					"request": map[string]interface{}{
						"path":    []string{path},
						"headers": map[string]interface{}{"Host": []string{host}},
					},
				},
			}
		}
	}

	outbound["streamSettings"] = stream
}

func (p *Parser) parseJSONConfigs(data []byte) ([]*models.ProxyConfig, error) {
	var configs []struct {
		Remarks   string            `json:"remarks"`
		Outbounds []json.RawMessage `json:"outbounds"`
	}

	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, fmt.Errorf("failed to parse JSON configs: %v", err)
	}

	logger.Debug("Parsed %d JSON configs", len(configs))

	var proxyConfigs []*models.ProxyConfig
	configIndex := 0

	for _, config := range configs {
		for _, outboundRaw := range config.Outbounds {
			proxyConfig, err := p.convertOutbound(outboundRaw, configIndex, nil)
			if err != nil {
				continue
			}
			if proxyConfig != nil {
				if config.Remarks != "" {
					proxyConfig.Name = config.Remarks
				}
				proxyConfigs = append(proxyConfigs, proxyConfig)
				configIndex++
			}
		}
	}

	if len(proxyConfigs) == 0 {
		return nil, fmt.Errorf("no valid proxy configurations found in JSON")
	}

	return proxyConfigs, nil
}

func (p *Parser) parseSingleJSONConfig(data []byte) ([]*models.ProxyConfig, error) {
	var config struct {
		Remarks   string            `json:"remarks"`
		Outbounds []json.RawMessage `json:"outbounds"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse single JSON config: %v", err)
	}

	logger.Debug("Parsed single JSON config with %d outbounds", len(config.Outbounds))

	var proxyConfigs []*models.ProxyConfig
	configIndex := 0

	for _, outboundRaw := range config.Outbounds {
		proxyConfig, err := p.convertOutbound(outboundRaw, configIndex, nil)
		if err != nil {
			continue
		}
		if proxyConfig != nil {
			if config.Remarks != "" {
				proxyConfig.Name = config.Remarks
			}
			proxyConfigs = append(proxyConfigs, proxyConfig)
			configIndex++
		}
	}

	if len(proxyConfigs) == 0 {
		return nil, fmt.Errorf("no valid proxy configurations found in single JSON config")
	}

	return proxyConfigs, nil
}

func (p *Parser) cleanEmptyLines(data []byte) []byte {
	decoded := p.tryDecodeBase64(data)

	lines := strings.Split(string(decoded), "\n")
	var cleanLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	return []byte(strings.Join(cleanLines, "\n"))
}

func (p *Parser) detectSourceType(source string) string {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return "url"
	}
	if strings.HasPrefix(source, "folder://") {
		return "folder"
	}
	if strings.HasPrefix(source, "file://") {
		return "file"
	}
	if strings.HasPrefix(source, "base64://") || strings.HasPrefix(source, "data:text/plain;base64,") {
		return "base64"
	}
	return "raw"
}

func (p *Parser) fetchURLContent(source string) (*fetchResult, error) {
	cleanURL, fragmentName := p.extractURLFragment(source)

	req, err := http.NewRequest("GET", cleanURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Xray-Checker")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Device-OS", "CheckerOS")
	req.Header.Set("X-Ver-OS", config.Version)
	req.Header.Set("X-Device-Model", "Xray-Checker Pro Max")
	req.Header.Set("X-Hwid", "0JLQq9Ca0JvQrtCn0Jgg0JHQm9Cp0KLQrCBIV0lE")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	name := fragmentName
	if name == "" {
		name = p.extractNameFromHeader(resp.Header.Get("profile-title"))
	}

	return &fetchResult{
		Content: content,
		Name:    name,
	}, nil
}

func (p *Parser) extractURLFragment(source string) (cleanURL string, name string) {
	if idx := strings.LastIndex(source, "#"); idx != -1 {
		name = strings.TrimSpace(source[idx+1:])
		cleanURL = source[:idx]
		if decoded, err := url.QueryUnescape(name); err == nil {
			name = decoded
		}
		return cleanURL, name
	}
	return source, ""
}

func (p *Parser) extractNameFromHeader(headerValue string) string {
	if headerValue == "" {
		return ""
	}

	headerValue = strings.TrimSpace(headerValue)

	if strings.HasPrefix(headerValue, "base64:") {
		encoded := strings.TrimPrefix(headerValue, "base64:")
		if decoded, err := p.decodeBase64(encoded); err == nil {
			return strings.TrimSpace(string(decoded))
		}
		return ""
	}

	if decoded, err := p.decodeBase64(headerValue); err == nil {
		decodedStr := string(decoded)
		if p.isPrintableString(decodedStr) {
			return strings.TrimSpace(decodedStr)
		}
	}

	return headerValue
}

func (p *Parser) isPrintableString(s string) bool {
	for _, r := range s {
		if r < 32 && r != '\t' && r != '\n' && r != '\r' {
			return false
		}
	}
	return true
}

func (p *Parser) parseOriginalLinks(rawData []byte) map[string]*originalLinkData {
	result := make(map[string]*originalLinkData)

	decoded := p.tryDecodeBase64(rawData)

	lines := strings.Split(string(decoded), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		data := p.parseShareLink(line)
		if data != nil {
			key := fmt.Sprintf("%s:%d", data.Server, data.Port)
			result[key] = &originalLinkData{
				Name:          data.Name,
				Encryption:    data.Encryption,
				Type:          data.Type,
				Path:          data.Path,
				Host:          data.Host,
				AllowInsecure: data.AllowInsecure,
			}
		}
	}

	return result
}

func (p *Parser) parseShareLink(link string) *parsedLink {
	if strings.HasPrefix(link, "vmess://") {
		return p.parseVMessLink(link)
	}

	u, err := url.Parse(link)
	if err != nil {
		return nil
	}

	result := &parsedLink{
		Name: u.Fragment,
	}

	host := u.Hostname()
	portStr := u.Port()
	if portStr == "" {
		return nil
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port == 0 {
		return nil
	}
	result.Server = host
	result.Port = port

	query := u.Query()
	result.Type = query.Get("type")
	result.Encryption = query.Get("encryption")
	result.Path = query.Get("path")
	result.Host = query.Get("host")
	result.AllowInsecure = query.Get("allowInsecure") == "1" || query.Get("allowInsecure") == "true"

	return result
}

func (p *Parser) parseVMessLink(link string) *parsedLink {
	encoded := strings.TrimPrefix(link, "vmess://")
	decoded, err := p.decodeBase64(encoded)
	if err != nil {
		return nil
	}

	var vmess map[string]interface{}
	if err := json.Unmarshal(decoded, &vmess); err != nil {
		return nil
	}

	result := &parsedLink{}

	if ps, ok := vmess["ps"].(string); ok {
		result.Name = ps
	}
	if add, ok := vmess["add"].(string); ok {
		result.Server = add
	}

	switch port := vmess["port"].(type) {
	case float64:
		result.Port = int(port)
	case string:
		if p, err := strconv.Atoi(port); err == nil {
			result.Port = p
		}
	}

	if result.Port == 0 {
		return nil
	}

	if net, ok := vmess["net"].(string); ok {
		result.Type = net
	}
	if host, ok := vmess["host"].(string); ok {
		result.Host = host
	}
	if path, ok := vmess["path"].(string); ok {
		result.Path = path
	}

	return result
}

func (p *Parser) convertOutbound(raw json.RawMessage, index int, originalData map[string]*originalLinkData) (*models.ProxyConfig, error) {
	var baseOutbound struct {
		Protocol       string                 `json:"protocol"`
		Tag            string                 `json:"tag"`
		SendThrough    string                 `json:"sendThrough"`
		Settings       json.RawMessage        `json:"settings"`
		StreamSettings *libXrayStreamSettings `json:"streamSettings"`
	}
	if err := json.Unmarshal(raw, &baseOutbound); err != nil {
		return nil, err
	}

	if baseOutbound.Protocol == "freedom" || baseOutbound.Protocol == "blackhole" || baseOutbound.Protocol == "dns" {
		return nil, nil
	}

	pc := &models.ProxyConfig{
		Index:    index,
		Name:     baseOutbound.SendThrough,
		Protocol: baseOutbound.Protocol,
	}

	if pc.Name == "" {
		pc.Name = baseOutbound.Tag
	}

	var flatSettings libXraySettings
	if err := json.Unmarshal(baseOutbound.Settings, &flatSettings); err == nil && flatSettings.Address != "" {
		pc.Server = flatSettings.Address
		pc.Port = flatSettings.Port

		switch baseOutbound.Protocol {
		case "vless":
			pc.UUID = flatSettings.ID
			pc.Flow = flatSettings.Flow
			pc.Encryption = flatSettings.Encryption
			pc.Level = flatSettings.Level
		case "vmess":
			pc.UUID = flatSettings.ID
			pc.AlterId = flatSettings.AlterId
			pc.Security = flatSettings.Security
			pc.Level = flatSettings.Level
		case "trojan":
			pc.Password = flatSettings.Password
		case "shadowsocks":
			pc.Password = flatSettings.Password
			pc.Method = flatSettings.Method
		}
	} else {
		var stdSettings xrayStandardSettings
		if err := json.Unmarshal(baseOutbound.Settings, &stdSettings); err != nil {
			return nil, fmt.Errorf("failed to parse settings: %v", err)
		}

		switch baseOutbound.Protocol {
		case "vless", "vmess":
			if len(stdSettings.Vnext) == 0 || len(stdSettings.Vnext[0].Users) == 0 {
				return nil, fmt.Errorf("no vnext/users found")
			}
			pc.Server = stdSettings.Vnext[0].Address
			pc.Port = stdSettings.Vnext[0].Port
			user := stdSettings.Vnext[0].Users[0]
			pc.UUID = user.ID
			pc.Flow = user.Flow
			pc.Encryption = user.Encryption
			pc.AlterId = user.AlterId
			pc.Level = user.Level
			if baseOutbound.Protocol == "vmess" {
				pc.Security = user.Security
			}
		case "trojan", "shadowsocks":
			if len(stdSettings.Servers) == 0 {
				return nil, fmt.Errorf("no servers found")
			}
			srv := stdSettings.Servers[0]
			pc.Server = srv.Address
			pc.Port = srv.Port
			pc.Password = srv.Password
			pc.Method = srv.Method
			pc.Flow = srv.Flow
		default:
			return nil, fmt.Errorf("unsupported protocol: %s", baseOutbound.Protocol)
		}
	}

	if pc.Server == "" || pc.Port == 0 {
		return nil, fmt.Errorf("failed to parse server/port")
	}

	if pc.Port == 0 || pc.Port == 1 {
		return nil, nil
	}

	if baseOutbound.StreamSettings != nil {
		ss := baseOutbound.StreamSettings
		pc.Type = ss.Network
		pc.Security = ss.Security

		if ss.TlsSettings != nil {
			pc.SNI = ss.TlsSettings.ServerName
			pc.AllowInsecure = ss.TlsSettings.AllowInsecure
			pc.Fingerprint = ss.TlsSettings.Fingerprint
			pc.ALPN = ss.TlsSettings.Alpn
		}

		if ss.RealitySettings != nil {
			pc.SNI = ss.RealitySettings.ServerName
			pc.Fingerprint = ss.RealitySettings.Fingerprint
			pc.PublicKey = ss.RealitySettings.PublicKey
			pc.ShortID = ss.RealitySettings.ShortId
		}

		if ss.Network == "raw" {
			pc.Type = "tcp"
		}

		if ss.RawSettings != nil && ss.RawSettings.Header != nil {
			pc.HeaderType = ss.RawSettings.Header.Type
			if ss.RawSettings.Header.Request != nil {
				if len(ss.RawSettings.Header.Request.Path) > 0 {
					pc.Path = ss.RawSettings.Header.Request.Path[0]
				}
				if ss.RawSettings.Header.Request.Headers != nil && len(ss.RawSettings.Header.Request.Headers.Host) > 0 {
					pc.Host = ss.RawSettings.Header.Request.Headers.Host[0]
				}
			}
		}

		if ss.WsSettings != nil {
			pc.Path = ss.WsSettings.Path
			if ss.WsSettings.Headers != nil {
				pc.Host = ss.WsSettings.Headers.Host
			}
			if pc.Host == "" {
				pc.Host = ss.WsSettings.Host
			}
		}

		if ss.GrpcSettings != nil {
			pc.ServiceName = ss.GrpcSettings.ServiceName
			pc.MultiMode = ss.GrpcSettings.MultiMode
		}

		if ss.HttpSettings != nil {
			pc.Path = ss.HttpSettings.Path
			if len(ss.HttpSettings.Host) > 0 {
				pc.Host = strings.Join(ss.HttpSettings.Host, ",")
			}
		}

		if ss.HttpupgradeSettings != nil {
			pc.Type = "httpupgrade"
			pc.Path = ss.HttpupgradeSettings.Path
			pc.Host = ss.HttpupgradeSettings.Host
		}

		if ss.Network == "xhttp" || ss.Network == "splithttp" {
			pc.Type = ss.Network

			var rawSettings json.RawMessage
			if len(ss.XhttpSettings) > 0 {
				rawSettings = ss.XhttpSettings
			} else if len(ss.SplithttpSettings) > 0 {
				rawSettings = ss.SplithttpSettings
			}

			if len(rawSettings) > 0 {
				pc.RawXhttpSettings = string(rawSettings)
				var parsed libXrayXhttpSettings
				if err := json.Unmarshal(rawSettings, &parsed); err == nil {
					pc.Path = parsed.Path
					pc.Host = parsed.Host
					pc.Mode = parsed.Mode
				}
			}
		}
	}

	key := fmt.Sprintf("%s:%d", pc.Server, pc.Port)
	if orig, ok := originalData[key]; ok {
		if pc.Name == "" && orig.Name != "" {
			pc.Name = orig.Name
		}
		if pc.Encryption == "" || pc.Encryption == "none" {
			if orig.Encryption != "" {
				pc.Encryption = orig.Encryption
			}
		}
		if orig.AllowInsecure {
			pc.AllowInsecure = true
		}
	}

	if pc.Name == "" {
		pc.Name = fmt.Sprintf("%s:%d", pc.Server, pc.Port)
	}

	if err := pc.Validate(); err != nil {
		return nil, err
	}

	pc.StableID = pc.GenerateStableID()

	return pc, nil
}

func (p *Parser) tryDecodeBase64(data []byte) []byte {
	text := strings.TrimSpace(string(data))

	if strings.HasPrefix(text, "vless://") || strings.HasPrefix(text, "vmess://") ||
		strings.HasPrefix(text, "trojan://") || strings.HasPrefix(text, "ss://") ||
		strings.HasPrefix(text, "{") || strings.HasPrefix(text, "[") {
		return data
	}

	decoded, err := p.decodeBase64(text)
	if err != nil {
		return data
	}

	return decoded
}

func (p *Parser) decodeBase64(text string) ([]byte, error) {
	text = strings.ReplaceAll(text, "-", "+")
	text = strings.ReplaceAll(text, "_", "/")

	if m := len(text) % 4; m != 0 {
		text += strings.Repeat("=", 4-m)
	}

	return base64.StdEncoding.DecodeString(text)
}

func (p *Parser) parseFolder(folderPath string) ([]*models.ProxyConfig, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read folder: %v", err)
	}

	var allConfigs []*models.ProxyConfig
	configIndex := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))
		if ext != ".json" {
			continue
		}

		filePath := filepath.Join(folderPath, fileName)
		data, err := os.ReadFile(filePath)
		if err != nil {
			logger.Warn("Failed to read file %s: %v", fileName, err)
			continue
		}

		configs, err := p.parseSingleConfigFile(data, configIndex)
		if err != nil {
			logger.Warn("Failed to parse file %s: %v", fileName, err)
			continue
		}

		for _, cfg := range configs {
			cfg.Index = configIndex
			allConfigs = append(allConfigs, cfg)
			configIndex++
		}

		logger.Debug("Parsed %d configs from %s", len(configs), fileName)
	}

	if len(allConfigs) == 0 {
		return nil, fmt.Errorf("no valid proxy configurations found in folder")
	}

	logger.Debug("Total configs from folder: %d", len(allConfigs))
	return allConfigs, nil
}

func (p *Parser) parseSingleConfigFile(data []byte, startIndex int) ([]*models.ProxyConfig, error) {
	trimmedData := strings.TrimSpace(string(data))

	if strings.HasPrefix(trimmedData, "[") {
		return p.parseJSONConfigs(data)
	}

	if strings.HasPrefix(trimmedData, "{") {
		var config struct {
			Remarks   string            `json:"remarks"`
			Outbounds []json.RawMessage `json:"outbounds"`
		}

		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %v", err)
		}

		var proxyConfigs []*models.ProxyConfig
		for _, outboundRaw := range config.Outbounds {
			proxyConfig, err := p.convertOutbound(outboundRaw, startIndex, nil)
			if err != nil {
				continue
			}
			if proxyConfig != nil {
				if config.Remarks != "" {
					proxyConfig.Name = config.Remarks
				}
				proxyConfigs = append(proxyConfigs, proxyConfig)
			}
		}

		if len(proxyConfigs) == 0 {
			return nil, fmt.Errorf("no valid proxy configurations found")
		}

		return proxyConfigs, nil
	}

	return nil, fmt.Errorf("unsupported config format")
}
