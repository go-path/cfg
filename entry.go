package cfg

// EntryKind representa os tipos de dados suportados na configuração, afim de manter compatibilidade total com arquivos JSON
type EntryKind uint

const (
	BoolKind   EntryKind = iota // bool, for JSON booleans
	NumberKind                  // float64, for JSON numbers
	StringKind                  // string, for JSON strings
	ArrayKind                   // []*Entry{}, for JSON arrays
	ObjectKind                  // map[string]*Entry{}, for JSON objects
)

// Entry todas as propriedades são mapeadas para um tipo de dado abaixo, afim de respeitar corretamente
// a integração com Json
type Entry struct {
	kind  EntryKind // Tipagem do valor
	value any       // Valor salvo (bool, float64, string, []*Entry, map[string]*Entry)
	expr  string    // Quando string com expressão (${var} | $var)
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

// Merge faz a mesclagem de dois objetos
func (e *Entry) Merge(other *Entry) {
	switch other.kind {
	case BoolKind, StringKind, NumberKind, ArrayKind:
		e.value = other.value
		break
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
		break
	}
}

// Clone faz um deep copy dessa entrada
func (e *Entry) Clone() *Entry {
	other := &Entry{kind: e.kind, value: e.value, expr: e.expr}

	switch other.kind {
	case ArrayKind:
		var value []*Entry
		for _, entry := range e.value.([]*Entry) {
			value = append(value, entry.Clone())
		}
		other.value = value
		break
	case ObjectKind:
		if e.value == nil {
			break
		}
		value := map[string]*Entry{}
		for key, src := range e.value.(map[string]*Entry) {
			value[key] = src.Clone()
		}
		other.value = value
		break
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
