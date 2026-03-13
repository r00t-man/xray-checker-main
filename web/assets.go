package web

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
	"xray-checker/logger"
)

type AssetLoader struct {
	basePath       string
	files          map[string][]byte
	customTemplate *template.Template
	hasCustomCSS   bool
	enabled        bool
}

var globalAssetLoader *AssetLoader

func InitAssetLoader(customPath string) error {
	loader := &AssetLoader{
		basePath: customPath,
		files:    make(map[string][]byte),
		enabled:  customPath != "",
	}

	if loader.enabled {
		if err := loader.load(); err != nil {
			return err
		}
	}

	globalAssetLoader = loader
	return nil
}

func GetAssetLoader() *AssetLoader {
	return globalAssetLoader
}

func (a *AssetLoader) load() error {
	info, err := os.Stat(a.basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("custom assets directory does not exist: %s", a.basePath)
		}
		return fmt.Errorf("error accessing custom assets directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("custom assets path is not a directory: %s", a.basePath)
	}

	logger.Info("Custom assets enabled: %s", a.basePath)

	entries, err := os.ReadDir(a.basePath)
	if err != nil {
		return fmt.Errorf("error reading custom assets directory: %w", err)
	}

	var loadedFiles []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		filePath := filepath.Join(a.basePath, name)

		data, err := os.ReadFile(filePath)
		if err != nil {
			logger.Warn("Failed to read custom asset %s: %v", name, err)
			continue
		}

		a.files[name] = data
		loadedFiles = append(loadedFiles, name)

		if name == "custom.css" {
			a.hasCustomCSS = true
		}
	}

	if len(loadedFiles) > 0 {
		logger.Info("Custom assets loaded:")
		for _, name := range loadedFiles {
			logger.Info("  ✓ %s", name)
		}
	} else {
		logger.Info("Custom assets loaded: (none)")
	}

	if data, exists := a.files["index.html"]; exists {
		funcMap := template.FuncMap{
			"formatLatency": func(d time.Duration) string {
				if d == 0 {
					return "n/a"
				}
				return fmt.Sprintf("%dms", d.Milliseconds())
			},
		}

		tmpl, err := template.New("index.html").Funcs(funcMap).Parse(string(data))
		if err != nil {
			return fmt.Errorf("failed to parse custom template index.html: %w", err)
		}

		a.customTemplate = tmpl
		logger.Info("Using custom template: index.html")
	} else {
		logger.Info("Using default template")
	}

	return nil
}

func (a *AssetLoader) IsEnabled() bool {
	return a.enabled
}

func (a *AssetLoader) HasCustomTemplate() bool {
	return a.customTemplate != nil
}

func (a *AssetLoader) GetCustomTemplate() *template.Template {
	return a.customTemplate
}

func (a *AssetLoader) HasCustomCSS() bool {
	return a.hasCustomCSS
}

func (a *AssetLoader) GetFile(name string) ([]byte, bool) {
	data, exists := a.files[name]
	return data, exists
}

func GetContentType(filename string) string {
	switch {
	case strings.HasSuffix(filename, ".css"):
		return "text/css; charset=utf-8"
	case strings.HasSuffix(filename, ".js"):
		return "application/javascript; charset=utf-8"
	case strings.HasSuffix(filename, ".html"):
		return "text/html; charset=utf-8"
	case strings.HasSuffix(filename, ".json"):
		return "application/json; charset=utf-8"
	case strings.HasSuffix(filename, ".woff2"):
		return "font/woff2"
	case strings.HasSuffix(filename, ".woff"):
		return "font/woff"
	case strings.HasSuffix(filename, ".ttf"):
		return "font/ttf"
	case strings.HasSuffix(filename, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(filename, ".png"):
		return "image/png"
	case strings.HasSuffix(filename, ".jpg"), strings.HasSuffix(filename, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(filename, ".gif"):
		return "image/gif"
	case strings.HasSuffix(filename, ".ico"):
		return "image/x-icon"
	case strings.HasSuffix(filename, ".webp"):
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
