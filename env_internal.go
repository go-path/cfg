package cfg

import (
	"fmt"
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
		c.logger.Warn("invalid config key. { key: %s }", key)
		return
	}

	object := map[string]interface{}{}
	var segment string
	var prev map[string]interface{}
	var curr = object

	if indexOpenBracket == -1 {
		// simple key (fast)

		if strings.IndexByte(key, ']') >= 0 {
			// "prop.array0]"
			c.logger.Warn("invalid brackets in config key. { key: %s }", key)
			return
		}

		for _, segment = range strings.Split(key, ".") {
			segment = strings.TrimSpace(segment)
			if segment == "" {
				// "prop..value" || "prop. .value"
				c.logger.Warn("invalid config key. { key: %s }", key)
				return
			}

			prev = curr
			curr[segment] = map[string]interface{}{}
			curr = curr[segment].(map[string]interface{})
		}

	} else {
		// complex key (array)
		if strings.IndexByte(key, ']') == -1 {
			// "prop.array[0"
			c.logger.Warn("invalid brackets in config key. { key: %s }", key)
			return
		}

		for _, segment = range strings.Split(key, ".") {
			segment = strings.TrimSpace(segment)
			if segment == "" {
				// "prop..value" || "prop. .value"
				c.logger.Warn("invalid config key. { key: %s }", key)
				return
			}

			// "prop.array[0]" => pkey = "array[0]"
			indexOpen := strings.IndexByte(key, '[')

			if indexOpen == 0 || (indexOpen >= 0 && indexOpen != len(segment)-1) {
				// "prop.[0]" || "prop.array[0]text"
				c.logger.Warn("invalid brackets in config key. { key: %s, part: %s }", key, segment)
				return
			}

			indexClose := strings.IndexByte(key, ']')

			if indexOpen != indexClose && indexClose-indexOpen < 0 {
				// "prop.array[0" || "prop.array0]"  || "prop.array0]tex[t" || "prop.array[]"
				c.logger.Warn("invalid brackets in config key. { key: %s, part: %s }", key, segment)
				return
			}

			if indexOpen >= 0 {

			} else {
				prev = curr
				curr[segment] = map[string]interface{}{}
				curr = curr[segment].(map[string]interface{})
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

// getValueUnsafe uso interno, chamada não bloqueante. Só usar quando o controle de acesso asíncrono estiver ativo, ver método get.
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

	entry := c.root
	parts := strings.Split(key, ".")
	for _, pkey := range parts {
		// "prop.array[0]" => pkey = "array[0]"
		if strings.HasSuffix(pkey, "]") {
			// array index
			pts := strings.Split(pkey, "[")
			if len(pts) != 2 {
				// formato inválido, espera "prop.array[0]"
				return nil, false
			} else if entry.kind != ObjectKind || entry.value == nil {
				// entry não é objeto, portanto não pode existir um filho do tipo array
				return nil, false
			} else if arrEntry, ok := (entry.value.(map[string]*Entry))[pts[0]]; !ok {
				// objeto não existe
				break
			} else if arrEntry.kind != ArrayKind || arrEntry.value == nil {
				// tipo de dado não é array ou é array vazio
				return nil, false
			} else if idx, err := strconv.Atoi(strings.TrimSuffix(pts[1], "]")); err != nil {
				// índice não é número
				return nil, false
			} else if idx < 0 {
				// não aceita índice negativo
				return nil, false
			} else if arr := arrEntry.value.([]*Entry); (len(arr) - 1) < idx {
				// array não possui objeto com o índice informado
				return nil, false
			} else {
				entry = arr[idx]
			}
		} else if entry.kind != ObjectKind || entry.value == nil {
			// entry não é objeto, portanto não pode existir um filho com a key informada
			return nil, false
		} else if e, ok := (entry.value.(map[string]*Entry))[pkey]; !ok {
			return nil, false
		} else {
			entry = e
		}
	}

	if entry != nil {
		exist = true

		// lazy string evaluation
		c.expand(entry)
		value = entry.Value()
	}
	return value, exist
}

// getEntryUnsafe uso interno, chamada não bloqueante.
// Só usar quando o controle de acesso asíncrono estiver ativo, ver método get.
func (c *Env) getEntryUnsafe(key string) *Entry {

	entry := c.root
	parts := strings.Split(key, ".")
	for _, pkey := range parts {
		// "prop.array[0]" => pkey = "array[0]"
		if strings.HasSuffix(pkey, "]") {
			// array index
			pts := strings.Split(pkey, "[")
			if len(pts) != 2 {
				// formato inválido, espera "prop.array[0]"
				return nil
			} else if entry.kind != ObjectKind || entry.value == nil {
				// entry não é objeto, portanto não pode existir um filho do tipo array
				return nil
			} else if arrEntry, ok := (entry.value.(map[string]*Entry))[pts[0]]; !ok {
				// objeto não existe
				break
			} else if arrEntry.kind != ArrayKind || arrEntry.value == nil {
				// tipo de dado não é array ou é array vazio
				return nil
			} else if idx, err := strconv.Atoi(strings.TrimSuffix(pts[1], "]")); err != nil {
				// índice não é número
				return nil
			} else if idx < 0 {
				// não aceita índice negativo
				return nil
			} else if arr := arrEntry.value.([]*Entry); (len(arr) - 1) < idx {
				// array não possui objeto com o índice informado
				return nil
			} else {
				entry = arr[idx]
			}
		} else if entry.kind != ObjectKind || entry.value == nil {
			// entry não é objeto, portanto não pode existir um filho com a key informada
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
		break
	case ArrayKind:
		for _, entry := range e.value.([]*Entry) {
			c.expand(entry)
		}
		break
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
