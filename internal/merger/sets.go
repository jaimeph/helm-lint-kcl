package merger

import (
	"fmt"
	"strconv"
	"strings"
)

func (m *Merger) Sets(sets []string) error {
	for _, set := range sets {
		keys, value, err := resolveSet(set)
		if err != nil {
			return err
		}
		addSetRecursive(m.chartValuesYaml, keys, value)
	}
	return nil
}

func addSetRecursive(content map[string]any, keys []string, value any) {
	if len(keys) == 0 {
		return
	}
	if len(keys) == 1 {
		content[keys[0]] = value
		return
	}
	if _, exists := content[keys[0]]; !exists {
		content[keys[0]] = make(map[string]any)
	}
	addSetRecursive(content[keys[0]].(map[string]any), keys[1:], value)
}

func resolveSet(set string) ([]string, any, error) {
	key, value, err := splitKeyValue(set)
	if err != nil {
		return nil, "", err
	}
	key = strings.ReplaceAll(key, "\\", "")
	keys := splitKey(key)
	return keys, value, nil
}

func splitKeyValue(value string) (string, any, error) {
	if !strings.Contains(value, "=") {
		return "", "", fmt.Errorf("invalid value format set '%s'", value)
	}
	parts := strings.Split(value, "=")
	return parts[0], isNumberReturn(parts[1]), nil
}

func splitKey(value string) []string {
	if !strings.Contains(value, ".") {
		return []string{value}
	}
	keys := strings.Split(value, ".")
	for i, k := range keys {
		if strings.EqualFold(k, "annotations") || strings.EqualFold(k, "labels") {
			return append(keys[:i], k, strings.Join(keys[i+1:], "."))
		}
	}
	return keys
}

func isNumberReturn(value string) any {
	number, err := strconv.Atoi(value)
	if err != nil {
		return value
	}
	return number
}
