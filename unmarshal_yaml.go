package cfg

import "gopkg.in/yaml.v3"

func YamlUnmarshal(content []byte) (map[string]any, error) {
	var config map[string]any
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, err
	}
	return config, nil
}
