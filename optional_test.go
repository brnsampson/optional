package optional_test

import (
	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
	"testing"
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
	assert.Equal(t, s, res.MustUnwrap())
}
