package cfg

// represents a object, used to create default settings
type O map[string]any

// Merge merges the current object with the other without maintaining a reference
func (o O) Merge(other O) {
	if other == nil {
		return
	}
	for key, value := range other {
		switch v := value.(type) {
		case O:
			vl := O{}
			vl.Merge(v)
			o[key] = vl
		case []string:
			o[key] = v[0:]
		case []O:
			var vl []O
			for _, it := range v {
				vli := O{}
				vli.Merge(it)
				vl = append(vl, vli)
			}
			o[key] = vl
		default:
			o[key] = value
		}
	}
}

type DefaultConfigFn func(extra ...O) O

// CreateDefaultConfigFn simplifies the creation of default configurations
func CreateDefaultConfigFn(defaultValue O) DefaultConfigFn {
	return func(extra ...O) O {
		config := O{}
		config.Merge(defaultValue)

		for _, o2 := range extra {
			config.Merge(o2)
		}
		return config
	}
}
