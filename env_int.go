package cfg

import (
	"strconv"
)

// Int get an integer value from a config
func (c *Env) Int(key string, def ...int) int {
	if v := c.Get(key); v != nil {
		return intValue(v)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func intValue(val any) int {
	var value int
	switch t := val.(type) {
	case float64:
		value = int(t)
	case bool:
		if t {
			value = 1
		} else {
			value = 0
		}
	case string:
		if i, errInt := strconv.ParseInt(t, 10, 64); errInt == nil {
			value = int(i)
		}
	}

	return value
}
