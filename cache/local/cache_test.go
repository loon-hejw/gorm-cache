package local

import (
	"reflect"
	"testing"
	"time"

	"github.com/hejw123/gorm-cache/cache/types"
)

type TestData struct {
	key string
	val string
}

func Test_Cache(t *testing.T) {

	cache := New(types.Local{
		Config: types.Config{
			Expiration: 5 * time.Minute,
		},
		DefaultExpiration: 30 * time.Minute,
		CleanupInterval:   10 * time.Minute,
	})

	test := []struct {
		key  string
		val  interface{}
		want interface{}
	}{
		{key: "a", val: "a", want: "a"},
		{key: "b", val: 1, want: 1},
		{key: "c", val: 2.1, want: 2.1},
		{key: "d", val: []byte("abcd"), want: []byte("abcd")},
		{key: "e", val: TestData{key: "1", val: "2"}, want: TestData{key: "1", val: "2"}},
	}

	for _, tt := range test {
		if err := cache.Set(tt.key, tt.val); err != nil {
			t.Error(err)
		}
	}

	for _, tt := range test {
		switch tt.key {
		case "a":
			var val string
			if err := cache.Get(tt.key, &val); err != nil {
				t.Error(err)
			}
			if ok := reflect.DeepEqual(val, tt.want); !ok {
				t.Errorf("want %v, return val %v", tt.want, val)
			}
		case "b":
			var val int
			if err := cache.Get(tt.key, &val); err != nil {
				t.Error(err)
			}
			if ok := reflect.DeepEqual(val, tt.want); !ok {
				t.Errorf("want %v, return val %v", tt.want, val)
			}
		case "c":
			var val float64
			if err := cache.Get(tt.key, &val); err != nil {
				t.Error(err)
			}
			if ok := reflect.DeepEqual(val, tt.want); !ok {
				t.Errorf("want %v, return val %v", tt.want, val)
			}
		case "d":
			var val []byte
			if err := cache.Get(tt.key, &val); err != nil {
				t.Error(err)
			}
			if ok := reflect.DeepEqual(val, tt.want); !ok {
				t.Errorf("want %v, return val %v", tt.want, val)
			}
		case "e":
			var val TestData
			if err := cache.Get(tt.key, &val); err != nil {
				t.Error(err)
			}
			if ok := reflect.DeepEqual(val, tt.want); !ok {
				t.Errorf("want %v, return val %v", tt.want, val)
			}
		}
	}
}
