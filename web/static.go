package web

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed static/*
var staticFiles embed.FS

func StaticHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filePath := strings.TrimPrefix(r.URL.Path, "/static/")
		if filePath == "" {
			http.NotFound(w, r)
			return
		}

		var data []byte
		var found bool

		loader := GetAssetLoader()
		if loader != nil && loader.IsEnabled() {
			data, found = loader.GetFile(filePath)

			if !found && (filePath == "logo-dark.svg" || filePath == "logo-light.svg") {
				variant := "logo-dark"
				if filePath == "logo-light.svg" {
					variant = "logo-light"
				}

				fallbacks := []struct {
					name        string
					contentType string
				}{
					{variant + ".png", "image/png"},
					{"logo.svg", "image/svg+xml"},
					{"logo.png", "image/png"},
				}

				for _, fb := range fallbacks {
					if logoData, ok := loader.GetFile(fb.name); ok {
						w.Header().Set("Content-Type", fb.contentType)
						w.Header().Set("Cache-Control", "public, max-age=31536000")
						w.Write(logoData)
						return
					}
				}
			}
		}

		if !found {
			embeddedData, err := fs.ReadFile(staticFiles, path.Join("static", filePath))
			if err != nil {
				if loader != nil && loader.IsEnabled() {
					if customData, customFound := loader.GetFile(filePath); customFound {
						data = customData
						found = true
					}
				}
				if !found {
					http.NotFound(w, r)
					return
				}
			} else {
				data = embeddedData
			}
		}

		contentType := GetContentType(filePath)
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Write(data)
	}
}
