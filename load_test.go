package orale_test

import (
	"os"
	"path/filepath"
	"testing"

	orale "github.com/RobertWHurst/Orale"
)

func TestLoad(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	defer os.Chdir(cwd)

	os.Chdir(testAssetsPath)

	t.Run("should load flags, environment variables, and configuration files, then correctly assign them to a struct", func(t *testing.T) {
		type TestConfig struct {
			A    string `config:"a"`
			Easy *struct {
				As []int `config:"as"`
			} `config:"easy"`
			Abc []struct {
				Baby string `config:"baby"`
				And  string `config:"and"`
			} `config:"abc"`
		}

		conf, err := orale.Load("test-application")
		if err != nil {
			t.Fatal(err)
		}

		testConf := TestConfig{}
		conf.MustGet("", &testConf)

		if testConf.A != "bc" {
			t.Fatalf("expected A to be bc, got %s", testConf.A)
		}
		if testConf.Easy == nil {
			t.Fatal("expected Easy to not be nil")
		}
		if len(testConf.Easy.As) != 3 {
			t.Fatalf("expected Easy.As to have 3 values, got %d", len(testConf.Easy.As))
		}
		if testConf.Easy.As[0] != 1 {
			t.Fatalf("expected Easy.As[0] to be 1, got %d", testConf.Easy.As[0])
		}
		if testConf.Easy.As[1] != 2 {
			t.Fatalf("expected Easy.As[1] to be 2, got %d", testConf.Easy.As[1])
		}
		if testConf.Easy.As[2] != 3 {
			t.Fatalf("expected Easy.As[2] to be 3, got %d", testConf.Easy.As[2])
		}
		if len(testConf.Abc) != 2 {
			t.Fatalf("expected Abc to have 2 values, got %d", len(testConf.Abc))
		}
		if testConf.Abc[0].Baby != "you" {
			t.Fatalf("expected Abc[0].Baby to be you, got %s", testConf.Abc[0].Baby)
		}
		if testConf.Abc[0].And != "me" {
			t.Fatalf("expected Abc[0].And to be me, got %s", testConf.Abc[0].And)
		}
		if testConf.Abc[1].Baby != "you" {
			t.Fatalf("expected Abc[1].Baby to be you, got %s", testConf.Abc[1].Baby)
		}
		if testConf.Abc[1].And != "me" {
			t.Fatalf("expected Abc[1].And to be me, got %s", testConf.Abc[1].And)
		}
	})
}

