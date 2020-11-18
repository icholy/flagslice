package flagslice

import (
	"flag"
	"reflect"
	"testing"
	"time"
)

type custom struct {
	Value string
}

func (c custom) String() string { return c.Value }

func (c custom) IsBoolFlag() bool { return false }

func (c *custom) Set(s string) error {
	c.Value = s
	return nil
}

type custom2 struct{}

func (c custom2) String() string     { return "" }
func (c custom2) Set(s string) error { return nil }

func TestValue(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		expect interface{}
	}{
		{
			name:   "bool",
			args:   []string{"true", "false"},
			expect: []bool{true, false},
		},
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
		{
			name:   "custom2 pointers",
			args:   []string{"a", "b"},
			expect: []*custom2{{}, {}},
		},
		{
			name:   "custom2 values",
			args:   []string{"a", "b"},
			expect: []custom2{{}, {}},
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

func TestFlagSet(t *testing.T) {
	fset := flag.NewFlagSet("", flag.ContinueOnError)
	var (
		bools    = []bool{}
		strings  = []string{}
		customs  = []custom{}
		customs2 = []custom2{}
	)
	fset.Var(Value(&strings), "s", "string")
	fset.Var(Value(&bools), "b", "bool")
	fset.Var(Value(&customs), "c", "custom")
	fset.Var(Value(&customs2), "c2", "custom2")
	if err := fset.Parse([]string{
		"-s", "foo", "-s", "bar",
		"-b", "-b=false",
		"-c", "thing",
		"-c2", "thing2",
	}); err != nil {
		t.Fatal(err)
	}
	fset.VisitAll(func(f *flag.Flag) { _ = f.Value.String() })
	if expect := []bool{true, false}; !reflect.DeepEqual(bools, expect) {
		t.Errorf("expected %v, got %v", expect, bools)
	}
	if expect := []string{"foo", "bar"}; !reflect.DeepEqual(strings, expect) {
		t.Errorf("expected %v, got %v", expect, strings)
	}
}
