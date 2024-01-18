package confopt_test

import (
	"github.com/brnsampson/optional"
	"github.com/brnsampson/optional/confopt"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

func TestAnd(t *testing.T) {
	S1 := confopt.NoStr()
	S2 := confopt.SomeStr("Testing this thinger")
	res := optional.And(&S1, &S2)

	assert.Assert(t, res.IsNone())
}

func TestOr(t *testing.T) {
	S1 := confopt.NoStr()
	S2 := confopt.SomeStr("Testing this thinger")
	res := optional.Or(&S1, &S2)

	assert.Assert(t, res.IsSome())
	assert.Equal(t, res.String(), "Testing this thinger")
}

func TestIntIsConfigOptional(t *testing.T) {
	// The real test is if we get compiler errors because &I does not implement ConfigOptional
	I := confopt.SomeInt(42)
	var OI confopt.ConfigOptional[int] = &I

	assert.Assert(t, OI.IsSome())

	I8 := confopt.SomeInt8(42)
	var OI8 confopt.ConfigOptional[int8] = &I8

	assert.Assert(t, OI8.IsSome())

	I16 := confopt.SomeInt16(42)
	var OI16 confopt.ConfigOptional[int16] = &I16

	assert.Assert(t, OI16.IsSome())

	I32 := confopt.SomeInt32(42)
	var OI32 confopt.ConfigOptional[int32] = &I32

	assert.Assert(t, OI32.IsSome())

	I64 := confopt.SomeInt64(42)
	var OI64 confopt.ConfigOptional[int64] = &I64

	assert.Assert(t, OI64.IsSome())
}

func TestUintIsConfigOptional(t *testing.T) {
	// The real test is if we get compiler errors because &UI does not implement ConfigOptional
	UI := confopt.SomeUint(42)
	var OUI confopt.ConfigOptional[uint] = &UI

	assert.Assert(t, OUI.IsSome())

	UI8 := confopt.SomeUint8(42)
	var OUI8 confopt.ConfigOptional[uint8] = &UI8

	assert.Assert(t, OUI8.IsSome())

	UI16 := confopt.SomeUint16(42)
	var OUI16 confopt.ConfigOptional[uint16] = &UI16

	assert.Assert(t, OUI16.IsSome())

	UI32 := confopt.SomeUint32(42)
	var OUI32 confopt.ConfigOptional[uint32] = &UI32

	assert.Assert(t, OUI32.IsSome())

	UI64 := confopt.SomeUint64(42)
	var OUI64 confopt.ConfigOptional[uint64] = &UI64

	assert.Assert(t, OUI64.IsSome())
}

func TestFloatIsConfigOptional(t *testing.T) {
	F32 := confopt.SomeFloat32(42.1)
	var OF32 confopt.ConfigOptional[float32] = &F32

	assert.Assert(t, OF32.IsSome())

	F64 := confopt.SomeFloat64(42.1)
	var OF64 confopt.ConfigOptional[float64] = &F64

	assert.Assert(t, OF64.IsSome())
}

func TestBoolIsConfigOptional(t *testing.T) {
	B := confopt.SomeBool(true)
	var OB confopt.ConfigOptional[bool] = &B

	assert.Assert(t, OB.IsSome())
}

func TestStringIsConfigOptional(t *testing.T) {
	S := confopt.SomeStr("testing")
	var OS confopt.ConfigOptional[string] = &S

	assert.Assert(t, OS.IsSome())
}

func TestTimeIsConfigOptional(t *testing.T) {
	T := confopt.SomeTime(time.Now())
	var OT confopt.ConfigOptional[time.Time] = &T

	assert.Assert(t, OT.IsSome())
}

func TestCertIsConfigOptional(t *testing.T) {
	T, err := confopt.SomeCert("/not/a/real/path")
	assert.NilError(t, err)

	var OT confopt.ConfigOptional[string] = &T
	assert.Assert(t, OT.IsSome())
}

func TestPubKeyIsConfigOptional(t *testing.T) {
	T, err := confopt.SomePubKey("/not/a/real/path")
	assert.NilError(t, err)

	var OT confopt.ConfigOptional[string] = &T

	assert.Assert(t, OT.IsSome())
}

func TestPrivateKeyIsConfigOptional(t *testing.T) {
	T, err := confopt.SomePrivateKey("/not/a/real/path")
	assert.NilError(t, err)

	var OT confopt.ConfigOptional[string] = &T

	assert.Assert(t, OT.IsSome())
}

func TestFileIsConfigOptional(t *testing.T) {
	T := confopt.SomeFile("/not/a/real/path")

	var OT confopt.ConfigOptional[string] = &T

	assert.Assert(t, OT.IsSome())
}
