package optional_test

import (
	"testing"

	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
)

func TestOptionIsOptional(t *testing.T) {
	// The real test is if we get compiler errors because &o does not implement Optional
	o := optional.Some(42)
	var op optional.Optional[int] = o
	assert.Assert(t, op.IsSome())
}

func TestEqual(t *testing.T) {
	s := "Testing this thinger"
	S1 := optional.None[string]()
	S2 := optional.Some(s)
	S3 := optional.Some(s)
	none := optional.None[string]()
	assert.Assert(t, !optional.Equal(S1, S2))
	assert.Assert(t, optional.Equal(S2, S3))
	assert.Assert(t, optional.Equal(&S2, &S3))
	assert.Assert(t, optional.Equal(S1, none))

}

func TestAnd(t *testing.T) {
	s := "Testing this thinger"
	S1 := optional.None[string]()
	S2 := optional.Some(s)
	res := optional.And(&S1, &S2)

	assert.Assert(t, res.IsNone())
}

func TestOr(t *testing.T) {
	s := "Testing this thinger"
	S1 := optional.None[string]()
	S2 := optional.Some(s)
	res := optional.Or(&S1, &S2)

	assert.Assert(t, res.IsSome())
	assert.Equal(t, s, res.MustGet())
}

// Skipping ClearIfMatch since it is literally just calling two covered functions
func TestGetOr(t *testing.T) {
	val := 42
	val2 := 49

	o := optional.Some(val)
	none := optional.None[int]()

	tmp := optional.GetOr(o, val2)
	assert.Equal(t, val, tmp)

	tmp = optional.GetOr(none, val2)
	assert.Equal(t, val2, tmp)
	assert.Assert(t, none.IsNone())
}

func TestGetOrInsert(t *testing.T) {
	val := 42
	val2 := 49

	o := optional.Some(val)
	none := optional.None[int]()

	tmp, err := optional.GetOrInsert(&o, val2)
	assert.NilError(t, err)
	assert.Equal(t, val, tmp)

	tmp, err = optional.GetOrInsert(&none, val2)
	assert.NilError(t, err)
	assert.Equal(t, val2, tmp)
	assert.Assert(t, none.IsSome())
	tmp = 0

	tmp, ok := none.Get()
	assert.Assert(t, ok)
	assert.Equal(t, val2, tmp)
}

func TestOptionMustPanics(t *testing.T) {
	defer func() { _ = recover() }()
	o := optional.None[int]()
	_ = o.MustGet()

	t.Errorf("Must() failed to panic on a None value Option")
}

func TestTransformOr(t *testing.T) {
	val := 42
	def := 99
	transform := func(x int) (int, error) { return x + 7, nil }
	after_val, err := transform(val)
	assert.NilError(t, err)
	after_def, err := transform(def)
	assert.NilError(t, err)

	o := optional.Some(val)
	assert.Assert(t, o.IsSome())
	assert.Assert(t, !o.IsNone())

	optional.TransformOr(&o, transform, def)
	tmp, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, after_val, tmp)

	// Transforming a None value makes it a Some(val)
	none := optional.None[int]()
	optional.TransformOr(&none, transform, def)

	tmp, ok = none.Get()
	assert.Assert(t, ok)
	assert.Equal(t, after_def, tmp)
}
