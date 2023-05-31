package cfg

import (
	"strings"
)

func ParseEntry(v any) *Entry {
	if v == nil {
		return nil
	}
	entry := &Entry{}
	switch o := v.(type) {
	case bool:
		entry.kind = BoolKind
		entry.value = o
	case string:
		entry.kind = StringKind
		entry.value = o
		if strings.IndexByte(o, '$') >= 0 {
			entry.expr = o
		}
	case int:
		parseNumber(float64(o), entry)
	case int8:
		parseNumber(float64(o), entry)
	case int16:
		parseNumber(float64(o), entry)
	case int32:
		parseNumber(float64(o), entry)
	case int64:
		parseNumber(float64(o), entry)
	case uint:
		parseNumber(float64(o), entry)
	case uint8:
		parseNumber(float64(o), entry)
	case uint16:
		parseNumber(float64(o), entry)
	case uint32:
		parseNumber(float64(o), entry)
	case uint64:
		parseNumber(float64(o), entry)
	case float32:
		parseNumber(float64(o), entry)
	case float64:
		parseNumber(o, entry)
	case map[string]bool:
		parseEntryMapBool(o, entry)
	case map[string]string:
		parseEntryMapString(o, entry)
	case map[string]int:
		parseEntryMapInt(o, entry)
	case map[string]int8:
		parseEntryMapInt8(o, entry)
	case map[string]int16:
		parseEntryMapInt16(o, entry)
	case map[string]int32:
		parseEntryMapInt32(o, entry)
	case map[string]int64:
		parseEntryMapInt64(o, entry)
	case map[string]uint:
		parseEntryMapUint(o, entry)
	case map[string]uint8:
		parseEntryMapUint8(o, entry)
	case map[string]uint16:
		parseEntryMapUint16(o, entry)
	case map[string]uint32:
		parseEntryMapUint32(o, entry)
	case map[string]uint64:
		parseEntryMapUint64(o, entry)
	case map[string]float32:
		parseEntryMapFloat32(o, entry)
	case map[string]float64:
		parseEntryMapFloat64(o, entry)
	case O:
		parseEntryMap(o, entry)
	case map[string]any:
		parseEntryMap(o, entry)
	case []bool:
		parseEntryArray[bool](o, entry)
	case []string:
		parseEntryArray[string](o, entry)
	case []int:
		parseEntryArray[int](o, entry)
	case []int8:
		parseEntryArray[int8](o, entry)
	case []int16:
		parseEntryArray[int16](o, entry)
	case []int32:
		parseEntryArray[int32](o, entry)
	case []int64:
		parseEntryArray[int64](o, entry)
	case []uint:
		parseEntryArray[uint](o, entry)
	case []uint8:
		parseEntryArray[uint8](o, entry)
	case []uint16:
		parseEntryArray[uint16](o, entry)
	case []uint32:
		parseEntryArray[uint32](o, entry)
	case []uint64:
		parseEntryArray[uint64](o, entry)
	case []float32:
		parseEntryArray[float32](o, entry)
	case []float64:
		parseEntryArray[float64](o, entry)
	case []O:
		parseEntryArray[O](o, entry)
	case []map[string]bool:
		parseEntryArray[map[string]bool](o, entry)
	case []map[string]string:
		parseEntryArray[map[string]string](o, entry)
	case []map[string]int:
		parseEntryArray[map[string]int](o, entry)
	case []map[string]int8:
		parseEntryArray[map[string]int8](o, entry)
	case []map[string]int16:
		parseEntryArray[map[string]int16](o, entry)
	case []map[string]int32:
		parseEntryArray[map[string]int32](o, entry)
	case []map[string]int64:
		parseEntryArray[map[string]int64](o, entry)
	case []map[string]uint:
		parseEntryArray[map[string]uint](o, entry)
	case []map[string]uint8:
		parseEntryArray[map[string]uint8](o, entry)
	case []map[string]uint16:
		parseEntryArray[map[string]uint16](o, entry)
	case []map[string]uint32:
		parseEntryArray[map[string]uint32](o, entry)
	case []map[string]uint64:
		parseEntryArray[map[string]uint64](o, entry)
	case []map[string]float32:
		parseEntryArray[map[string]float32](o, entry)
	case []map[string]float64:
		parseEntryArray[map[string]float64](o, entry)
	case []map[string]any:
		parseEntryArray[map[string]any](o, entry)
	case []any:
		parseEntryArray[any](o, entry)
	default:
		return nil
	}
	return entry
}

func parseNumber(o float64, entry *Entry) {
	entry.kind = NumberKind
	entry.value = o
}

func parseEntryArray[T any](o []T, entry *Entry) {
	var list []*Entry
	for _, oc := range o {
		if ce := ParseEntry(oc); ce != nil {
			list = append(list, ce)
		}
	}
	entry.kind = ArrayKind
	entry.value = list
}

var (
	parseEntryMap        func(obj map[string]any, entry *Entry)
	parseEntryMapBool    func(obj map[string]bool, entry *Entry)
	parseEntryMapString  func(obj map[string]string, entry *Entry)
	parseEntryMapInt     func(obj map[string]int, entry *Entry)
	parseEntryMapInt8    func(obj map[string]int8, entry *Entry)
	parseEntryMapInt16   func(obj map[string]int16, entry *Entry)
	parseEntryMapInt32   func(obj map[string]int32, entry *Entry)
	parseEntryMapInt64   func(obj map[string]int64, entry *Entry)
	parseEntryMapUint    func(obj map[string]uint, entry *Entry)
	parseEntryMapUint8   func(obj map[string]uint8, entry *Entry)
	parseEntryMapUint16  func(obj map[string]uint16, entry *Entry)
	parseEntryMapUint32  func(obj map[string]uint32, entry *Entry)
	parseEntryMapUint64  func(obj map[string]uint64, entry *Entry)
	parseEntryMapFloat32 func(obj map[string]float32, entry *Entry)
	parseEntryMapFloat64 func(obj map[string]float64, entry *Entry)
)

func init() {
	parseEntryMap = entryMapParserFn[any]()
	parseEntryMapBool = entryMapParserFn[bool]()
	parseEntryMapString = entryMapParserFn[string]()
	parseEntryMapInt = entryMapParserFn[int]()
	parseEntryMapInt8 = entryMapParserFn[int8]()
	parseEntryMapInt16 = entryMapParserFn[int16]()
	parseEntryMapInt32 = entryMapParserFn[int32]()
	parseEntryMapInt64 = entryMapParserFn[int64]()
	parseEntryMapUint = entryMapParserFn[uint]()
	parseEntryMapUint8 = entryMapParserFn[uint8]()
	parseEntryMapUint16 = entryMapParserFn[uint16]()
	parseEntryMapUint32 = entryMapParserFn[uint32]()
	parseEntryMapUint64 = entryMapParserFn[uint64]()
	parseEntryMapFloat32 = entryMapParserFn[float32]()
	parseEntryMapFloat64 = entryMapParserFn[float64]()
}

func entryMapParserFn[T any]() func(obj map[string]T, entry *Entry) {
	return func(obj map[string]T, entry *Entry) {
		empty := true
		entries := map[string]*Entry{}
		for k, v := range obj {
			key := strings.TrimSpace(k)
			if strings.IndexByte(key, '#') == 0 {
				// comments
				continue
			}

			if e := ParseEntry(v); e != nil {
				entries[key] = e
				empty = false
			}
		}

		entry.kind = ObjectKind
		if !empty {
			entry.value = entries
		}
	}
}
