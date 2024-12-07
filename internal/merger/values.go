package merger

import (
	"os"
)

func (m *Merger) Values(values []string) error {
	for _, value := range values {
		content, err := os.ReadFile(value)
		if err != nil {
			return err
		}
		contentYaml, err := unmarshal(content)
		if err != nil {
			return err
		}
		mergeMaps(m.chartValuesYaml, contentYaml)
	}
	return nil
}

func mergeMaps(map1, map2 map[string]any) map[string]any {
	for k, v2 := range map2 {
		if v1, ok := map1[k]; ok {
			// Si ambos son mapas, fusiónalos recursivamente
			if mapV1, ok1 := v1.(map[string]any); ok1 {
				if mapV2, ok2 := v2.(map[string]any); ok2 {
					map1[k] = mergeMaps(mapV1, mapV2)
					continue
				}
			}
		}
		// Si no, sobrescribe el valor en map1
		map1[k] = v2
	}
	return map1
}
