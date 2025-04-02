package optional_test

import (
	"reflect"
	"testing"

	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
)

func TestSecretType(t *testing.T) {
	o := optional.SomeStr("A dumb test string")
	s := optional.MakeSecret(&o)
	assert.Equal(t, reflect.TypeOf(s).Name(), s.Type())
	assert.Assert(t, o.IsNone())
}

func TestSecretString(t *testing.T) {
	str := "***REDACTED***"
	o := optional.SomeStr(str)
	s := optional.MakeSecret(&o)
	assert.Equal(t, str, s.String())
	assert.Assert(t, o.IsNone())
}

func TestSecretMarshalText(t *testing.T) {
	str := "testing this tester with the testing module"
	o := optional.SomeSecret(str)
	s, err := o.MarshalText()

	assert.NilError(t, err)
	assert.Equal(t, str, string(s))
}

func TestSecretUnmarshalText(t *testing.T) {
	str := "testing this tester with the testing module"
	nullStr := "null"
	intStr := "42"

	// Text sucessful unmarshaling
	o := optional.NoSecret()
	err := o.UnmarshalText([]byte(str))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, str, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// There is no string that we cannot unmarshal into a string, but we should check that other types actually
	// end up as the string version as expected I guess...
	err = o.UnmarshalText([]byte(intStr))
	assert.NilError(t, err)

	ret, ok = o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, intStr, ret)
}
