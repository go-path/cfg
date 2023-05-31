package cfg

import (
	"reflect"
	"testing"
	"time"
)

type testKW struct {
	key  string
	want any
}

var testSignedK = []string{"int", "int8", "int16", "int32", "int64", "float32", "float64"}
var testUnsignedK = []string{"uint", "uint8", "uint16", "uint32", "uint64"}
var testConfig = O{
	"# ignored": "ignored",
	"nil":       nil,
	"bool":      map[string]bool{"true": true, "false": false},
	"string":    map[string]string{"empty": "", "value": "1", "template": "${string.value}"},
	"int":       map[string]int{"pos": int(1), "zero": int(0), "neg": int(-1)},
	"int8":      map[string]int8{"pos": int8(1), "zero": int8(0), "neg": int8(-1)},
	"int16":     map[string]int16{"pos": int16(1), "zero": int16(0), "neg": int16(-1)},
	"int32":     map[string]int32{"pos": int32(1), "zero": int32(0), "neg": int32(-1)},
	"int64":     map[string]int64{"pos": int64(1), "zero": int64(0), "neg": int64(-1)},
	"uint":      map[string]uint{"pos": uint(1), "zero": uint(0)},
	"uint8":     map[string]uint8{"pos": uint8(1), "zero": uint8(0)},
	"uint16":    map[string]uint16{"pos": uint16(1), "zero": uint16(0)},
	"uint32":    map[string]uint32{"pos": uint32(1), "zero": uint32(0)},
	"uint64":    map[string]uint64{"pos": uint64(1), "zero": uint64(0)},
	"float32":   map[string]float32{"pos": float32(1), "zero": float32(0), "neg": float32(-1)},
	"float64":   map[string]float64{"pos": float64(1), "zero": float64(0), "neg": float64(-1), "pi": 3.1415927410125732},
	"map":       map[string]any{"key": "value"},
	"object":    O{"key": "value"},
	"duration": O{
		"milis": "300ms",
		"hour":  "2h45m",
	},
	"array": O{
		"map":     []map[string]any{{"key": "value"}},
		"object":  []O{{"key": "value"}},
		"bool":    []bool{true, false},
		"string":  []string{"", "1", "${array.string[1]}"},
		"int":     []int{int(-1), int(0), int(1)},
		"int8":    []int8{int8(-1), int8(0), int8(1)},
		"int16":   []int16{int16(-1), int16(0), int16(1)},
		"int32":   []int32{int32(-1), int32(0), int32(1)},
		"int64":   []int64{int64(-1), int64(0), int64(1)},
		"uint":    []uint{uint(0), uint(1)},
		"uint8":   []uint8{uint8(0), uint8(1)},
		"uint16":  []uint16{uint16(0), uint16(1)},
		"uint32":  []uint32{uint32(0), uint32(1)},
		"uint64":  []uint64{uint64(0), uint64(1)},
		"float32": []float32{float32(-1), float32(0), float32(1)},
		"float64": []float64{float64(-1), float64(0), float64(1)},
		"complex": O{
			"bool":    []map[string]bool{{"true": true, "false": false}},
			"string":  []map[string]string{{"empty": "", "value": "1", "template": "${string.value}"}},
			"int":     []map[string]int{{"pos": int(1), "zero": int(0), "neg": int(-1)}},
			"int8":    []map[string]int8{{"pos": int8(1), "zero": int8(0), "neg": int8(-1)}},
			"int16":   []map[string]int16{{"pos": int16(1), "zero": int16(0), "neg": int16(-1)}},
			"int32":   []map[string]int32{{"pos": int32(1), "zero": int32(0), "neg": int32(-1)}},
			"int64":   []map[string]int64{{"pos": int64(1), "zero": int64(0), "neg": int64(-1)}},
			"uint":    []map[string]uint{{"pos": uint(1), "zero": uint(0)}},
			"uint8":   []map[string]uint8{{"pos": uint8(1), "zero": uint8(0)}},
			"uint16":  []map[string]uint16{{"pos": uint16(1), "zero": uint16(0)}},
			"uint32":  []map[string]uint32{{"pos": uint32(1), "zero": uint32(0)}},
			"uint64":  []map[string]uint64{{"pos": uint64(1), "zero": uint64(0)}},
			"float32": []map[string]float32{{"pos": float32(1), "zero": float32(0), "neg": float32(-1)}},
			"float64": []map[string]float64{{"pos": float64(1), "zero": float64(0), "neg": float64(-1), "pi": 3.1415927410125732}},
			"map":     []map[string]any{{"key": "value"}},
		},
	},
}

