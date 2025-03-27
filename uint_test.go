package optional_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
)

func TestUintType(t *testing.T) {
	o := optional.SomeUint(42)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestUintString(t *testing.T) {
	var i uint = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint(i)
	assert.Equal(t, iStr, o.String())
}

func TestUintMarshalText(t *testing.T) {
	var i uint = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint(i)

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
	o := optional.NoUint()
	err := o.UnmarshalText([]byte(iStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, i, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-uint
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestUint8Type(t *testing.T) {
	o := optional.SomeUint8(42)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestUint8String(t *testing.T) {
	var i uint8 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint8(i)
	assert.Equal(t, iStr, o.String())
}

func TestUint8MarshalText(t *testing.T) {
	var i uint8 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint8(i)

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
	o := optional.NoUint8()
	err := o.UnmarshalText([]byte(iStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, i, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-uint
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestUint16Type(t *testing.T) {
	o := optional.SomeUint16(42)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestUint16String(t *testing.T) {
	var i uint16 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint16(i)
	assert.Equal(t, iStr, o.String())
}

func TestUint16MarshalText(t *testing.T) {
	var i uint16 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint16(i)

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
	o := optional.NoUint16()
	err := o.UnmarshalText([]byte(iStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, i, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-uint
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestUint32Type(t *testing.T) {
	o := optional.SomeUint32(42)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestUint32String(t *testing.T) {
	var f uint32 = 42
	fStr := strconv.FormatUint(uint64(f), 10)
	o := optional.SomeUint32(f)
	assert.Equal(t, fStr, o.String())
}

func TestUint32MarshalText(t *testing.T) {
	var f uint32 = 42
	fStr := strconv.FormatUint(uint64(f), 10)
	o := optional.SomeUint32(f)

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
	o := optional.NoUint32()
	err := o.UnmarshalText([]byte(iStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, i, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-uint
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestUint64Type(t *testing.T) {
	o := optional.SomeUint64(42)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestUint64String(t *testing.T) {
	var i uint64 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint64(i)
	assert.Equal(t, iStr, o.String())
}

func TestUint64MarshalText(t *testing.T) {
	var i uint64 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint64(i)

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
	o := optional.NoUint64()
	err := o.UnmarshalText([]byte(iStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, i, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-uint
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}
