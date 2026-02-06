package optional_test

import (
	"reflect"
	"testing"

	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
)

func TestByteType(t *testing.T) {
	o := optional.SomeByte('r')
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestByteString(t *testing.T) {
	byteString := "6e"
	o := optional.SomeByte('n')
	assert.Equal(t, byteString, o.String())
}

func TestByteMarshalText(t *testing.T) {
	byteString := "6e"
	o := optional.SomeByte('n')

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, byteString, string(s))
}

func TestByteUnmarshalText(t *testing.T) {
	byteString := "6e"
	nullString := "null"
	longString := "this is more than one byte"

	// Text sucessful unmarshaling
	o := optional.NoByte()
	err := o.UnmarshalText([]byte(byteString))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, byte('n'), ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullString))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-byte
	err = o.UnmarshalText([]byte(longString))
	assert.Assert(t, err != nil)
}
