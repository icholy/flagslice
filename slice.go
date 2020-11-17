package flagslice

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type sliceValue struct {
	slice reflect.Value
	set   func(s string) (interface{}, error)
}

func (sv sliceValue) String() string {
	ss := make([]string, sv.slice.Len())
	for i := 0; i < sv.slice.Len(); i++ {
		ss[i] = fmt.Sprint(sv.slice.Index(i))
	}
	return strings.Join(ss, ", ")
}

func (sv sliceValue) Set(s string) error {
	v, err := sv.set(s)
	if err != nil {
		return err
	}
	sv.slice.Set(reflect.Append(sv.slice, reflect.ValueOf(v)))
	return nil
}

// Value accepts a pointer to slice and returns a flag.Value which appends
// a value to the slice for each call to Set.
func Value(slice interface{}) flag.Value {
	p := reflect.ValueOf(slice)
	if p.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("expected pointer to slice, got %s", p.Type()))
	}
	s := p.Elem()
	if s.Kind() != reflect.Slice {
		panic(fmt.Sprintf("expected pointer to slice, got %s", p.Type()))
	}
	if s.Type().Elem() == reflect.TypeOf(time.Duration(0)) {
		return sliceValue{slice: s, set: func(s string) (interface{}, error) {
			return time.ParseDuration(s)
		}}
	}
	switch s.Type().Elem().Kind() {
	case reflect.Bool:
		return sliceValue{slice: s, set: func(s string) (interface{}, error) {
			return strconv.ParseBool(s)
		}}
	case reflect.Float64:
		return sliceValue{slice: s, set: func(s string) (interface{}, error) {
			return strconv.ParseFloat(s, 64)
		}}
	case reflect.Int:
		return sliceValue{slice: s, set: func(s string) (interface{}, error) {
			return strconv.Atoi(s)
		}}
	case reflect.Int64:
		return sliceValue{slice: s, set: func(s string) (interface{}, error) {
			return strconv.ParseInt(s, 10, 64)
		}}
	case reflect.Uint:
		return sliceValue{slice: s, set: func(s string) (interface{}, error) {
			x, err := strconv.ParseUint(s, 10, 64)
			return uint(x), err
		}}
	case reflect.Uint64:
		return sliceValue{slice: s, set: func(s string) (interface{}, error) {
			return strconv.ParseUint(s, 10, 64)
		}}
	case reflect.String:
		return sliceValue{slice: s, set: func(s string) (interface{}, error) {
			return s, nil
		}}
	default:
		panic(fmt.Sprintf("unsupported slice type %s", s.Type()))
	}
}
