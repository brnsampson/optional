package optional_test

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"testing"

	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
)

func TestOptionZeroValueIsNone(t *testing.T) {
	var o optional.Option[int]
	assert.Assert(t, o.IsNone())
	assert.Assert(t, !o.IsSome())
}

func TestOptionPointerIsMutableOptional(t *testing.T) {
	// The real test is if we get compiler errors because &o does not implement Optional
	o := optional.Some(42)
	var op optional.MutableOptional[int] = &o
	assert.Assert(t, op.IsSome())
}

func TestOptionalClone(t *testing.T) {
	o := optional.Some(42)
	assert.Assert(t, o.IsSome())
	clone := o.Clone()
	o.Replace(49)
	assert.Assert(t, clone.IsSome())
	assert.Equal(t, 42, clone.MustGet())
}

func TestMutableOptionalMutableClone(t *testing.T) {
	o := optional.Some(42)
	var op optional.MutableOptional[int] = &o
	clone := op.MutableClone()
	clone.Replace(49)
	assert.Assert(t, op.IsSome())
	assert.Equal(t, 42, op.MustGet())
}

func TestOptionBasics(t *testing.T) {
	// Covers IsSome, IsNone, Clear, Default, Replace, and Get
	val := 42
	val2 := 49
	val3 := 66

	o := optional.Some(val)
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	tmp, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, val, tmp)
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	o.Clear()
	assert.Assert(t, !o.IsSome())
	assert.Assert(t, o.IsNone())

	tmp, ok = o.Get()
	assert.Assert(t, !ok)

	o.Default(val2)
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	replaced := o.Replace(val3)
	assert.Assert(t, replaced.Match(val2))
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	tmp, ok = o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, val3, tmp)
}

func TestOptionMatch(t *testing.T) {
	val := 42
	val2 := 49
	none := optional.None[int]()
	o1 := optional.Some(val)
	o2 := optional.Some(val2)
	assert.Assert(t, o1.Match(val))
	assert.Assert(t, !o2.Match(val))
	assert.Assert(t, !none.Match(val))
}

func TestOptionEquality(t *testing.T) {
	val := 42
	val2 := 49
	none := optional.None[int]()
	none2 := optional.None[int]()
	o1 := optional.Some(val)
	o2 := optional.Some(val)
	o3 := optional.Some(val2)
	assert.Equal(t, o1, o2)
	assert.Equal(t, none, none2)
	assert.Assert(t, o1 != o3)
	assert.Assert(t, o2 != none)
}

func TestOptionTransform(t *testing.T) {
	val := 42
	err_string := "Error"
	transform := func(x int) (int, error) { return x + 7, nil }
	err_transform := func(x int) (int, error) { return x + 9, fmt.Errorf(err_string) }
	after_val, _ := transform(val)

	o := optional.Some(val)
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	// Transforming a None value results in a None
	none := optional.None[int]()
	err := none.Transform(transform)
	assert.NilError(t, err)
	assert.Assert(t, none.IsNone())

	// Transforming eith a function that doesn't err is the same as a normal transform
	err = o.Transform(transform)
	assert.NilError(t, err)

	tmp, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, after_val, tmp)

	// transforming with a function that produces an error should return the same error
	err = o.Transform(err_transform)
	assert.Error(t, err, err_string)

	// The error should prevent the transform from actually being applied
	tmp, ok = o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, after_val, tmp)
}

func TestOptionMarshalJSONNumber(t *testing.T) {
	null := "null"
	tmp := optional.None[int]()

	res, err := json.Marshal(tmp)
	assert.NilError(t, err)
	assert.Equal(t, null, string(res))

	intTest := 42
	expected := strconv.Itoa(intTest)
	tmp = optional.Some(intTest)
	res, err = json.Marshal(tmp)
	assert.NilError(t, err)
	assert.Equal(t, expected, string(res))

	var uintTest uint64 = 42
	expected = strconv.FormatUint(uintTest, 10)
	tmp2 := optional.Some(uintTest)
	res, err = json.Marshal(tmp2)
	assert.NilError(t, err)
	assert.Equal(t, expected, string(res))

	floatTest := 42.1
	expected = strconv.FormatFloat(floatTest, 'f', 1, 64)
	tmp3 := optional.Some(floatTest)
	res, err = json.Marshal(tmp3)
	assert.NilError(t, err)
	assert.Equal(t, expected, string(res))
}

func TestOptionMarshalJSONString(t *testing.T) {
	null := "null"
	tmp := optional.None[string]()

	res, err := json.Marshal(tmp)
	assert.NilError(t, err)
	assert.Equal(t, null, string(res))

	test := "this is a test in an option"
	expected := "\"" + test + "\""
	tmp = optional.Some(test)
	res, err = json.Marshal(tmp)
	assert.NilError(t, err)
	assert.Equal(t, expected, string(res))
}

func TestOptionMarshalJSONBool(t *testing.T) {
	null := "null"
	tmp := optional.None[bool]()

	res, err := json.Marshal(tmp)
	assert.NilError(t, err)
	assert.Equal(t, null, string(res))

	test := true
	expected := strconv.FormatBool(test)
	tmp = optional.Some(test)
	res, err = json.Marshal(tmp)
	assert.NilError(t, err)
	assert.Equal(t, expected, string(res))
}

func TestOptionMarshalJSONArray(t *testing.T) {
	null := "null"
	var tmp []optional.Option[int]

	res, err := json.Marshal(tmp)
	assert.NilError(t, err)
	assert.Equal(t, null, string(res))

	tmp = []optional.Option[int]{optional.Some(42), optional.Some(49), optional.None[int]()}
	expected := string([]byte(`[42,49,null]`))
	res, err = json.Marshal(tmp)
	assert.NilError(t, err)
	assert.Equal(t, expected, string(res))
}

