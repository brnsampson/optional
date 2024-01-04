package config_test

import (
	"github.com/brnsampson/optional/config"
	"gotest.tools/v3/assert"
	"testing"
)

func TestBoolString(t *testing.T) {
	trueString := "true"
	o := config.SomeBool(true)
	assert.Equal(t, trueString, o.String())
}

func TestBoolMarshalText(t *testing.T) {
	trueString := "true"
	o := config.SomeBool(true)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, trueString, string(s))
}

func TestBoolUnmarshalText(t *testing.T) {
	trueString := "true"
	nullString := "null"
	intString := "42"

	// Text sucessful unmarshaling
	o := config.NoBool()
	err := o.UnmarshalText([]byte(trueString))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, true, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullString))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-bool
	err = o.UnmarshalText([]byte(intString))
	assert.Assert(t, err != nil)
}
