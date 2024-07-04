package cfg

import "strings"

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

func Escape(key string) string {
	if strings.IndexByte(key, '.') == -1 {
		return key
	}

	var (
		segmentSize int
		out         string
	)

	for {
		segmentSize = strings.IndexByte(key, '.')
		if segmentSize == -1 {
			segmentSize = len(key)
		} else if key[segmentSize-1] == '\\' {
			if out != "" {
				out += "\\."
			}
			out += key[:segmentSize-1]
			key = key[segmentSize+1:]
			continue
		}

		if out != "" {
			out += "\\."
		}

		out += key[:segmentSize]
		if segmentSize == len(key) {
			break
		}
		key = key[segmentSize+1:]
	}
	return out
}

// Segments extract all key segments.
// (Ex. "assets.alias.bootstrap\.css.filepath" = ["assets", "alias", "bootstrap.css", "filepath"])
func Segments(key string) (segments []string) {
	var (
		segmentSize int
		segment     string
	)

	for {
		segmentSize = strings.IndexByte(key, '.')
		if segmentSize == -1 {
			segmentSize = len(key)
		} else if key[segmentSize-1] == '\\' {
			segment += key[:segmentSize-1] + "."
			key = key[segmentSize+1:]
			continue
		}

		segments = append(segments, segment+key[:segmentSize])
		if segmentSize == len(key) {
			break
		}
		segment = ""
		key = key[segmentSize+1:]
	}
	return
}
