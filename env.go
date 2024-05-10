package cfg

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// cacheEntry entry
type cacheEntry struct {
	value any
	exist bool
}

// Env propriedades de configuração de um site ou globais.
type Env struct {
	mutex      sync.RWMutex
	fs         FileSystem
	root       *Entry
	cache      map[string]*cacheEntry
	logger     Logger
	fileExts   map[string]UnmarshalFn
	filePaths  []string
	profileKey string
}

// New 6) Default config
func New(defaults ...O) *Env {
	config := &Env{
		root:       &Entry{kind: ObjectKind, value: map[string]*Entry{}},
		cache:      map[string]*cacheEntry{},
		logger:     defaultLogger(),
		fileExts:   map[string]UnmarshalFn{"json": JsonUnmarshal},
		filePaths:  []string{"config"},
		profileKey: "profiles",
	}

	if len(defaults) > 0 {
		for _, cfg := range defaults {
			config.LoadObject(cfg)
		}
	}

	return config
}

// Get a configuration property
func (c *Env) Get(key string) any {
	if e, exist := c.get(key); exist {
		return e
	}
	return nil
}

// Set a configuration property
func (c *Env) Set(key string, value any) {
	if strings.IndexByte(key, '.') == -1 {
		c.LoadObject(O{key: value})
	} else {
		c.set(key, value)
	}
}

// Bool get a boolean value
func (c *Env) Bool(key string) bool {
	v := c.Get(key)
	if v == nil {
		return false
	}
	switch s := v.(type) {
	case bool:
		return s
	case float64:
		return s > 0
	case string:
		return s != ""
	default:
		return true
	}
}

// String get a string value
func (c *Env) String(key string, def ...string) string {
	v := c.Get(key)

	return c.toString(v, def...)
}

func (c *Env) Strings(key string, def ...[]string) []string {
	v := c.Get(key)
	if v == nil {
		if len(def) > 0 {
			return def[0]
		}
		return []string{}
	}
	switch s := v.(type) {
	case []interface{}:
		var list []string
		for _, it := range s {
			list = append(list, c.toString(it))
		}
		return list
	default:
		return []string{c.toString(s)}
	}
}

func (c *Env) toString(v any, def ...string) string {
	var out string
	if v != nil {
		switch s := v.(type) {
		case string:
			out = s
		case float64:
			out = fmt.Sprintf("%v", s)
		case map[string]interface{}:
			if b, err := json.Marshal(s); err != nil {
				c.logger.Warn("cannot convert value to string using json.Marshal", s, err)
				out = strings.TrimPrefix(fmt.Sprintf("%#v", s), "map[string]interface {}")
			} else {
				out = string(b)
			}
		case []interface{}:
			if b, err := json.Marshal(s); err != nil {
				c.logger.Error("cannot convert value to string using json.Marshal", s, err)
			} else {
				out = string(b)
			}
		default:
			out = fmt.Sprintf("%#v", s)
		}
	}

	out = strings.TrimSpace(out)
	if out == "" {
		for _, o := range def {
			if o != "" {
				out = o
				break
			}
		}
	}

	return out
}

// Duration obtém uma duração da config
func (c *Env) Duration(key string, def ...time.Duration) time.Duration {
	value := c.String(key)
	if out, err := time.ParseDuration(value); err == nil {
		return out
	} else {
		c.logger.Error("could not be converted to time.Duration. { key: %s, value: %s }", key, value, err)
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

// Time obtém uma duração da config
func (c *Env) Time(key string, def ...time.Time) time.Time {
	return c.TimeLayout(key, time.RFC3339, def...)
}

func (c *Env) DateTime(key string, def ...time.Time) time.Time {
	return c.TimeLayout(key, time.DateTime, def...)
}

func (c *Env) DateOnly(key string, def ...time.Time) time.Time {
	return c.TimeLayout(key, time.DateOnly, def...)
}

func (c *Env) TimeOnly(key string, def ...time.Time) time.Time {
	return c.TimeLayout(key, time.TimeOnly, def...)
}

// TimeLayout obtém uma duração da config
func (c *Env) TimeLayout(key string, layout string, def ...time.Time) time.Time {
	value := c.String(key)
	if out, err := time.Parse(layout, value); err == nil {
		return out
	} else {
		c.logger.Error("could not be converted to time.Time. { key: %s, value: %s, layout:  }", key, value, layout, err)
		if len(def) > 0 {
			return def[0]
		}
		return time.Time{}
	}
}

// Clone make a copy of the config
func (c *Env) Clone() *Env {
	o := New()
	o.root = c.root.Clone()
	return o
}

func (c *Env) Keys(key string) []string {
	unlock := c.lock(true)
	defer unlock()

	entry := c.getEntryUnsafe(key)
	if entry == nil {
		return nil
	}
	switch entry.kind {
	case BoolKind, StringKind, NumberKind:
		return nil
	case ArrayKind:
		var list []string
		value := entry.value.([]*Entry)
		for i, _ := range value {
			list = append(list, "["+strconv.Itoa(i)+"]")
		}
		return list
	default:
		var list []string
		value := entry.value.(map[string]*Entry)
		for k, _ := range value {
			list = append(list, k)
		}
		return list
	}
}

// Merge faz o merge de src na config atual
func (c *Env) Merge(src *Env) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.root.Merge(src.root)
	c.cache = map[string]*cacheEntry{}

	// clear expressions
	c.root.Walk(func(entry *Entry) {
		if entry.expr != "" {
			entry.value = entry.expr
		}
	})
}
