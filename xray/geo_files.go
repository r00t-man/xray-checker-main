package xray

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"xray-checker/logger"
)

const (
	geoSiteURL  = "https://github.com/v2fly/domain-list-community/releases/latest/download/dlc.dat"
	geoIPURL    = "https://github.com/v2fly/geoip/releases/latest/download/geoip.dat"
	geoSiteFile = "geo/geosite.dat"
	geoIPFile   = "geo/geoip.dat"
)

type GeoFileManager struct {
	baseDir string
}

func NewGeoFileManager(baseDir string) *GeoFileManager {
	if baseDir == "" {
		if wd, err := os.Getwd(); err == nil {
			baseDir = wd
		} else {
			baseDir = "."
		}
	}

	return &GeoFileManager{
		baseDir: baseDir,
	}
}

func (gfm *GeoFileManager) EnsureGeoFiles() error {
	if err := gfm.ensureFile(geoSiteFile, geoSiteURL); err != nil {
		return fmt.Errorf("failed to ensure geosite.dat: %v", err)
	}

	if err := gfm.ensureFile(geoIPFile, geoIPURL); err != nil {
		return fmt.Errorf("failed to ensure geoip.dat: %v", err)
	}

	return nil
}

func (gfm *GeoFileManager) ensureFile(filename, url string) error {
	filePath := filepath.Join(gfm.baseDir, filename)

	if _, err := os.Stat(filePath); err == nil {
		return nil
	}

	logger.Info("Downloading %s...", filename)

	fileDir := filepath.Dir(filePath)
	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	if err := gfm.downloadFile(url, filePath); err != nil {
		return fmt.Errorf("failed to download %s: %v", filename, err)
	}

	logger.Info("Downloaded %s", filename)
	return nil
}

func (gfm *GeoFileManager) downloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
