package cfg

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func (c *Env) get(key string) (any, bool) {
	unlock := c.lock(true)

	if e, exist := c.cache[key]; exist { // 1st check
		unlock()
		return e.value, e.exist
	} else {
		unlock()
	}

	unlock = c.lock(false)
	defer unlock()

	if e, exist := c.cache[key]; exist { // 2nd (double) check
		return e.value, e.exist
	}

	return c.getValueUnsafe(key)
}

func (c *Env) set(key string, value any) {

	indexOpenBracket := strings.IndexByte(key, '[')

	if indexOpenBracket == 0 {
		slog.Warn("[cfg] invalid config key.", slog.String("key", key))
		return
	}

	object := map[string]any{}
	var segment string
	var prev map[string]any
	var curr = object

	if indexOpenBracket == -1 {
		// simple key (fast)

		if strings.IndexByte(key, ']') >= 0 {
			// "prop.array0]"
			slog.Warn("[cfg] invalid brackets in config key.", slog.String("key", key))
			return
		}

		for _, segment = range Segments(key) {
			segment = strings.TrimSpace(segment)
			if segment == "" {
				// "prop..value" || "prop. .value"
				slog.Warn("[cfg] invalid config key.", slog.String("key", key))
				return
			}

			prev = curr
			curr[segment] = map[string]any{}
			curr = curr[segment].(map[string]any)
		}

	} else {
		// complex key (array)
		if strings.IndexByte(key, ']') == -1 {
			// "prop.array[0"
			slog.Warn("[cfg] invalid brackets in config key.", slog.String("key", key))
			return
		}

		for _, segment = range Segments(key) {
			segment = strings.TrimSpace(segment)
			if segment == "" {
				// "prop..value" || "prop. .value"
				slog.Warn("[cfg] invalid config key.", slog.String("key", key), slog.String("segment", segment))
				return
			}

			// "prop.array[0]" => pkey = "array[0]"
			indexOpen := strings.IndexByte(key, '[')

			if indexOpen == 0 || (indexOpen >= 0 && indexOpen != len(segment)-1) {
				// "prop.[0]" || "prop.array[0]text"
				slog.Warn("[cfg] invalid brackets in config key.", slog.String("key", key), slog.String("segment", segment))
				return
			}

			indexClose := strings.IndexByte(key, ']')

			if indexOpen != indexClose && indexClose-indexOpen < 0 {
				// "prop.array[0" || "prop.array0]"  || "prop.array0]tex[t" || "prop.array[]"
				slog.Warn("[cfg] invalid brackets in config key.", slog.String("key", key), slog.String("segment", segment))
				return
			}

			if indexOpen >= 0 {

			} else {
				prev = curr
				curr[segment] = map[string]any{}
				curr = curr[segment].(map[string]any)
			}

		}
	}

	prev[segment] = value

	c.LoadObject(object)
}

func (c *Env) getStringUnsafe(key string) string {
	if v, exist := c.getValueUnsafe(key); !exist || v == nil {
		return ""
	} else {
		switch s := v.(type) {
		case string:
			return s
		case []string:
			return strings.Join(s, ",")
		default:
			return fmt.Sprintf("%v", s)
		}
	}
}

// getValueUnsafe internal use, non-blocking call. Only use
// when asynchronous access control is active, see get method.
func (c *Env) getValueUnsafe(key string) (any, bool) {

	if e, exist := c.cache[key]; exist {
		return e.value, e.exist
	}

	exist := false
	var value any

	defer func() {
		// Deferred function calls are pushed onto a stack.
		// When a function returns, its deferred calls are executed in last-in-first-out order.
		// https://go.dev/tour/flowcontrol/13
		c.cache[key] = &cacheEntry{value: value, exist: exist}
	}()

	entry := c.getEntryUnsafe(key)

	if entry != nil {
		exist = true
		// lazy string evaluation
		c.expand(entry)
		value = entry.Value()
	}
	return value, exist
}

// getEntryUnsafe internal use, non-blocking call. Only use
// when asynchronous access control is active, see get method.
func (c *Env) getEntryUnsafe(key string) *Entry {

	entry := c.root
	for _, pkey := range Segments(key) {
		// "prop.array[0]" => pkey = "array[0]"
		if strings.HasSuffix(pkey, "]") {
			// array index
			pts := strings.Split(pkey, "[")
			if len(pts) != 2 {
				// invalid format, expects "prop.array[0]"
				return nil
			} else if entry.kind != ObjectKind || entry.value == nil {
				// entry is not an object, so there cannot be a child of the array type
				return nil
			} else if arrEntry, ok := (entry.value.(map[string]*Entry))[pts[0]]; !ok {
				// object does not exist
				break
			} else if arrEntry.kind != ArrayKind || arrEntry.value == nil {
				// data type is not array or is empty array
				return nil
			} else if idx, err := strconv.Atoi(strings.TrimSuffix(pts[1], "]")); err != nil {
				// index is not number
				return nil
			} else if idx < 0 {
				// does not accept negative index
				return nil
			} else if arr := arrEntry.value.([]*Entry); (len(arr) - 1) < idx {
				// array does not have an object with the given index
				return nil
			} else {
				entry = arr[idx]
			}
		} else if entry.kind != ObjectKind || entry.value == nil {
			// entry is not an object, therefore there cannot be a child with the given key
			return nil
		} else if e, ok := (entry.value.(map[string]*Entry))[pkey]; !ok {
			return nil
		} else {
			entry = e
		}
	}

	return entry
}

// expand replaces ${var} or $var in the strings based on the mapping function.
func (c *Env) expand(e *Entry) {
	switch e.kind {
	case StringKind:
		if e.expr != "" && strings.IndexByte(e.value.(string), '$') >= 0 {
			// replaces ${var} or $var in the string
			e.value = os.Expand(e.expr, c.getStringUnsafe)
		}
	case ArrayKind:
		for _, entry := range e.value.([]*Entry) {
			c.expand(entry)
		}
	case ObjectKind:
		for _, entry := range e.value.(map[string]*Entry) {
			c.expand(entry)
		}
	}
}

func (c *Env) lock(read bool) func() {
	if read {
		c.mutex.RLock()
	} else {
		c.mutex.Lock()
	}

	return func() {
		if read {
			c.mutex.RUnlock()
		} else {
			c.mutex.Unlock()
		}
	}
}