var testKeys []string

func init() {

	testKeys = append(testKeys, "nil", "# ignored", "array")
	testKeys = append(testKeys, "map", "map.key", "array.map", "array.map[0]", "array.map[0].key")
	testKeys = append(testKeys,
		"array.complex.map",
		"array.complex.map[0]",
		"array.complex.map[0].key",
	)

	testKeys = append(testKeys, "object", "object.key", "array.object", "array.object[0]", "array.object[0].key")

	testKeys = append(testKeys, "bool", "bool.true", "bool.false")
	testKeys = append(testKeys, "array.bool", "array.bool[0]", "array.bool[1]")
	testKeys = append(testKeys,
		"array.complex.bool",
		"array.complex.bool[0]",
		"array.complex.bool[0].true",
		"array.complex.bool[0].false",
	)

	testKeys = append(testKeys, "string", "string.empty", "string.value", "string.template")
	testKeys = append(testKeys, "array.string", "array.string[0]", "array.string[1]", "array.string[2]")
	testKeys = append(testKeys,
		"array.complex.string",
		"array.complex.string[0]",
		"array.complex.string[0].empty",
		"array.complex.string[0].value",
		"array.complex.string[0].template",
	)

	for _, name := range testSignedK {
		testKeys = append(testKeys, name)
		testKeys = append(testKeys, name+".neg")
		testKeys = append(testKeys, name+".zero")
		testKeys = append(testKeys, name+".pos")
		testKeys = append(testKeys, "array."+name+"[0]")
		testKeys = append(testKeys, "array."+name+"[1]")
		testKeys = append(testKeys, "array."+name+"[2]")

		testKeys = append(testKeys, "array.complex"+name+"[0]")
		testKeys = append(testKeys, "array.complex"+name+"[0].neg")
		testKeys = append(testKeys, "array.complex"+name+"[0].zero")
		testKeys = append(testKeys, "array.complex"+name+"[0].pos")
	}

	for _, name := range testUnsignedK {
		testKeys = append(testKeys, name)
		testKeys = append(testKeys, name+".zero")
		testKeys = append(testKeys, name+".pos")

		testKeys = append(testKeys, "array."+name+"[0]")
		testKeys = append(testKeys, "array."+name+"[1]")

		testKeys = append(testKeys, "array.complex"+name+"[0]")
		testKeys = append(testKeys, "array.complex"+name+"[0].zero")
		testKeys = append(testKeys, "array.complex"+name+"[0].pos")
	}
}

func testNumbers(positive, zero, negative any) []testKW {
	var values []testKW
	for _, name := range testSignedK {
		values = append(values, testKW{key: name + ".neg", want: negative})
		values = append(values, testKW{key: name + ".zero", want: zero})
		values = append(values, testKW{key: name + ".pos", want: positive})

		values = append(values, testKW{key: "array." + name + "[0]", want: negative})
		values = append(values, testKW{key: "array." + name + "[1]", want: zero})
		values = append(values, testKW{key: "array." + name + "[2]", want: positive})
	}

	for _, name := range testUnsignedK {
		values = append(values, testKW{key: name + ".zero", want: zero})
		values = append(values, testKW{key: name + ".pos", want: positive})

		values = append(values, testKW{key: "array." + name + "[0]", want: zero})
		values = append(values, testKW{key: "array." + name + "[1]", want: positive})
	}

	return values
}