func TestLoadFromValues(t *testing.T) {
	t.Parallel()

	t.Run("should load flag and short flag values", func(t *testing.T) {
		t.Parallel()

		programArgs := []string{
			"--flag1=value1",
			"--flag2=value2",
			"-f=value3",
			"-g=value4",
		}

		conf, err := orale.LoadFromValues(programArgs, "", []string{}, "", []string{})
		if err != nil {
			t.Fatal(err)
		}

		if len(conf.FlagValues) != 4 {
			t.Fatalf("expected 4 flag values, got %d", len(conf.FlagValues))
		}

		if conf.FlagValues["flag1"][0] != "value1" {
			t.Fatalf("expected flag1 to be value1, got %s", conf.FlagValues["flag1"])
		}

		if conf.FlagValues["flag2"][0] != "value2" {
			t.Fatalf("expected flag2 to be value2, got %s", conf.FlagValues["flag2"])
		}

		if conf.FlagValues["f"][0] != "value3" {
			t.Fatalf("expected f to be value3, got %s", conf.FlagValues["f"])
		}

		if conf.FlagValues["g"][0] != "value4" {
			t.Fatalf("expected g to be value4, got %s", conf.FlagValues["g"])
		}
	})

	t.Run("should load environment values", func(t *testing.T) {
		t.Parallel()

		envVars := []string{
			"TEST__ENV1=value1",
			"TEST__ENV2=value2",
		}

		conf, err := orale.LoadFromValues([]string{}, "TEST", envVars, "", []string{})
		if err != nil {
			t.Fatal(err)
		}

		if len(conf.EnvironmentValues) != 2 {
			t.Fatalf("expected 2 environment values, got %d", len(conf.EnvironmentValues))
		}

		if conf.EnvironmentValues["env1"][0] != "value1" {
			t.Fatalf("expected env1 to be value1, got %s", conf.EnvironmentValues["env1"])
		}

		if conf.EnvironmentValues["env2"][0] != "value2" {
			t.Fatalf("expected env2 to be value2, got %s", conf.EnvironmentValues["env2"])
		}
	})

	t.Run("should load configuration files", func(t *testing.T) {
		t.Parallel()

		configSearchStartPath := filepath.Join(testAssetsPath, "search-dir")
		configFileNames := []string{
			"test-config-1.toml",
			"test-config-2.toml",
		}

		conf, err := orale.LoadFromValues([]string{}, "", []string{}, configSearchStartPath, configFileNames)
		if err != nil {
			t.Fatal(err)
		}

		if len(conf.ConfigurationFiles) != 2 {
			t.Fatalf("expected 2 configuration files, got %d", len(conf.ConfigurationFiles))
		}

		config1Path := filepath.Join(testAssetsPath, "search-dir/test-config-1.toml")
		config2Path := filepath.Join(testAssetsPath, "test-config-2.toml")

		if len(conf.ConfigurationFiles) != 2 {
			t.Fatalf("expected 2 values in test_config, got %d", len(conf.ConfigurationFiles))
		}

		if conf.ConfigurationFiles[0].Path != config1Path {
			t.Fatalf("expected first configuration file path to be %x, got %x", config1Path, conf.ConfigurationFiles[0].Path)
		}
		if conf.ConfigurationFiles[0].Values["test_val_1"][0] != int64(1) {
			t.Fatalf("expected test_val_1 to be 3, got %s", conf.ConfigurationFiles[0].Values["test_val_1"])
		}
		if conf.ConfigurationFiles[0].Values["test_val_2"][0] != int64(2) {
			t.Fatalf("expected test_val_3 to be 4, got %s", conf.ConfigurationFiles[0].Values["test_val_2"])
		}

		if conf.ConfigurationFiles[1].Path != config2Path {
			t.Fatalf("expected second configuration file path to be %s, got %s", config2Path, conf.ConfigurationFiles[1].Path)
		}
		if conf.ConfigurationFiles[1].Values["test_val_1"][0] != int64(3) {
			t.Fatalf("expected test_val_1 to be 1, got %s", conf.ConfigurationFiles[1].Values["test_val_1"])
		}
		if conf.ConfigurationFiles[1].Values["test_val_3"][0] != int64(4) {
			t.Fatalf("expected test_val_2 to be 2, got %s", conf.ConfigurationFiles[1].Values["test_val_3"])
		}
	})

	t.Run("should handle multi entry values", func(t *testing.T) {
		t.Parallel()

		programArgs := []string{
			"--flag=value1",
			"--flag=value2",
			"-f=value3",
			"-f=value4",
		}
		envVars := []string{
			"TEST__ENV=value1",
			"TEST__ENV=value2",
		}

		conf, err := orale.LoadFromValues(programArgs, "TEST", envVars, "", []string{})
		if err != nil {
			t.Fatal(err)
		}

		if len(conf.FlagValues) != 2 {
			t.Fatalf("expected 2 flag values, got %d", len(conf.FlagValues))
		}

		if conf.FlagValues["flag"][0] != "value1" && conf.FlagValues["flag"][1] != "value2" {
			t.Fatalf("expected flag to be value1 and value2, got %v", conf.FlagValues["flag"])
		}

		if conf.FlagValues["f"][0] != "value3" && conf.FlagValues["f"][1] != "value4" {
			t.Fatalf("expected f to be value3 and value4, got %v", conf.FlagValues["f"])
		}

		if len(conf.EnvironmentValues) != 1 {
			t.Fatalf("expected 1 environment value, got %d", len(conf.EnvironmentValues))
		}

		if conf.EnvironmentValues["env"][0] != "value1" && conf.EnvironmentValues["env"][1] != "value2" {
			t.Fatalf("expected env to be value1 and value2, got %v", conf.EnvironmentValues["env"])
		}
	})
}
