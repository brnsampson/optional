package config_test

import (
	"encoding/json"
	"github.com/brnsampson/optional/config"
	"gotest.tools/v3/assert"
	"reflect"
	"slices"
	"testing"
	"time"
)

func TestTimeType(t *testing.T) {
	o := config.SomeTime(time.Now())
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestTimeString(t *testing.T) {
	now := time.Now().Truncate(0)
	nowString := now.Format(time.RFC3339Nano)
	o := config.SomeTime(now).WithFormats(time.RFC3339Nano)
	assert.Equal(t, nowString, o.String())
}

func TestTimeMarshalText(t *testing.T) {
	now := time.Now().Truncate(0)
	nowString := now.Format(time.RFC3339Nano)
	o := config.SomeTime(now).WithFormats(time.RFC3339Nano)
	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, nowString, string(s))
}

func TestTimeUnmarshalText(t *testing.T) {
	now := time.Now().Truncate(0)
	nowString := now.Format(time.RFC3339Nano)
	wait := 10 * time.Second
	later := now.Add(wait)

	// Text sucessful unmarshaling
	o := config.NoTime().WithFormats(time.RFC3339Nano)
	err := o.UnmarshalText([]byte(nowString))
	assert.NilError(t, err)

	ret, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, now, ret)
	assert.Equal(t, later, ret.Add(wait))

	// Test unmarshaling null
	assert.Assert(t, o.IsSome())
	err = o.UnmarshalText([]byte("null"))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-date
	err = o.UnmarshalText([]byte("this is not a date!"))
	assert.Assert(t, err != nil)

	// Test unmarshaling different format
	o2 := o.WithFormats(time.UnixDate)
	err = o2.UnmarshalText([]byte(nowString))
	assert.Assert(t, err != nil)
}

func TestTimeMarshalJson(t *testing.T) {
	now := time.Now().Truncate(0)
	nowString := now.Format(time.RFC3339Nano)
	nowJson := "\"" + nowString + "\""

	o := config.SomeTime(now).WithFormats(time.RFC3339Nano)
	res, err := json.Marshal(o)
	assert.NilError(t, err)
	assert.Equal(t, nowJson, string(res))
}

func TestTimeUnmarshalJson(t *testing.T) {
	now := time.Now().Truncate(0)
	nowString := now.Format(time.RFC3339Nano)
	nowJson := "\"" + nowString + "\""
	nullJson := []byte(`null`)

	//wait := 10 * time.Second
	//later := now.Add(wait)
	nowUnixString := now.Format(time.UnixDate)
	nowUnixJson := "\"" + nowUnixString + "\""

	// Text null case
	var n config.Time
	expected := config.NoTime()
	expectedFormats := expected.Formats()
	json.Unmarshal(nullJson, &n)

	// Need to test == instead of expected.Eq(&n) because there is some additional state in the Time struct in addition
	// to the embedded Option.
	assert.Assert(t, n.IsNone())
	assert.Assert(t, expected.Eq(&n))
	formats := n.Formats()
	for _, f := range expectedFormats {
		assert.Assert(t, slices.Contains(formats, f))
	}

	// Test valid case
	var o config.Time
	err := json.Unmarshal([]byte(nowJson), &o)
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	res, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, now, res)

	// Test invalid data case
	var p config.Time
	err = json.Unmarshal([]byte(nowUnixJson), &p)
	assert.Assert(t, err != nil)

	var q config.Time
	err = json.Unmarshal([]byte("this is not a date"), &q)
	assert.Assert(t, err != nil)
}
