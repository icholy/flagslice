package flagslice

import (
	"reflect"
	"testing"
	"time"
)

type custom struct {
	Value string
}

func (c custom) String() string { return c.Value }

func (c *custom) Set(s string) error {
	c.Value = s
	return nil
}

func TestValue(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		expect interface{}
	}{
		{
			name:   "int",
			args:   []string{"1", "42", "999"},
			expect: []int{1, 42, 999},
		},
		{
			name:   "int64",
			args:   []string{"1", "42", "999"},
			expect: []int64{1, 42, 999},
		},
		{
			name:   "uint",
			args:   []string{"1", "42", "999"},
			expect: []uint{1, 42, 999},
		},
		{
			name:   "uint64",
			args:   []string{"1", "42", "999"},
			expect: []uint64{1, 42, 999},
		},
		{
			name:   "uint64",
			args:   []string{"1", "42", "999"},
			expect: []float64{1, 42, 999},
		},
		{
			name:   "string",
			args:   []string{"foo", "bar", "baz"},
			expect: []string{"foo", "bar", "baz"},
		},
		{
			name:   "duration",
			args:   []string{"2s", "10m", "0"},
			expect: []time.Duration{2 * time.Second, 10 * time.Minute, 0},
		},
		{
			name:   "custom pointers",
			args:   []string{"a", "b"},
			expect: []*custom{{"a"}, {"b"}},
		},
		{
			name:   "custom values",
			args:   []string{"a", "b"},
			expect: []custom{{"a"}, {"b"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := reflect.MakeSlice(reflect.TypeOf(tt.expect), 0, 0)
			ptr := reflect.New(slice.Type())
			ptr.Elem().Set(slice)
			v := Value(ptr.Interface())
			for _, a := range tt.args {
				if err := v.Set(a); err != nil {
					t.Errorf("arg %q: %v", a, err)
				}
			}
			if s := ptr.Elem().Interface(); !reflect.DeepEqual(s, tt.expect) {
				t.Errorf("expected %v, got %v", tt.expect, s)
			}
		})
	}
}

func TestPanicOnInvalidSlice(t *testing.T) {
	t.Run("non-pointer", func(t *testing.T) {
		defer func() { recover() }()
		Value([]string{})
		t.Error("did not panic")
	})
	t.Run("non-slice", func(t *testing.T) {
		defer func() { recover() }()
		Value(new(int))
		t.Error("did not panic")
	})
	t.Run("non-supported", func(t *testing.T) {
		defer func() { recover() }()
		var floats []float32
		Value(&floats)
		t.Error("did not panic")
	})
}
