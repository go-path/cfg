package cfg

import (
	"strconv"
)

// Float obtém um valor inteiro de uma config
func (c *Env) Float(key string, def ...float64) float64 {
	if v := c.Get(key); v != nil {
		return floatValue(v)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func floatValue(val interface{}) float64 {
	var value float64
	switch t := val.(type) {
	case float64:
		value = t
	case bool:
		if t {
			value = 1
		} else {
			value = 0
		}
	case string:
		if i, errInt := strconv.ParseFloat(t, 64); errInt == nil {
			value = i
		}
	}

	return value
}
