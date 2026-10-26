package config

import "testing"

func TestConfigFiles(t *testing.T) {
	for _, file := range []string{"config_local.yml", "config_test.yml", "config.yml"} {
		t.Run(file, func(t *testing.T) {
			if err := Init(file); err != nil {
				t.Fatalf("Init(%q) error = %v", file, err)
			}
		})
	}
}
