package optional_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"slices"
	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
	"testing"
)

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
	o.Set(49)
	assert.Assert(t, clone.IsSome())
	assert.Equal(t, 42, clone.Must())
}

func TestMutableOptionalMutableClone(t *testing.T) {
	o := optional.Some(42)
	var op optional.MutableOptional[int] = &o
	clone := op.MutableClone()
	clone.Set(49)
	assert.Assert(t, op.IsSome())
	assert.Equal(t, 42, op.MustUnwrap())
}

func TestOptionBasics(t *testing.T) {
	// Covers IsSome, IsNone, Clear, Set, and Get
	val := 42
	val2 := 49
	errString := "Attempted to Get Option with None value"

	o := optional.Some(val)
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	tmp, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, val, tmp)
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	o.Clear()
	assert.Assert(t, !o.IsSome())
	assert.Assert(t, o.IsNone())

	tmp, err = o.Get()
	assert.Error(t, err, errString)

	o.Set(val2)
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	tmp, err = o.Get()
	assert.NilError(t, err)
	assert.Equal(t, val2, tmp)
}

func TestOptionGetOr(t *testing.T) {
	val := 42
	val2 := 49

	o := optional.Some(val)
	none := optional.None[int]()

	tmp := o.GetOr(val2)
	assert.Equal(t, val, tmp)

	tmp = none.GetOr(val2)
	assert.Equal(t, val2, tmp)
	assert.Assert(t, none.IsNone())
}

func TestOptionGetOrInsert(t *testing.T) {
	val := 42
	val2 := 49

	o := optional.Some(val)
	none := optional.None[int]()

	tmp := o.GetOrInsert(val2)
	assert.Equal(t, val, tmp)

	tmp = none.GetOrInsert(val2)
	assert.Equal(t, val2, tmp)
	assert.Assert(t, none.IsSome())
	tmp = 0

	tmp, err := none.Get()
	assert.NilError(t, err)
	assert.Equal(t, val2, tmp)
}

func TestOptionMustPanics(t *testing.T) {
	defer func() { _ = recover() }()
	o := optional.None[int]()
	_ = o.Must()

	t.Errorf("Must() failed to panic on a None value Option")
}

func TestOptionUnwrap(t *testing.T) {
	// Covers Unwrap, MustUnwrap, UnwrapOr, and UnwrapOrElse
	val := 42
	errString := "Attempted to Get Option with None value"
	o := optional.Some(val)
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	tmp, err := o.Unwrap()
	assert.NilError(t, err)
	assert.Equal(t, val, tmp)
	assert.Assert(t, !o.IsSome())
	assert.Assert(t, o.IsNone())

	tmp, err = o.Unwrap()
	assert.Error(t, err, errString)
}

func TestOptionMustUnwrapPanics(t *testing.T) {
	defer func() { _ = recover() }()
	o := optional.None[int]()
	_ = o.MustUnwrap()

	t.Errorf("MustUnwrap() failed to panic on a None value Option")
}

func TestOptionMustUnwrap(t *testing.T) {
	// Covers Unwrap, MustUnwrap, UnwrapOr, and UnwrapOrElse
	val := 42

	o := optional.Some(val)
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	tmp := o.MustUnwrap()
	assert.Equal(t, val, tmp)
	assert.Assert(t, !o.IsSome())
	assert.Assert(t, o.IsNone())
}

func TestOptionUnwrapOr(t *testing.T) {
	// Covers Unwrap, MustUnwrap, UnwrapOr, and UnwrapOrElse
	val := 42
	def_val := 49

	o := optional.Some(val)
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	tmp := o.UnwrapOr(def_val)
	assert.Equal(t, val, tmp)
	assert.Assert(t, !o.IsSome())
	assert.Assert(t, o.IsNone())

	tmp = o.UnwrapOr(def_val)
	assert.Equal(t, def_val, tmp)
	assert.Assert(t, !o.IsSome())
	assert.Assert(t, o.IsNone())
}

func TestOptionUnwrapOrElse(t *testing.T) {
	// Covers Unwrap, MustUnwrap, UnwrapOr, and UnwrapOrElse
	val := 42
	def_val := 49
	def_func := func() int { return def_val }

	o := optional.Some(val)
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	tmp := o.UnwrapOrElse(def_func)
	assert.Equal(t, val, tmp)
	assert.Assert(t, !o.IsSome())
	assert.Assert(t, o.IsNone())

	tmp = o.UnwrapOrElse(def_func)
	assert.Equal(t, def_val, tmp)
	assert.Assert(t, !o.IsSome())
	assert.Assert(t, o.IsNone())
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

func TestOptionEq(t *testing.T) {
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

func TestOptionAnd(t *testing.T) {
	val1 := 42
	val2 := 49
	o1 := optional.Some(val1)
	o2 := optional.Some(val2)
	none := optional.None[int]()

	tmpO := o1.And(&none)
	assert.Assert(t, tmpO.IsNone())

	tmpO = none.And(&o1)
	assert.Assert(t, tmpO.IsNone())

	tmpO = o1.And(&o2)
	tmp := tmpO.Must()
	assert.Equal(t, val2, tmp)

	tmpO = o2.And(&o1)
	tmp = tmpO.Must()
	assert.Equal(t, val1, tmp)
}

func TestOptionOr(t *testing.T) {
	val1 := 42
	val2 := 49
	o1 := optional.Some(val1)
	o2 := optional.Some(val2)
	none := optional.None[int]()
	none2 := optional.None[int]()

	tmpO := none.Or(&o1)
	tmp := tmpO.Must()
	assert.Equal(t, val1, tmp)

	tmpO = o1.Or(&none)
	tmp = tmpO.Must()
	assert.Equal(t, val1, tmp)

	tmpO = o2.Or(&o1)
	tmp = tmpO.Must()
	assert.Equal(t, val2, tmp)

	tmpO = none.Or(&none2)
	assert.Assert(t, tmpO.IsNone())
}

func TestOptionTransform(t *testing.T) {
	val := 42
	transform := func(x int) int { return x + 7 }
	after_val := transform(val)

	o := optional.Some(val)
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	o.Transform(transform)

	tmp, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, after_val, tmp)

	// Transforming a None value should still result in a None.
	none := optional.None[int]()
	none.Transform(transform)
	assert.Assert(t, none.IsNone())
}

func TestOptionTransformOr(t *testing.T) {
	val := 42
	def := 99
	transform := func(x int) int { return x + 7 }
	after_val := transform(val)

	o := optional.Some(val)
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	o.TransformOr(transform, def)
	tmp, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, after_val, tmp)

	// Transforming a None value makes it a Some(val)
	none := optional.None[int]()
	none.TransformOr(transform, def)

	tmp, err = none.Get()
	assert.NilError(t, err)
	assert.Equal(t, def, tmp)
}

