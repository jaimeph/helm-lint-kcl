package merger

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Merger struct {
	chartValuesYaml map[string]any
}

func New(chartValues []byte) (*Merger, error) {
	chartValuesYaml, err := unmarshal(chartValues)
	if err != nil {
		return nil, err
	}

	return &Merger{chartValuesYaml: chartValuesYaml}, nil
}

func (m *Merger) Merged() ([]byte, error) {
	merged, err := marshal(m.chartValuesYaml)
	if err != nil {
		return nil, err
	}
	return merged, nil
}

func unmarshal(content []byte) (map[string]any, error) {
	values := make(map[string]any)
	if err := yaml.Unmarshal(content, &values); err != nil {
		return nil, fmt.Errorf("failed to parse content: %v", err)
	}
	return values, nil
}

func marshal(content map[string]interface{}) ([]byte, error) {
	val, err := yaml.Marshal(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse content: %v", err)
	}
	return val, nil
}
