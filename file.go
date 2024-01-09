package orale

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type File struct {
	Path   string
	Values map[string][]any
}

func maybeLoadFile(maybeConfigFilePath string) (*File, error) {
	fileBytes, err := os.ReadFile(maybeConfigFilePath)
	if err != nil {
		switch {
		case os.IsNotExist(err), os.IsPermission(err):
			return nil, nil
		default:
			return nil, err
		}
	}

	fileStr := string(fileBytes)

	var hierarchicalFileValues map[string]any
	if _, err := toml.Decode(fileStr, &hierarchicalFileValues); err != nil {
		return nil, err
	}
	fileValues := map[string][]any{}
	flattenFileValues(nil, hierarchicalFileValues, fileValues)

	return &File{
		Path:   maybeConfigFilePath,
		Values: fileValues,
	}, nil
}

func flattenFileValues(pathChunks []string, hierarchicalValues map[string]any, flattenedValues map[string][]any) {
	if pathChunks == nil {
		pathChunks = []string{}
	}

	// TODO: deal with slice indexes
	for key, value := range hierarchicalValues {
		keyPathChunks := append(pathChunks, key)
		keyPath := strings.Join(keyPathChunks, ".")

		switch val := value.(type) {
		case []map[string]any:
			for i, v := range val {
				subKeyPathChunks := []string{}
				for j, chunk := range keyPathChunks {
					if j == len(keyPathChunks)-1 {
						chunk = fmt.Sprintf("%s[%d]", chunk, i)
					}
					subKeyPathChunks = append(subKeyPathChunks, chunk)
				}
				flattenFileValues(subKeyPathChunks, v, flattenedValues)
			}
		case []any:
			for i, v := range val {
				subKeyPathChunks := []string{}
				for j, chunk := range keyPathChunks {
					if j == len(keyPathChunks)-1 {
						chunk = fmt.Sprintf("%s[%d]", chunk, i)
					}
					subKeyPathChunks = append(subKeyPathChunks, chunk)
				}
				keyPath = strings.Join(subKeyPathChunks, ".")
				if _, ok := flattenedValues[keyPath]; !ok {
					flattenedValues[keyPath] = []any{}
				}
				flattenedValues[keyPath] = append(flattenedValues[keyPath], v)
			}
		case map[string]any:
			flattenFileValues(keyPathChunks, val, flattenedValues)
		default:
			if _, ok := flattenedValues[keyPath]; !ok {
				flattenedValues[keyPath] = []any{}
			}
			flattenedValues[keyPath] = append(flattenedValues[keyPath], value)
		}
	}
}