func TestOptionTransformOrError(t *testing.T) {
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
	err := none.TransformOrError(transform)
	assert.NilError(t, err)
	assert.Assert(t, none.IsNone())

	// Transforming eith a function that doesn't err is the same as a normal transform
	err = o.TransformOrError(transform)
	assert.NilError(t, err)

	tmp, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, after_val, tmp)

	// transforming with a function that produces an error should return the same error
	err = o.TransformOrError(err_transform)
	assert.Error(t, err, err_string)

	// The error should prevent the transform from actually being applied
	tmp, err = o.Get()
	assert.NilError(t, err)
	assert.Equal(t, after_val, tmp)
}

func TestOptionBinaryTransform(t *testing.T) {
	val := 42
	inc := 7

	transform := func(x, y int) int { return x + y }
	after_val := transform(val, inc)

	o := optional.Some(val)
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	o.BinaryTransform(inc, transform)
	tmp, err := o.Get()
	assert.NilError(t, err)
	assert.Equal(t, after_val, tmp)

	// Transforming a None value leaves it as a None
	none := optional.None[int]()
	none.BinaryTransform(inc, transform)
	assert.Assert(t, none.IsNone())
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
	assert.Assert(t, expected.Eq(&n))

	intTest := []byte("42")
	var o optional.Option[int]
	expected = optional.Some(42)
	err = json.Unmarshal(intTest, &o)
	assert.NilError(t, err)
	assert.Assert(t, expected.Eq(&o))

	var p optional.Option[uint]
	uintExpected := optional.Some[uint](42)
	err = json.Unmarshal(intTest, &p)
	assert.NilError(t, err)
	assert.Assert(t, uintExpected.Eq(&p))

	floatTest := []byte("42.1")
	var q optional.Option[float64]
	floatExpected := optional.Some(42.1)
	err = json.Unmarshal(floatTest, &q)
	assert.NilError(t, err)
	assert.Assert(t, floatExpected.Eq(&q))
}

func TestOptionUnMarshalJSONString(t *testing.T) {
	var n optional.Option[string]
	expected := optional.None[string]()
	nullTest := []byte("null")
	err := json.Unmarshal(nullTest, &n)
	assert.NilError(t, err)
	assert.Assert(t, expected.Eq(&n))

	test := []byte(`"this is a json string"`)
	var o optional.Option[string]
	expected = optional.Some("this is a json string")
	err = json.Unmarshal(test, &o)
	assert.NilError(t, err)
	assert.Assert(t, expected.Eq(&o))
}

func TestOptionUnMarshalJSONBool(t *testing.T) {
	var n optional.Option[bool]
	expected := optional.None[bool]()
	nullTest := []byte("null")
	err := json.Unmarshal(nullTest, &n)
	assert.NilError(t, err)
	assert.Assert(t, expected.Eq(&n))

	test := []byte("true")
	var o optional.Option[bool]
	expected = optional.Some(true)
	err = json.Unmarshal(test, &o)
	assert.NilError(t, err)
	assert.Assert(t, expected.Eq(&o))
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
		assert.Assert(t, opt.Eq(&val))
	}
}

func TestOptionUnMarshalJSON(t *testing.T) {
	type Inner struct {
		Number optional.Option[int]
	}

	type S struct {
		Animal optional.Option[string]
		Mineral optional.Option[string]
		Man optional.Option[string]
		Other Inner
	}

	str := []byte(`"test1"`)
	expectedStr := optional.Some("test1")
	none := []byte(`null`)
	expectedNone := optional.None[string]()

	var o optional.Option[string]
	err := json.Unmarshal(str, &o)
	assert.NilError(t, err)
	assert.Assert(t, o.Eq(&expectedStr))

	var p optional.Option[string]
	err = json.Unmarshal([]byte(none), &p)
	assert.NilError(t, err)
	assert.Assert(t, p.Eq(&expectedNone))

	rawJSON := []byte(`{"animal": "monkey", "mineral": null, "man": "Lincon", "other": { "number": 42 }}`)
	expected := S{ optional.Some("monkey"), optional.None[string](), optional.Some("Lincon"), Inner { optional.Some(42) } }

	s := S{}
	err = json.Unmarshal(rawJSON, &s)
	assert.NilError(t, err)
	assert.Equal(t, expected, s)
}