func TestOptionMarshalJSONObject(t *testing.T) {
	null := "null"
	var tmp map[string]optional.Option[int]

	res, err := json.Marshal(tmp)
	assert.NilError(t, err)
	assert.Equal(t, null, string(res))

	tmp = map[string]optional.Option[int]{"one": optional.Some(42), "two": optional.Some(49), "three": optional.None[int]()}
	expected := string([]byte(`{"one":42,"three":null,"two":49}`))
	res, err = json.Marshal(tmp)
	assert.NilError(t, err)
	assert.Equal(t, expected, string(res))
}

func TestOptionMarshalJSONStruct(t *testing.T) {
	rawJSON := []byte(`{"A":42,"B":null,"C":"testing the tester","D":null,"E":{"Animal":"monkey"},"F":null}`)
	type MyStruct struct {
		Animal string
	}

	test1 := struct {
		A optional.Option[int]
		B optional.Option[int]
		C optional.Option[string]
		D optional.Option[string]
		E optional.Option[MyStruct]
		F optional.Option[MyStruct]
	}{
		optional.Some(42),
		optional.None[int](),
		optional.Some("testing the tester"),
		optional.None[string](),
		optional.Some(MyStruct{"monkey"}),
		optional.None[MyStruct](),
	}

	res, err := json.Marshal(test1)
	assert.NilError(t, err)
	assert.Equal(t, string(rawJSON), string(res))
}

func TestOptionUnMarshalJSONNumber(t *testing.T) {
	var n optional.Option[int]
	expected := optional.None[int]()
	nullTest := []byte("null")
	err := json.Unmarshal(nullTest, &n)
	assert.NilError(t, err)
	assert.Assert(t, optional.Equal(expected, n))

	intTest := []byte("42")
	var o optional.Option[int]
	expected = optional.Some(42)
	err = json.Unmarshal(intTest, &o)
	assert.NilError(t, err)
	assert.Assert(t, optional.Equal(expected, o))

	var p optional.Option[uint]
	uintExpected := optional.Some[uint](42)
	err = json.Unmarshal(intTest, &p)
	assert.NilError(t, err)
	assert.Assert(t, optional.Equal(uintExpected, p))

	floatTest := []byte("42.1")
	var q optional.Option[float64]
	floatExpected := optional.Some(42.1)
	err = json.Unmarshal(floatTest, &q)
	assert.NilError(t, err)
	assert.Assert(t, optional.Equal(floatExpected, q))
}

func TestOptionUnMarshalJSONString(t *testing.T) {
	var n optional.Option[string]
	expected := optional.None[string]()
	nullTest := []byte("null")
	err := json.Unmarshal(nullTest, &n)
	assert.NilError(t, err)
	assert.Assert(t, optional.Equal(expected, n))

	test := []byte(`"this is a json string"`)
	var o optional.Option[string]
	expected = optional.Some("this is a json string")
	err = json.Unmarshal(test, &o)
	assert.NilError(t, err)
	assert.Assert(t, optional.Equal(expected, o))
}

func TestOptionUnMarshalJSONBool(t *testing.T) {
	var n optional.Option[bool]
	expected := optional.None[bool]()
	nullTest := []byte("null")
	err := json.Unmarshal(nullTest, &n)
	assert.NilError(t, err)
	assert.Assert(t, optional.Equal(expected, n))

	test := []byte("true")
	var o optional.Option[bool]
	expected = optional.Some(true)
	err = json.Unmarshal(test, &o)
	assert.NilError(t, err)
	assert.Assert(t, optional.Equal(expected, o))
}

func TestOptionUnMarshalJSONArray(t *testing.T) {
	var s []optional.Option[int]
	expected := []optional.Option[int]{optional.Some(42), optional.Some(49), optional.None[int]()}
	test := []byte("[42,49,null]")
	err := json.Unmarshal(test, &s)
	assert.NilError(t, err)
	for _, opt := range expected {
		assert.Assert(t, slices.Contains(s, opt))
	}
}

func TestOptionUnMarshalJSONObject(t *testing.T) {
	var m map[string]optional.Option[int]
	expected := map[string]optional.Option[int]{"one": optional.Some(42), "two": optional.Some(49), "three": optional.None[int]()}
	test := []byte(`{"one":42,"two":49,"three":null}`)
	err := json.Unmarshal(test, &m)
	assert.NilError(t, err)
	for name, opt := range expected {
		val, ok := m[name]
		assert.Assert(t, ok)
		assert.Assert(t, optional.Equal(opt, val))
	}
}

func TestOptionUnMarshalJSON(t *testing.T) {
	type Inner struct {
		Number optional.Option[int]
	}

	type S struct {
		Animal  optional.Option[string]
		Mineral optional.Option[string]
		Man     optional.Option[string]
		Other   Inner
	}

	str := []byte(`"test1"`)
	expectedStr := optional.Some("test1")
	none := []byte(`null`)
	expectedNone := optional.None[string]()

	var o optional.Option[string]
	err := json.Unmarshal(str, &o)
	assert.NilError(t, err)
	assert.Assert(t, optional.Equal(expectedStr, o))

	var p optional.Option[string]
	err = json.Unmarshal([]byte(none), &p)
	assert.NilError(t, err)
	assert.Assert(t, optional.Equal(expectedNone, p))

	rawJSON := []byte(`{"animal": "monkey", "mineral": null, "man": "Lincon", "other": { "number": 42 }}`)
	expected := S{optional.Some("monkey"), optional.None[string](), optional.Some("Lincon"), Inner{optional.Some(42)}}

	s := S{}
	err = json.Unmarshal(rawJSON, &s)
	assert.NilError(t, err)
	assert.Equal(t, expected, s)
}
