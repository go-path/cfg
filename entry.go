package cfg

// EntryKind represents the data types supported in the
// configuration, in order to maintain full compatibility
// with JSON files
type EntryKind uint

const (
	BoolKind   EntryKind = iota // bool, for JSON booleans
	NumberKind                  // float64, for JSON numbers
	StringKind                  // string, for JSON strings
	ArrayKind                   // []*Entry{}, for JSON arrays
	ObjectKind                  // map[string]*Entry{}, for JSON objects
)

// Entry all properties are mapped to a data type below,
// in order to correctly respect the integration with Json
type Entry struct {
	kind  EntryKind // Value typing
	value any       // Saved value (bool, float64, string, []*Entry, map[string]*Entry)
	expr  string    // When string with expression (${var} | $var)
}

func (e *Entry) Kind() EntryKind {
	return e.kind
}

func (e *Entry) Value() any {
	switch e.kind {
	case BoolKind, StringKind, NumberKind:
		return e.value
	case ArrayKind:
		var list []any
		value := e.value.([]*Entry)
		for _, entry := range value {
			list = append(list, entry.Value())
		}
		return list
	default:
		obj := map[string]any{}
		value := e.value.(map[string]*Entry)
		for key, entry := range value {
			obj[key] = entry.Value()
		}
		return obj
	}
}

// Merge merges two objects
func (e *Entry) Merge(other *Entry) {
	switch other.kind {
	case BoolKind, StringKind, NumberKind, ArrayKind:
		e.value = other.value
	case ObjectKind:
		if e.value == nil || e.kind != ObjectKind {
			e.value = map[string]*Entry{}
		}
		e.kind = ObjectKind

		target := e.value.(map[string]*Entry)
		source := other.value.(map[string]*Entry)
		for key, src := range source {
			if dest, exist := target[key]; !exist {
				target[key] = src
			} else if dest == nil || src.kind != dest.kind {
				target[key] = src
			} else if src.value == nil {
				delete(target, key)
			} else {
				dest.Merge(src)
			}
		}
	}
}

// Clone makes a deep copy of the entry
func (e *Entry) Clone() *Entry {
	other := &Entry{kind: e.kind, value: e.value, expr: e.expr}

	switch other.kind {
	case ArrayKind:
		var value []*Entry
		for _, entry := range e.value.([]*Entry) {
			value = append(value, entry.Clone())
		}
		other.value = value
	case ObjectKind:
		if e.value == nil {
			break
		}
		value := map[string]*Entry{}
		for key, src := range e.value.(map[string]*Entry) {
			value[key] = src.Clone()
		}
		other.value = value
	}

	return other
}

func (e *Entry) Walk(visitor func(*Entry)) {
	switch e.kind {
	case BoolKind, StringKind, NumberKind:
		visitor(e)
	case ArrayKind:
		for _, entry := range e.value.([]*Entry) {
			entry.Walk(visitor)
		}
	default:
		if e.value != nil {
			for _, entry := range e.value.(map[string]*Entry) {
				entry.Walk(visitor)
			}
		}
	}
}
