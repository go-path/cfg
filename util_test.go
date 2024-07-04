package cfg

import (
	"reflect"
	"testing"
)

func Test_Segments(t *testing.T) {
	tests := []testAny{
		{key: "app.name", want: []string{"app", "name"}},
		{key: "app.description", want: []string{"app", "description"}},
		{key: "dev.watch.interval", want: []string{"dev", "watch", "interval"}},
		{key: "dev.watch.array[0].id", want: []string{"dev", "watch", "array[0]", "id"}},
		{key: "dev.watch.array[1].id", want: []string{"dev", "watch", "array[1]", "id"}},
		{key: "assets.alias.bootstrap\\.css.filepath", want: []string{"assets", "alias", "bootstrap.css", "filepath"}},
		{key: "assets.required.bootstrap\\.js", want: []string{"assets", "required", "bootstrap.js"}},
	}
	for _, tt := range tests {
		t.Run(tt.key+"+def", func(t *testing.T) {
			if got := Segments(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Escape(t *testing.T) {
	tests := []testAny{
		{key: "bootstrap", want: "bootstrap"},
		{key: "bootstrap.js", want: "bootstrap\\.js"},
		{key: "bootstrap\\.js", want: "bootstrap\\.js"},
		{key: "app.name", want: "app\\.name"},
		{key: "app.description", want: "app\\.description"},
		{key: "dev.watch.interval", want: "dev\\.watch\\.interval"},
		{key: "dev.watch.array[0].id", want: "dev\\.watch\\.array[0]\\.id"},
		{key: "assets.alias.bootstrap.css.filepath", want: "assets\\.alias\\.bootstrap\\.css\\.filepath"},
		{key: "assets\\.alias.bootstrap\\.css.filepath", want: "assets\\.alias\\.bootstrap\\.css\\.filepath"},
		{key: "assets.required.bootstrap.js", want: "assets\\.required\\.bootstrap\\.js"},
	}
	for _, tt := range tests {
		t.Run(tt.key+"+def", func(t *testing.T) {
			if got := Escape(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
