package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 加载配置
	config, err := LoadConfig(".")
	if err != nil {
		t.Fatalf("Loading config failed: %v", err)
	}

	t.Logf("config: %+v", config)

}
