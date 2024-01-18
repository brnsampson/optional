package confopt_test

import (
	"github.com/brnsampson/optional/confopt"
	"gotest.tools/v3/assert"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFileType(t *testing.T) {
	o := confopt.SomeFile("/not/a/real/path")
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestFileGet(t *testing.T) {
	path := "../tls/rsa/cert.pem"
	o := confopt.SomeFile(path)
	abs, err := filepath.Abs(path)
	// an error here doesn't mean our library is broken, just that the path we chose to test with isn't valid.
	assert.NilError(t, err)

	tmp, err := o.Get()
	assert.NilError(t, err)

	assert.Equal(t, path, tmp)
	a, err := o.Abs()
	assert.NilError(t, err)

	tmp, err = a.Get()
	assert.NilError(t, err)
	assert.Equal(t, abs, tmp)
}

func TestFileString(t *testing.T) {
	path := "../tls/rsa/cert.pem"
	noneStr := "None[File]"

	o := confopt.SomeFile(path)
	assert.Equal(t, path, o.String())

	// Test None case displays correctly
	o = confopt.NoFile()
	assert.Equal(t, noneStr, o.String())
}

func TestFileMarshalText(t *testing.T) {
	path := "../tls/rsa/cert.pem"

	o := confopt.SomeFile(path)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, path, string(s))
}

func TestFileUnmarshalText(t *testing.T) {
	path := "../tls/rsa/cert.pem"
	nullFile := "null"
	intFile := "42"

	// Text sucessful unmarshaling
	o := confopt.NoFile()
	err := o.UnmarshalText([]byte(path))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, path, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullFile))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-path. This will work because it should interpret this as a file named "41", which you could
	// theoretically have.
	err = o.UnmarshalText([]byte(intFile))
	assert.NilError(t, err)

	ret, err = o.Get()
	assert.NilError(t, err)
	assert.Equal(t, intFile, ret)
}
