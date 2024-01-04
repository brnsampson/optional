package config_test

import (
	"strconv"
	"github.com/brnsampson/optional/config"
	"gotest.tools/v3/assert"
	"testing"
)

func TestIntString(t *testing.T) {
	var i int = 42
	iStr := strconv.Itoa(i)
	o := config.SomeInt(i)
	assert.Equal(t, iStr, o.String())
}

func TestIntMarshalText(t *testing.T) {
	var i int = 42
	iStr := strconv.Itoa(i)
	o := config.SomeInt(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestIntUnmarshalText(t *testing.T) {
	var i int = 42
	iStr := strconv.Itoa(i)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := config.NoInt()
	err := o.UnmarshalText([]byte(iStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, i, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-int
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestInt8String(t *testing.T) {
	var i int8 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := config.SomeInt8(i)
	assert.Equal(t, iStr, o.String())
}

func TestInt8MarshalText(t *testing.T) {
	var i int8 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := config.SomeInt8(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestInt8UnmarshalText(t *testing.T) {
	var i int8 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := config.NoInt8()
	err := o.UnmarshalText([]byte(iStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, i, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-int
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestInt16String(t *testing.T) {
	var i int16 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := config.SomeInt16(i)
	assert.Equal(t, iStr, o.String())
}

func TestInt16MarshalText(t *testing.T) {
	var i int16 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := config.SomeInt16(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestInt16UnmarshalText(t *testing.T) {
	var i int16 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := config.NoInt16()
	err := o.UnmarshalText([]byte(iStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, i, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-int
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestInt32String(t *testing.T) {
	var i int32 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := config.SomeInt32(i)
	assert.Equal(t, iStr, o.String())
}

func TestInt32MarshalText(t *testing.T) {
	var i int32 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := config.SomeInt32(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestInt32UnmarshalText(t *testing.T) {
	var i int32 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := config.NoInt32()
	err := o.UnmarshalText([]byte(iStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, i, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-int
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestInt64String(t *testing.T) {
	var i int64 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := config.SomeInt64(i)
	assert.Equal(t, iStr, o.String())
}

func TestInt64MarshalText(t *testing.T) {
	var i int64 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := config.SomeInt64(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestInt64UnmarshalText(t *testing.T) {
	var i int64 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := config.NoInt64()
	err := o.UnmarshalText([]byte(iStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, i, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-int
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}
