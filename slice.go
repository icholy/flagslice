package flagslice

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// sliceValue is a flag.Value implementation
type sliceValue struct {
	elem  reflect.Type
	slice reflect.Value
	parse parseFn
}

// IsBoolFlag implements an optional interface which allows
// passing bool flags without values.
func (sv sliceValue) IsBoolFlag() bool {
	if _, fv, ok := toFlagValue(sv.elem); ok {
		if bf, ok := fv.(interface{ IsBoolFlag() bool }); ok {
			return bf.IsBoolFlag()
		}
		return false
	}
	return sv.elem.Kind() == reflect.Bool
}

// String implements flag.Value
func (sv sliceValue) String() string {
	if !sv.slice.IsValid() {
		return ""
	}
	ss := make([]string, sv.slice.Len())
	for i := 0; i < sv.slice.Len(); i++ {
		ss[i] = fmt.Sprint(sv.slice.Index(i))
	}
	return strings.Join(ss, ", ")
}

// Set implements flag.Value
func (sv sliceValue) Set(s string) error {
	v, err := sv.parse(s)
	if err != nil {
		return err
	}
	sv.slice.Set(reflect.Append(sv.slice, reflect.ValueOf(v)))
	return nil
}

// Value accepts a pointer to slice and returns a flag.Value which appends
// a value to the slice for each call to Set. This function will panic if
// the slice parameter is invalid.
func Value(slice interface{}) flag.Value {
	p := reflect.ValueOf(slice)
	if p.Kind() != reflect.Ptr || p.Elem().Kind() != reflect.Slice {
		panic(fmt.Sprintf("expected pointer to slice, got %s", p.Type()))
	}
	s := p.Elem()
	et := s.Type().Elem()
	parse, ok := toParseFn(et)
	if !ok {
		panic(fmt.Sprintf("unsupported slice type %s", s.Type()))
	}
	return sliceValue{slice: s, parse: parse, elem: et}
}

// conv parses an argument string
type parseFn func(s string) (interface{}, error)

// toParseFn returns a parsing function for the provided type t.
func toParseFn(t reflect.Type) (parseFn, bool) {
	// check if the element type implements flag.Value
	if _, _, ok := toFlagValue(t); ok {
		return func(s string) (interface{}, error) {
			v, fv, _ := toFlagValue(t)
			err := fv.Set(s)
			return v.Interface(), err
		}, true
	}
	// special case time.Duration
	if t == reflect.TypeOf(time.Duration(0)) {
		return func(s string) (interface{}, error) {
			return time.ParseDuration(s)
		}, true
	}
	// check for supported underlying types
	switch t.Kind() {
	case reflect.Bool:
		return func(s string) (interface{}, error) {
			return strconv.ParseBool(s)
		}, true
	case reflect.Float64:
		return func(s string) (interface{}, error) {
			return strconv.ParseFloat(s, 64)
		}, true
	case reflect.Int:
		return func(s string) (interface{}, error) {
			x, err := strconv.ParseInt(s, 0, strconv.IntSize)
			return int(x), err
		}, true
	case reflect.Int64:
		return func(s string) (interface{}, error) {
			return strconv.ParseInt(s, 0, 64)
		}, true
	case reflect.Uint:
		return func(s string) (interface{}, error) {
			x, err := strconv.ParseUint(s, 0, strconv.IntSize)
			return uint(x), err
		}, true
	case reflect.Uint64:
		return func(s string) (interface{}, error) {
			return strconv.ParseUint(s, 0, 64)
		}, true
	case reflect.String:
		return func(s string) (interface{}, error) {
			return s, nil
		}, true
	default:
		return nil, false
	}
}

// toFlagValue creates a flag.Value if t implements flag.Value
func toFlagValue(t reflect.Type) (reflect.Value, flag.Value, bool) {
	fvt := reflect.TypeOf((*flag.Value)(nil)).Elem()
	if t.Implements(fvt) {
		var v reflect.Value
		if t.Kind() == reflect.Ptr {
			v = reflect.New(t.Elem())
		} else {
			v = reflect.Zero(t)
		}
		return v, v.Interface().(flag.Value), true
	}
	if reflect.PtrTo(t).Implements(fvt) {
		v := reflect.New(t)
		return v.Elem(), v.Interface().(flag.Value), true
	}
	return reflect.Value{}, nil, false
}