func TestEnv_Get(t *testing.T) {
	env := New(map[string]interface{}{
		"# ignored": "",
		"app": map[string]interface{}{
			"name":        "My App",
			"description": "${app.name} is a Syntax application",
		},
		"dev": map[string]interface{}{
			"port":       float64(8443),
			"active":     true,
			"debug":      true,
			"livereload": true,
			"watch": map[string]interface{}{
				"interval": float64(120),
				"quiet":    float64(60),
				"array": []interface{}{
					map[string]interface{}{"id": float64(2)},
					map[string]interface{}{"id": float64(4)},
				},
				"names": []interface{}{"Darwin", "Gabriel", "Mara"},
			},
		},
	})

	// merge
	env.LoadObject(map[string]interface{}{
		"app": map[string]interface{}{
			"name":   "My Super App",
			"author": "Alex Rodin",
		},
		"dev": map[string]interface{}{
			"livereload": false,
			"watch": map[string]interface{}{
				"interval": float64(240),
				"array": []interface{}{
					map[string]interface{}{"id": float64(8)},
					map[string]interface{}{"id": float64(16)},
				},
				"names": []interface{}{"Darwin R.", "Gabriel R.", "Mara R.", "${app.author}"},
			},
		},
	})

	tests := []struct {
		key  string
		want any
	}{
		{key: "app.name", want: "My Super App"},
		{key: "app.description", want: "My Super App is a Syntax application"},
		{key: "dev.livereload", want: false},
		{key: "dev.watch.interval", want: float64(240)},
		{key: "dev.watch.array[0].id", want: float64(8)},
		{key: "dev.watch.array[1].id", want: float64(16)},
		{key: "dev.watch.array[2].id", want: nil},
		{key: "dev.watch.names", want: []interface{}{"Darwin R.", "Gabriel R.", "Mara R.", "Alex Rodin"}},
		{key: "dev.watch.names[0]", want: "Darwin R."},
		{key: "dev.watch.names[1]", want: "Gabriel R."},
		{key: "dev.watch.names[2]", want: "Mara R."},
		{key: "dev.watch.names[3]", want: "Alex Rodin"},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := env.Get(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_Set(t *testing.T) {

	env := New(nil)

	tests := []struct {
		key   string
		value any
		want  any
	}{
		{"name", "Ramon", "Ramon"},
		{"app.name", "My App", "My App"},
		{"my.big.path.variable", "Value", "Value"},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			env.Set(tt.key, tt.value)
			if got := env.Get(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_Clone(t *testing.T) {
	env := New(testConfig)
	clone := env.Clone()

	for _, key := range testKeys {
		t.Run(key, func(t *testing.T) {
			got := clone.Get(key)
			want := env.Get(key)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Get() = %v, want %v", got, want)
			}
		})
	}
}

func TestEnv_Merge(t *testing.T) {
	envA := New(testConfig)
	envB := New(O{"string": map[string]string{"value": "2"}})

	if got := envA.Get("string.template"); !reflect.DeepEqual(got, "1") {
		t.Errorf("Get() = %v, want %v", got, "1")
	}

	if got := envB.Get("string.value"); !reflect.DeepEqual(got, "2") {
		t.Errorf("Get() = %v, want %v", got, "2")
	}

	envA.Merge(envB)

	if got := envA.Get("string.value"); !reflect.DeepEqual(got, "2") {
		t.Errorf("Get() = %v, want %v", got, "2")
	}

	if got := envA.Get("string.template"); !reflect.DeepEqual(got, "2") {
		t.Errorf("Get() = %v, want %v", got, "2")
	}
}

func TestEnv_Bool(t *testing.T) {
	env := New(testConfig)
	tests := []testKW{
		{key: "nil", want: false},
		{key: "object", want: true},
		{key: "bool.true", want: true},
		{key: "bool.false", want: false},
		{key: "object", want: true},
		{key: "array.bool[0]", want: true},
		{key: "array.bool[1]", want: false},
		{key: "string.empty", want: false},
		{key: "string.value", want: true},
		{key: "array.string[0]", want: false},
		{key: "array.string[1]", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := env.Bool(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range testNumbers(true, false, false) {
		t.Run(tt.key, func(t *testing.T) {
			if got := env.Bool(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_String(t *testing.T) {
	env := New(testConfig)
	tests := []testKW{
		{key: "nil", want: ""},
		{key: "bool.true", want: "true"},
		{key: "bool.false", want: "false"},
		{key: "string.empty", want: ""},
		{key: "string.value", want: "1"},
		{key: "string.template", want: "1"},
		{key: "map", want: `{"key":"value"}`},
		{key: "map.key", want: "value"},
		{key: "object", want: `{"key":"value"}`},
		{key: "object.key", want: "value"},
		{key: "array.map", want: `[{"key":"value"}]`},
		{key: "array.object", want: `[{"key":"value"}]`},
		{key: "array.bool", want: `[true,false]`},
		{key: "array.string", want: `["","1","1"]`},
		{key: "array.int", want: `[-1,0,1]`},
		{key: "array.int8", want: `[-1,0,1]`},
		{key: "array.int16", want: `[-1,0,1]`},
		{key: "array.int32", want: `[-1,0,1]`},
		{key: "array.int64", want: `[-1,0,1]`},
		{key: "array.uint", want: `[0,1]`},
		{key: "array.uint8", want: `[0,1]`},
		{key: "array.uint16", want: `[0,1]`},
		{key: "array.uint32", want: `[0,1]`},
		{key: "array.uint64", want: `[0,1]`},
		{key: "array.float32", want: `[-1,0,1]`},
		{key: "array.float64", want: `[-1,0,1]`},
		{key: "float64.pi", want: `3.1415927410125732`},
		{key: "bool", want: `{"false":false,"true":true}`},
		{key: "string", want: `{"empty":"","template":"1","value":"1"}`},
		{key: "int", want: `{"neg":-1,"pos":1,"zero":0}`},
		{key: "int8", want: `{"neg":-1,"pos":1,"zero":0}`},
		{key: "int16", want: `{"neg":-1,"pos":1,"zero":0}`},
		{key: "int32", want: `{"neg":-1,"pos":1,"zero":0}`},
		{key: "int64", want: `{"neg":-1,"pos":1,"zero":0}`},
		{key: "float32", want: `{"neg":-1,"pos":1,"zero":0}`},
		{key: "float64", want: `{"neg":-1,"pi":3.1415927410125732,"pos":1,"zero":0}`},
		{key: "uint", want: `{"pos":1,"zero":0}`},
		{key: "uint8", want: `{"pos":1,"zero":0}`},
		{key: "uint16", want: `{"pos":1,"zero":0}`},
		{key: "uint32", want: `{"pos":1,"zero":0}`},
		{key: "uint64", want: `{"pos":1,"zero":0}`},
		{key: "map", want: `{"key":"value"}`},
		{key: "object", want: `{"key":"value"}`},
		{key: "array.complex.bool", want: `[{"false":false,"true":true}]`},
		{key: "array.complex.string", want: `[{"empty":"","template":"1","value":"1"}]`},
		{key: "array.complex.int", want: `[{"neg":-1,"pos":1,"zero":0}]`},
		{key: "array.complex.int8", want: `[{"neg":-1,"pos":1,"zero":0}]`},
		{key: "array.complex.int16", want: `[{"neg":-1,"pos":1,"zero":0}]`},
		{key: "array.complex.int32", want: `[{"neg":-1,"pos":1,"zero":0}]`},
		{key: "array.complex.int64", want: `[{"neg":-1,"pos":1,"zero":0}]`},
		{key: "array.complex.float32", want: `[{"neg":-1,"pos":1,"zero":0}]`},
		{key: "array.complex.float64", want: `[{"neg":-1,"pi":3.1415927410125732,"pos":1,"zero":0}]`},
		{key: "array.complex.uint", want: `[{"pos":1,"zero":0}]`},
		{key: "array.complex.uint8", want: `[{"pos":1,"zero":0}]`},
		{key: "array.complex.uint16", want: `[{"pos":1,"zero":0}]`},
		{key: "array.complex.uint32", want: `[{"pos":1,"zero":0}]`},
		{key: "array.complex.uint64", want: `[{"pos":1,"zero":0}]`},
		{key: "array.complex.map", want: `[{"key":"value"}]`},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := env.String(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}

	// default value
	testsDefault := []string{
		"nil",
		"map.invalid.key",
	}
	for _, tt := range testsDefault {
		want := "2"
		t.Run(tt, func(t *testing.T) {
			if got := env.String(tt, want); !reflect.DeepEqual(got, want) {
				t.Errorf("Get() = %v, want %v", got, want)
			}
		})
	}

	for _, tt := range testNumbers("1", "0", "-1") {
		t.Run(tt.key, func(t *testing.T) {
			if got := env.String(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_Strings(t *testing.T) {
	env := New(testConfig)
	tests := []testKW{
		{key: "nil", want: []string{}},
		{key: "bool.true", want: []string{"true"}},
		{key: "bool.false", want: []string{"false"}},
		{key: "string.empty", want: []string{""}},
		{key: "string.value", want: []string{"1"}},
		{key: "string.template", want: []string{"1"}},
		{key: "map", want: []string{`{"key":"value"}`}},
		{key: "map.key", want: []string{"value"}},
		{key: "object", want: []string{`{"key":"value"}`}},
		{key: "object.key", want: []string{"value"}},
		{key: "array.map", want: []string{`{"key":"value"}`}},
		{key: "array.object", want: []string{`{"key":"value"}`}},
		{key: "array.bool", want: []string{"true", "false"}},
		{key: "array.string", want: []string{"", "1", "1"}},
		{key: "array.int", want: []string{"-1", "0", "1"}},
		{key: "array.int8", want: []string{"-1", "0", "1"}},
		{key: "array.int16", want: []string{"-1", "0", "1"}},
		{key: "array.int32", want: []string{"-1", "0", "1"}},
		{key: "array.int64", want: []string{"-1", "0", "1"}},
		{key: "array.uint", want: []string{"0", "1"}},
		{key: "array.uint8", want: []string{"0", "1"}},
		{key: "array.uint16", want: []string{"0", "1"}},
		{key: "array.uint32", want: []string{"0", "1"}},
		{key: "array.uint64", want: []string{"0", "1"}},
		{key: "array.float32", want: []string{"-1", "0", "1"}},
		{key: "array.float64", want: []string{"-1", "0", "1"}},
		{key: "float64.pi", want: []string{`3.1415927410125732`}},
		{key: "bool", want: []string{`{"false":false,"true":true}`}},
		{key: "string", want: []string{`{"empty":"","template":"1","value":"1"}`}},
		{key: "int", want: []string{`{"neg":-1,"pos":1,"zero":0}`}},
		{key: "int8", want: []string{`{"neg":-1,"pos":1,"zero":0}`}},
		{key: "int16", want: []string{`{"neg":-1,"pos":1,"zero":0}`}},
		{key: "int32", want: []string{`{"neg":-1,"pos":1,"zero":0}`}},
		{key: "int64", want: []string{`{"neg":-1,"pos":1,"zero":0}`}},
		{key: "float32", want: []string{`{"neg":-1,"pos":1,"zero":0}`}},
		{key: "float64", want: []string{`{"neg":-1,"pi":3.1415927410125732,"pos":1,"zero":0}`}},
		{key: "uint", want: []string{`{"pos":1,"zero":0}`}},
		{key: "uint8", want: []string{`{"pos":1,"zero":0}`}},
		{key: "uint16", want: []string{`{"pos":1,"zero":0}`}},
		{key: "uint32", want: []string{`{"pos":1,"zero":0}`}},
		{key: "uint64", want: []string{`{"pos":1,"zero":0}`}},
		{key: "map", want: []string{`{"key":"value"}`}},
		{key: "object", want: []string{`{"key":"value"}`}},
		{key: "array.complex.bool", want: []string{`{"false":false,"true":true}`}},
		{key: "array.complex.string", want: []string{`{"empty":"","template":"1","value":"1"}`}},
		{key: "array.complex.int", want: []string{`{"neg":-1,"pos":1,"zero":0}`}},
		{key: "array.complex.int8", want: []string{`{"neg":-1,"pos":1,"zero":0}`}},
		{key: "array.complex.int16", want: []string{`{"neg":-1,"pos":1,"zero":0}`}},
		{key: "array.complex.int32", want: []string{`{"neg":-1,"pos":1,"zero":0}`}},
		{key: "array.complex.int64", want: []string{`{"neg":-1,"pos":1,"zero":0}`}},
		{key: "array.complex.float32", want: []string{`{"neg":-1,"pos":1,"zero":0}`}},
		{key: "array.complex.float64", want: []string{`{"neg":-1,"pi":3.1415927410125732,"pos":1,"zero":0}`}},
		{key: "array.complex.uint", want: []string{`{"pos":1,"zero":0}`}},
		{key: "array.complex.uint8", want: []string{`{"pos":1,"zero":0}`}},
		{key: "array.complex.uint16", want: []string{`{"pos":1,"zero":0}`}},
		{key: "array.complex.uint32", want: []string{`{"pos":1,"zero":0}`}},
		{key: "array.complex.uint64", want: []string{`{"pos":1,"zero":0}`}},
		{key: "array.complex.map", want: []string{`{"key":"value"}`}},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := env.Strings(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}

	// default value
	testsDefault := []string{
		"nil",
		"map.invalid.key",
	}
	for _, tt := range testsDefault {
		want := []string{"2"}
		t.Run(tt, func(t *testing.T) {
			if got := env.Strings(tt, want); !reflect.DeepEqual(got, want) {
				t.Errorf("Get() = %v, want %v", got, want)
			}
		})
	}

	for _, tt := range testNumbers([]string{"1"}, []string{"0"}, []string{"-1"}) {
		t.Run(tt.key, func(t *testing.T) {
			if got := env.Strings(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_Float(t *testing.T) {
	env := New(testConfig)
	tests := []testKW{
		{key: "nil", want: float64(0)},
		{key: "object", want: float64(0)},
		{key: "bool.true", want: float64(1)},
		{key: "bool.false", want: float64(0)},
		{key: "string.empty", want: float64(0)},
		{key: "string.value", want: float64(1)},
		{key: "map", want: float64(0)},
		{key: "map.key", want: float64(0)},
		{key: "object", want: float64(0)},
		{key: "object.key", want: float64(0)},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := env.Float(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}

	// default value
	testsDefault := []string{
		"nil",
		"map.invalid.key",
	}
	for _, tt := range testsDefault {
		want := float64(1)
		t.Run(tt, func(t *testing.T) {
			if got := env.Float(tt, want); !reflect.DeepEqual(got, want) {
				t.Errorf("Get() = %v, want %v", got, want)
			}
		})
	}

	for _, tt := range testNumbers(float64(1), float64(0), float64(-1)) {
		t.Run(tt.key, func(t *testing.T) {
			if got := env.Float(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_Int(t *testing.T) {
	env := New(testConfig)
	tests := []testKW{
		{key: "nil", want: int(0)},
		{key: "object", want: int(0)},
		{key: "bool.true", want: int(1)},
		{key: "bool.false", want: int(0)},
		{key: "string.empty", want: int(0)},
		{key: "string.value", want: int(1)},
		{key: "map", want: int(0)},
		{key: "map.key", want: int(0)},
		{key: "object", want: int(0)},
		{key: "object.key", want: int(0)},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := env.Int(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}

	// default value
	testsDefault := []string{
		"nil",
		"map.invalid.key",
	}
	for _, tt := range testsDefault {
		want := int(1)
		t.Run(tt, func(t *testing.T) {
			if got := env.Int(tt, want); !reflect.DeepEqual(got, want) {
				t.Errorf("Get() = %v, want %v", got, want)
			}
		})
	}

	for _, tt := range testNumbers(int(1), int(0), int(-1)) {
		t.Run(tt.key, func(t *testing.T) {
			if got := env.Int(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_Duration(t *testing.T) {
	env := New(testConfig)
	tests := []testKW{
		{key: "nil", want: time.Millisecond * 0},
		{key: "duration.milis", want: time.Millisecond * 300},
		{key: "duration.hour", want: time.Hour*2 + time.Minute*45},
		{key: "object", want: time.Millisecond * 0},
	}
	for _, tt := range tests {
		t.Run(tt.key+"+def", func(t *testing.T) {
			if got := env.Duration(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}

	def := time.Millisecond * 5
	testsDef := []testKW{
		{key: "nil", want: def},
		{key: "duration.milis", want: time.Millisecond * 300},
		{key: "duration.hour", want: time.Hour*2 + time.Minute*45},
		{key: "object", want: def},
	}
	for _, tt := range testsDef {
		t.Run(tt.key+"+def", func(t *testing.T) {
			if got := env.Duration(tt.key, def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
