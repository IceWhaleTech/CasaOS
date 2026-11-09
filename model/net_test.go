package model

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/shirou/gopsutil/v3/net"
)

// The counters used to be reinterpreted with unsafe.Pointer, which read past the
// end of gopsutil's smaller struct. Copying them by name keeps every field
// gopsutil has, and leaves the two fields it does not have alone.
func TestIOCountersFromCopiesEveryFieldGopsutilHas(t *testing.T) {
	source := net.IOCountersStat{}
	fields := reflect.ValueOf(&source).Elem()
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		switch field.Kind() {
		case reflect.String:
			field.SetString(fmt.Sprintf("value of %s", fields.Type().Field(i).Name))
		case reflect.Uint64:
			field.SetUint(uint64(i) + 1)
		default:
			t.Fatalf("%s is a %s, which this test does not know how to fill", fields.Type().Field(i).Name, field.Kind())
		}
	}

	copied := reflect.ValueOf(IOCountersFrom(source))
	for i := 0; i < fields.NumField(); i++ {
		name := fields.Type().Field(i).Name
		got := copied.FieldByName(name)
		if !got.IsValid() {
			t.Fatalf("%s is not copied", name)
		}
		if got.Interface() != fields.Field(i).Interface() {
			t.Errorf("%s: got %v, want %v", name, got.Interface(), fields.Field(i).Interface())
		}
	}

	own := IOCountersFrom(net.IOCountersStat{})
	if own.State != "" || own.Time != 0 {
		t.Errorf("the fields gopsutil does not have must stay empty: state %q, time %d", own.State, own.Time)
	}
}
