package orale

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func Load(applicationName string) (*Config, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	applicationNameRunes := []rune(applicationName)

	envPrefixRunes := []rune{}
	for i := 0; i < len(applicationNameRunes); i += 1 {
		currentChar := applicationNameRunes[i]
		var nextChar rune
		if i+1 < len(applicationNameRunes) {
			nextChar = applicationNameRunes[i+1]
		}
		if currentChar == '-' {
			continue
		}
		if unicode.IsLower(currentChar) {
			envPrefixRunes = append(envPrefixRunes, unicode.ToUpper(currentChar))
			if unicode.IsUpper(nextChar) {
				envPrefixRunes = append(envPrefixRunes, '_')
			}
		} else {
			envPrefixRunes = append(envPrefixRunes, currentChar)
		}
	}
	envPrefix := string(envPrefixRunes)

	configNameRunes := []rune{}
	for i := 0; i < len(applicationNameRunes); i += 1 {
		currentChar := applicationNameRunes[i]
		var nextChar rune
		if i+1 < len(applicationNameRunes) {
			nextChar = applicationNameRunes[i+1]
		}
		if currentChar == '_' {
			continue
		}
		if unicode.IsUpper(currentChar) {
			configNameRunes = append(configNameRunes, unicode.ToLower(currentChar))
		} else {
			configNameRunes = append(configNameRunes, currentChar)
			if unicode.IsLower(currentChar) && unicode.IsUpper(nextChar) {
				configNameRunes = append(configNameRunes, '-')
			}
		}
	}
	configName := string(configNameRunes)
	configName = fmt.Sprintf("%s.config.toml", configName)

	return LoadFromValues(
		os.Args[1:],
		envPrefix,
		os.Environ(),
		workingDir,
		[]string{configName},
	)
}

func LoadFromValues(programArgs []string, envVarPrefix string, envVars []string, configSearchStartPath string, configFileNames []string) (*Config, error) {
	flagValues := loadFlags(programArgs)
	environmentValues := loadEnvironment(envVarPrefix, envVars)
	configurationFiles, err := loadConfigurationFiles(configSearchStartPath, configFileNames)
	if err != nil {
		return nil, err
	}

	return &Config{
		FlagValues:         flagValues,
		EnvironmentValues:  environmentValues,
		ConfigurationFiles: configurationFiles,
	}, nil
}

// NOTE: programArgs should not include the program name - os.Args[1:]
// would be appropriate
func loadFlags(programArgs []string) map[string][]any {
	flagValues := map[string][]any{}
	// short flags
	for _, arg := range programArgs {
		isShortFlag := arg[0] == '-' && arg[1] != '-'
		isFlag := !isShortFlag && arg[0:2] == "--"

		var startIndex int
		switch {
		case isShortFlag:
			startIndex = 1
		case isFlag:
			startIndex = 2
		default:
			continue
		}

		splitIndex := -1
		for i := startIndex; i < len(arg); i += 1 {
			if arg[i] == '=' {
				splitIndex = i
				break
			}
		}
		if splitIndex == -1 {
			continue
		}

		key := arg[startIndex:splitIndex]
		value := arg[splitIndex+1:]

		key = strings.ToLower(key)
		key = strings.Replace(key, ".", "\\.", -1)
		key = strings.Replace(key, "--", ".", -1)
		key = strings.Replace(key, "-", "_", -1)

		if _, ok := flagValues[key]; !ok {
			flagValues[key] = []any{}
		}
		flagValues[key] = append(flagValues[key], value)
	}

	return flagValues
}

// NOTE: envVariables should be in the same format as the returned value from
// os.Environ()
func loadEnvironment(variablePrefix string, envVariables []string) map[string][]any {
	variablePrefix += "__"
	environmentValues := map[string][]any{}

	for _, envVariable := range envVariables {
		if len(envVariable) >= len(variablePrefix) && envVariable[0:len(variablePrefix)] == variablePrefix {
			splitIndex := -1
			for j := len(variablePrefix); j < len(envVariable); j += 1 {
				if envVariable[j] == '=' {
					splitIndex = j
					break
				}
			}
			if splitIndex == -1 {
				continue
			}

			key := envVariable[len(variablePrefix):splitIndex]
			value := envVariable[splitIndex+1:]

			key = strings.ToLower(key)
			key = strings.Replace(key, ".", "\\.", -1)
			key = strings.Replace(key, "__", ".", -1)

			if _, ok := environmentValues[key]; !ok {
				environmentValues[key] = []any{}
			}
			environmentValues[key] = append(environmentValues[key], value)
		}
	}

	return environmentValues
}

func loadConfigurationFiles(startPath string, configNames []string) ([]*File, error) {
	currentPathChunks := strings.Split(startPath, string(filepath.Separator))

	configFiles := []*File{}
	for i := len(currentPathChunks); i > 0; i -= 1 {
		isAbsPath := currentPathChunks[0] == ""
		currentPath := filepath.Join(currentPathChunks[:i]...)
		if isAbsPath {
			currentPath = string(filepath.Separator) + currentPath
		}

		for _, configName := range configNames {
			maybeConfigFilePath := filepath.Join(currentPath, configName)
			maybeConfigFile, err := maybeLoadFile(maybeConfigFilePath)
			if err != nil {
				return nil, err
			}
			if maybeConfigFile == nil {
				continue
			}

			configFiles = append(configFiles, maybeConfigFile)
		}
	}

	return configFiles, nil
}
