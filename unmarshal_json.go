package cfg

import "encoding/json"

func JsonUnmarshal(content []byte) (map[string]any, error) {
	var config map[string]any
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, err
	}
	return config, nil
}
