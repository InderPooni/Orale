package orale

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

func (l *Loader) Get(path string, target any) error {
	targetRefVal := reflect.ValueOf(target)
	if targetRefVal.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}
	targetRefVal = targetRefVal.Elem()

	return getFromLoader(l, path, targetRefVal, 0)
}

func (l *Loader) MustGet(path string, target any) {
	err := l.Get(path, target)
	if err != nil {
		panic(err)
	}
}

func getFromLoader(l *Loader, currentPath string, targetRefVal reflect.Value, index int) error {
	switch targetRefVal.Kind() {
	case reflect.Ptr:
		if targetRefVal.IsNil() {
			targetRefVal.Set(reflect.New(targetRefVal.Type().Elem()))
		}
		return getFromLoader(l, currentPath, targetRefVal.Elem(), 0)

	case reflect.Struct:
		for i := 0; i < targetRefVal.NumField(); i += 1 {
			field := targetRefVal.Type().Field(i)
			fieldTag := field.Tag.Get("config")
			if fieldTag == "" {
				fieldTag = calDefaultFieldTag(field.Name)
			}
			if currentPath != "" {
				fieldTag = currentPath + "." + fieldTag
			}
			if err := getFromLoader(l, fieldTag, targetRefVal.Field(i), 0); err != nil {
				return err
			}
		}

	case reflect.Slice:
		if targetRefVal.IsNil() {
			targetRefVal.Set(reflect.MakeSlice(targetRefVal.Type(), 0, 0))
		}
		valueLen, err := resolvePathLen(l, currentPath)
		if err != nil {
			return err
		}
		if valueLen > 0 {
			targetRefVal.Set(reflect.MakeSlice(targetRefVal.Type(), valueLen, valueLen))
			for i := 0; i < valueLen; i += 1 {
				if err := getFromLoader(l, fmt.Sprintf("%s[%d]", currentPath, i), targetRefVal.Index(i), 0); err != nil {
					return err
				}
			}
		} else {
			value, err := resolveValue(l, currentPath)
			if err != nil {
				return err
			}
			if value != nil {
				targetRefVal.Set(reflect.MakeSlice(targetRefVal.Type(), len(value), len(value)))
				for i := 0; i < len(value); i += 1 {
					if err := getFromLoader(l, currentPath, targetRefVal.Index(i), i); err != nil {
						return err
					}
				}
			} else {
				targetRefVal.Set(reflect.MakeSlice(targetRefVal.Type(), 0, 0))
			}
		}

	case reflect.String:
		value, err := resolveValue(l, currentPath)
		if err != nil {
			return err
		}
		if len(value) > index {
			strValue, ok := value[index].(string)
			if ok {
				targetRefVal.SetString(strValue)
			} else {
				targetRefVal.SetString("")
			}
		} else {
			targetRefVal.SetString("")
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := resolveValue(l, currentPath)
		if err != nil {
			return err
		}
		if len(value) > index {
			int64Value, ok := value[index].(int64)
			if ok {
				targetRefVal.SetInt(int64Value)
			} else {
				targetRefVal.SetInt(0)
			}
		} else {
			targetRefVal.SetInt(0)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := resolveValue(l, currentPath)
		if err != nil {
			return err
		}
		if len(value) > index {
			uint64Value, ok := value[index].(uint64)
			if ok && len(value) > 0 {
				targetRefVal.SetUint(uint64Value)
			} else {
				targetRefVal.SetUint(0)
			}
		} else {
			targetRefVal.SetUint(0)
		}

	case reflect.Float32, reflect.Float64:
		value, err := resolveValue(l, currentPath)
		if err != nil {
			return err
		}
		if len(value) > index {
			float64Value, ok := value[index].(float64)
			if ok {
				targetRefVal.SetFloat(float64Value)
			} else {
				targetRefVal.SetFloat(0)
			}
		} else {
			targetRefVal.SetFloat(0)
		}

	case reflect.Bool:
		value, err := resolveValue(l, currentPath)
		if err != nil {
			return err
		}
		if len(value) > index {
			if len(value) > 0 {
				targetRefVal.SetBool(value[index].(bool))
			} else {
				targetRefVal.SetBool(false)
			}
		} else {
			targetRefVal.SetBool(false)
		}

	default:
		return fmt.Errorf("unsupported type %s", targetRefVal.Kind())
	}
	return nil
}

func resolveValue(l *Loader, targetPath string) ([]any, error) {
	if targetPath == "" {
		return nil, fmt.Errorf("target path cannot be empty")
	}
	if value, ok := l.FlagValues[targetPath]; ok {
		return value, nil
	} else if value, ok := l.EnvironmentValues[targetPath]; ok {
		return value, nil
	} else {
		for _, file := range l.ConfigurationFiles {
			if value, ok := file.Values[targetPath]; ok {
				return value, nil
			}
		}
	}
	return nil, nil
}

func resolvePathLen(l *Loader, targetPath string) (int, error) {
	if targetPath == "" {
		return 0, fmt.Errorf("target path cannot be empty")
	}

	flagPaths := map[string]bool{}
	for flagPath := range l.FlagValues {
		slicePath := getSlicePathFromSubjectAndTargetPaths(flagPath, targetPath)
		if slicePath != "" {
			flagPaths[slicePath] = true
		}
	}
	if len(flagPaths) != 0 {
		return len(flagPaths), nil
	}

	environmentPaths := map[string]bool{}
	for environmentPath := range l.EnvironmentValues {
		slicePath := getSlicePathFromSubjectAndTargetPaths(environmentPath, targetPath)
		if slicePath != "" {
			environmentPaths[slicePath] = true
		}
	}
	if len(environmentPaths) != 0 {
		return len(environmentPaths), nil
	}

	for _, file := range l.ConfigurationFiles {
		filePaths := map[string]bool{}
		for filePath := range file.Values {
			slicePath := getSlicePathFromSubjectAndTargetPaths(filePath, targetPath)
			if slicePath != "" {
				filePaths[slicePath] = true
			}
		}
		if len(filePaths) != 0 {
			return len(filePaths), nil
		}
	}

	return 0, nil
}

func getSlicePathFromSubjectAndTargetPaths(subjectPath, targetPath string) string {
	if len(subjectPath) < len(targetPath)+3 {
		return ""
	}
	remainingPath := subjectPath[len(targetPath):]
	if remainingPath[0] != '[' {
		return ""
	}
	endIndexOffset := 0
	for i, r := range remainingPath {
		if r == ']' {
			endIndexOffset = i
			break
		}
	}
	return subjectPath[:len(targetPath)+endIndexOffset+1]
}

func calDefaultFieldTag(fieldName string) string {
	fieldTag := ""
	for i, r := range fieldName {
		if unicode.IsUpper(r) {
			if i != 0 {
				fieldTag += "_"
			}
			fieldTag += strings.ToLower(string(r))
		} else {
			fieldTag += string(r)
		}
	}
	return fieldTag
}
