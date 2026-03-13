package xray

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"xray-checker/logger"

	"github.com/xtls/xray-core/common/log"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf/serial"
	_ "github.com/xtls/xray-core/main/distro/all"
)

type filteredLogHandler struct{}

func (h *filteredLogHandler) Handle(msg log.Message) {
	msgStr := msg.String()
	if strings.Contains(msgStr, "deprecated") {
		return
	}
	if strings.HasPrefix(msgStr, "[Warning]") {
		return
	}
	logger.Debug("xray: %s", msgStr)
}

func init() {
	log.RegisterHandler(&filteredLogHandler{})
}

type Runner struct {
	instance   *core.Instance
	configFile string
}

func NewRunner(configFile string) *Runner {
	return &Runner{
		configFile: configFile,
	}
}

func (r *Runner) Start() error {
	configBytes, err := os.ReadFile(r.configFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	xrayConfig, err := serial.DecodeJSONConfig(bytes.NewReader(configBytes))
	if err != nil {
		return fmt.Errorf("error decoding config: %v", err)
	}

	coreConfig, err := xrayConfig.Build()
	if err != nil {
		return fmt.Errorf("error building config: %v", err)
	}

	instance, err := core.New(coreConfig)
	if err != nil {
		return fmt.Errorf("error creating Xray instance: %v", err)
	}

	if err := instance.Start(); err != nil {
		return fmt.Errorf("error starting Xray: %v", err)
	}

	r.instance = instance
	logger.Debug("Xray instance started")

	return nil
}

func (r *Runner) Stop() error {
	if r.instance != nil {
		err := r.instance.Close()
		r.instance = nil
		if err != nil {
			return fmt.Errorf("error stopping Xray: %v", err)
		}
		logger.Debug("Xray instance stopped")
	}
	return nil
}
