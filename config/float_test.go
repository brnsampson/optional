package config_test

import (
	"reflect"
	"strconv"
	"github.com/brnsampson/optional/config"
	"gotest.tools/v3/assert"
	"testing"
)

func TestFloat32Type(t *testing.T) {
	o := config.SomeFloat32(42.0)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestFloat32String(t *testing.T) {
	var f float32 = 42.1
	fStr := strconv.FormatFloat(float64(f), 'g', 3, 32)
	o := config.SomeFloat32(f)
	assert.Equal(t, fStr, o.String())
}

func TestFloat32MarshalText(t *testing.T) {
	var f float32 = 42.1
	fStr := strconv.FormatFloat(float64(f), 'g', 3, 32)
	o := config.SomeFloat32(f)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, fStr, string(s))
}

func TestFloat32UnmarshalText(t *testing.T) {
	var f float32 = 42.1
	fStr := strconv.FormatFloat(float64(f), 'g', 3, 32)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := config.NoFloat32()
	err := o.UnmarshalText([]byte(fStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, f, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-float
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestFloat64Type(t *testing.T) {
	o := config.SomeFloat64(42.0)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestFloat64String(t *testing.T) {
	var f float64 = 42.1
	fStr := strconv.FormatFloat(float64(f), 'g', 3, 64)
	o := config.SomeFloat64(f)
	assert.Equal(t, fStr, o.String())
}

func TestFloat64MarshalText(t *testing.T) {
	var f float64 = 42.1
	fStr := strconv.FormatFloat(float64(f), 'g', 3, 64)
	o := config.SomeFloat64(f)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, fStr, string(s))
}

func TestFloat64UnmarshalText(t *testing.T) {
	var f float64 = 42.1
	fStr := strconv.FormatFloat(float64(f), 'g', 3, 64)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := config.NoFloat64()
	err := o.UnmarshalText([]byte(fStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, f, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-float
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}
