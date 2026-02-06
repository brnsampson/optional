package optional_test

import (
	"encoding/json"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
)

func TestTimeType(t *testing.T) {
	o := optional.SomeTime(time.Now())
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestTimeString(t *testing.T) {
	now := time.Now().Truncate(0)
	nowString := now.Format(time.RFC3339Nano)
	o := optional.SomeTime(now).WithFormats(time.RFC3339Nano)
	assert.Equal(t, nowString, o.String())
}

func TestTimeMarshalText(t *testing.T) {
	now := time.Now().Truncate(0)
	nowString := now.Format(time.RFC3339Nano)
	o := optional.SomeTime(now).WithFormats(time.RFC3339Nano)
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
	o := optional.NoTime().WithFormats(time.RFC3339Nano)
	err := o.UnmarshalText([]byte(nowString))
	assert.NilError(t, err)

	ret, ok := o.Get()
	assert.Assert(t, ok)
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

	o := optional.SomeTime(now).WithFormats(time.RFC3339Nano)
	res, err := json.Marshal(o)
	assert.NilError(t, err)
	assert.Equal(t, nowJson, string(res))
}

func TestTimeUnmarshalJson(t *testing.T) {
	now := time.Now().Truncate(0)
	nowString := now.Format(time.RFC3339Nano)
	nowJson := "\"" + nowString + "\""
	nullJson := []byte(`null`)

	// wait := 10 * time.Second
	// later := now.Add(wait)
	nowUnixString := now.Format(time.UnixDate)
	nowUnixJson := "\"" + nowUnixString + "\""

	// Text null case
	var n optional.Time
	expected := optional.NoTime()
	expectedFormats := expected.Formats()
	json.Unmarshal(nullJson, &n)

	// Need to test both optional.Equal(expected, n) and that the format slice is the same.
	assert.Assert(t, n.IsNone())
	assert.Assert(t, optional.Equal(expected, n))
	formats := n.Formats()
	for _, f := range expectedFormats {
		assert.Assert(t, slices.Contains(formats, f))
	}

	// Test valid case
	var o optional.Time
	err := json.Unmarshal([]byte(nowJson), &o)
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	res, ok := o.Get()
	assert.Assert(t, ok)
	assert.NilError(t, err)
	assert.Equal(t, now, res)

	// Test invalid data case
	var p optional.Time
	err = json.Unmarshal([]byte(nowUnixJson), &p)
	assert.NilError(t, err)

	var q optional.Time
	err = json.Unmarshal([]byte("this is not a date"), &q)
	assert.Assert(t, err != nil)
}

func TestDurationType(t *testing.T) {
	d, err := time.ParseDuration("300ms")
	assert.NilError(t, err)
	o := optional.SomeDuration(d)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestDurationString(t *testing.T) {
	d, err := time.ParseDuration("300ms")
	assert.NilError(t, err)
	dStr := d.String()
	o := optional.SomeDuration(d)
	assert.Equal(t, dStr, o.String())
}

func TestDurationMarshalText(t *testing.T) {
	d, err := time.ParseDuration("300ms")
	assert.NilError(t, err)
	dStr := d.String()
	o := optional.SomeDuration(d)
	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, dStr, string(s))
}

func TestDurationUnmarshalText(t *testing.T) {
	dStr := "300ms"
	d, err := time.ParseDuration(dStr)
	assert.NilError(t, err)

	wait := 10 * time.Second
	later := d + wait

	// Text sucessful unmarshaling
	o := optional.NoDuration()
	err = o.UnmarshalText([]byte(dStr))
	assert.NilError(t, err)

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, d, ret)
	assert.Equal(t, later, ret+wait)

	// Test unmarshaling null
	assert.Assert(t, o.IsSome())
	err = o.UnmarshalText([]byte("null"))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-date
	err = o.UnmarshalText([]byte("this is not a date!"))
	assert.Assert(t, err != nil)
}

func TestDurationMarshalJson(t *testing.T) {
	d, err := time.ParseDuration("300ms")
	assert.NilError(t, err)
	dStr := d.String()

	dJson := "\"" + dStr + "\""

	o := optional.SomeDuration(d)
	res, err := json.Marshal(o)
	assert.NilError(t, err)
	assert.Equal(t, dJson, string(res))
}

func TestDurationUnmarshalJson(t *testing.T) {
	dStr := "300ms"
	d, err := time.ParseDuration(dStr)
	assert.NilError(t, err)

	dJson := "\"" + dStr + "\""
	nullJson := []byte(`null`)

	// Text null case
	var n optional.Duration
	expected := optional.NoDuration()
	json.Unmarshal(nullJson, &n)

	// Need to test both optional.Equal(expected, n) and that the format slice is the same.
	assert.Assert(t, n.IsNone())
	assert.Assert(t, optional.Equal(expected, n))

	// Test valid case
	var o optional.Duration
	err = json.Unmarshal([]byte(dJson), &o)
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	res, ok := o.Get()
	assert.Assert(t, ok)
	assert.NilError(t, err)
	assert.Equal(t, d, res)

	// Test invalid data case
	var q optional.Duration
	err = json.Unmarshal([]byte("this is not a duration"), &q)
	assert.Assert(t, err != nil)
}
