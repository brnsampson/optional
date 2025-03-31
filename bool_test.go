package optional_test

import (
	"reflect"
	"testing"

	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
)

func TestBoolType(t *testing.T) {
	o := optional.SomeBool(true)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestBoolString(t *testing.T) {
	trueString := "true"
	o := optional.SomeBool(true)
	assert.Equal(t, trueString, o.String())
}

func TestBoolMarshalText(t *testing.T) {
	trueString := "true"
	o := optional.SomeBool(true)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, trueString, string(s))
}

func TestBoolUnmarshalText(t *testing.T) {
	trueString := "true"
	nullString := "null"
	intString := "42"

	// Text sucessful unmarshaling
	o := optional.NoBool()
	err := o.UnmarshalText([]byte(trueString))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, true, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullString))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-bool
	err = o.UnmarshalText([]byte(intString))
	assert.Assert(t, err != nil)
}

func TestBoolTrue(t *testing.T) {
	o := optional.NoBool()
	assert.Assert(t, !o.True())

	o = optional.SomeBool(false)
	assert.Assert(t, !o.True())

	o = optional.SomeBool(true)
	assert.Assert(t, o.True())
}
