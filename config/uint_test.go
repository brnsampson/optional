package config_test

import (
	"strconv"
	"github.com/brnsampson/optional/config"
	"gotest.tools/v3/assert"
	"testing"
)

func TestUintString(t *testing.T) {
	var i uint = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := config.SomeUint(i)
	assert.Equal(t, iStr, o.String())
}

func TestUintMarshalText(t *testing.T) {
	var i uint = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := config.SomeUint(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestUintUnmarshalText(t *testing.T) {
	var i uint = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := config.NoUint()
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

	// Test unmarshaling non-uint
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestUint8String(t *testing.T) {
	var i uint8 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := config.SomeUint8(i)
	assert.Equal(t, iStr, o.String())
}

func TestUint8MarshalText(t *testing.T) {
	var i uint8 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := config.SomeUint8(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestUint8UnmarshalText(t *testing.T) {
	var i uint8 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := config.NoUint8()
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

	// Test unmarshaling non-uint
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestUint16String(t *testing.T) {
	var i uint16 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := config.SomeUint16(i)
	assert.Equal(t, iStr, o.String())
}

func TestUint16MarshalText(t *testing.T) {
	var i uint16 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := config.SomeUint16(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestUint16UnmarshalText(t *testing.T) {
	var i uint16 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := config.NoUint16()
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

	// Test unmarshaling non-uint
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestUint32String(t *testing.T) {
	var f uint32 = 42
	fStr := strconv.FormatUint(uint64(f), 10)
	o := config.SomeUint32(f)
	assert.Equal(t, fStr, o.String())
}

func TestUint32MarshalText(t *testing.T) {
	var f uint32 = 42
	fStr := strconv.FormatUint(uint64(f), 10)
	o := config.SomeUint32(f)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, fStr, string(s))
}

func TestUint32UnmarshalText(t *testing.T) {
	var i uint32 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := config.NoUint32()
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

	// Test unmarshaling non-uint
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestUint64String(t *testing.T) {
	var i uint64 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := config.SomeUint64(i)
	assert.Equal(t, iStr, o.String())
}

func TestUint64MarshalText(t *testing.T) {
	var i uint64 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := config.SomeUint64(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestUint64UnmarshalText(t *testing.T) {
	var i uint64 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := config.NoUint64()
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

	// Test unmarshaling non-uint
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}
