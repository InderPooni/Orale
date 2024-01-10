package orale_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/RobertWHurst/orale"
)

var testAssetsPath string

func init() {
	_, currentFilePath, _, _ := runtime.Caller(0)
	testAssetsPath = filepath.Join(filepath.Dir(currentFilePath), "test-assets")
}

func TestLoadAndGetConfig(t *testing.T) {
	t.Parallel()

	t.Run("should load flags, environment variables, and configuration files, then correctly assign them to a struct", func(t *testing.T) {
		type TestConfig struct {
			Database struct {
				ConnectionUri string `config:"connection_uri"`
			} `config:"database"`
			Server struct {
				Port int `config:"port"`
			} `config:"server"`
			Channels []struct {
				Name string `config:"name"`
				Id   int    `config:"id"`
			} `config:"channels"`
		}

		programArgs := []string{
			"--database.connection_uri=postgres://localhost:5432",
		}
		envVars := []string{
			"TEST__SERVER__PORT=8080",
		}
		configSearchStartPath := testAssetsPath
		configFileNames := []string{"test-config-3.toml"}

		loader, err := orale.LoadFromValues(programArgs, "TEST", envVars, configSearchStartPath, configFileNames)
		if err != nil {
			t.Fatal(err)
		}

		var testConfig TestConfig
		if err := loader.Get("", &testConfig); err != nil {
			t.Fatal(err)
		}
	})
}
