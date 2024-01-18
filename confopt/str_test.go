package confopt_test

import (
	"github.com/brnsampson/optional/confopt"
	"gotest.tools/v3/assert"
	"reflect"
	"testing"
)

func TestStrType(t *testing.T) {
	o := confopt.SomeStr("A dumb test string")
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestStrString(t *testing.T) {
	str := "testing this tester with the testing module"
	o := confopt.SomeStr(str)
	assert.Equal(t, str, o.String())
}

func TestStrMarshalText(t *testing.T) {
	str := "testing this tester with the testing module"
	o := confopt.SomeStr(str)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, str, string(s))
}

func TestStrUnmarshalText(t *testing.T) {
	str := "testing this tester with the testing module"
	nullStr := "null"
	intStr := "42"

	// Text sucessful unmarshaling
	o := confopt.NoStr()
	err := o.UnmarshalText([]byte(str))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, str, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// There is no string that we cannot unmarshal into a string, but we should check that other types actually
	// end up as the string version as expected I guess...
	err = o.UnmarshalText([]byte(intStr))
	assert.NilError(t, err)

	ret, err = o.Get()
	assert.NilError(t, err)
	assert.Equal(t, intStr, ret)
}
