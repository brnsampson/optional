package config_test

import (
	"time"
	"github.com/brnsampson/optional/config"
	"gotest.tools/v3/assert"
	"testing"
)

func TestIntIsConfigOptional(t *testing.T) {
	// The real test is if we get compiler errors because &I does not implement ConfigOptional
	I := config.SomeInt(42)
	var OI config.ConfigOptional[int] = &I

	assert.Assert(t, OI.IsSome())

	I8 := config.SomeInt8(42)
	var OI8 config.ConfigOptional[int8] = &I8

	assert.Assert(t, OI8.IsSome())

	I16 := config.SomeInt16(42)
	var OI16 config.ConfigOptional[int16] = &I16

	assert.Assert(t, OI16.IsSome())

	I32 := config.SomeInt32(42)
	var OI32 config.ConfigOptional[int32] = &I32

	assert.Assert(t, OI32.IsSome())

	I64 := config.SomeInt64(42)
	var OI64 config.ConfigOptional[int64] = &I64

	assert.Assert(t, OI64.IsSome())
}

func TestUintIsConfigOptional(t *testing.T) {
	// The real test is if we get compiler errors because &UI does not implement ConfigOptional
	UI := config.SomeUint(42)
	var OUI config.ConfigOptional[uint] = &UI

	assert.Assert(t, OUI.IsSome())

	UI8 := config.SomeUint8(42)
	var OUI8 config.ConfigOptional[uint8] = &UI8

	assert.Assert(t, OUI8.IsSome())

	UI16 := config.SomeUint16(42)
	var OUI16 config.ConfigOptional[uint16] = &UI16

	assert.Assert(t, OUI16.IsSome())

	UI32 := config.SomeUint32(42)
	var OUI32 config.ConfigOptional[uint32] = &UI32

	assert.Assert(t, OUI32.IsSome())

	UI64 := config.SomeUint64(42)
	var OUI64 config.ConfigOptional[uint64] = &UI64

	assert.Assert(t, OUI64.IsSome())
}

func TestFloatIsConfigOptional(t *testing.T) {
	F32 := config.SomeFloat32(42.1)
	var OF32 config.ConfigOptional[float32] = &F32

	assert.Assert(t, OF32.IsSome())

	F64 := config.SomeFloat64(42.1)
	var OF64 config.ConfigOptional[float64] = &F64

	assert.Assert(t, OF64.IsSome())
}

func TestBoolIsConfigOptional(t *testing.T) {
	B := config.SomeBool(true)
	var OB config.ConfigOptional[bool] = &B

	assert.Assert(t, OB.IsSome())
}

func TestStringIsConfigOptional(t *testing.T) {
	S := config.SomeStr("testing")
	var OS config.ConfigOptional[string] = &S

	assert.Assert(t, OS.IsSome())
}

func TestTimeIsConfigOptional(t *testing.T) {
	T := config.SomeTime(time.Now())
	var OT config.ConfigOptional[time.Time] = &T

	assert.Assert(t, OT.IsSome())
}

func TestCertIsCertOptional(t *testing.T) {
	T := config.SomeCert("/not/a/real/path")
	var OT config.CertOptional = &T

	assert.Assert(t, OT.IsSome())
}

func TestPubKeyIsPubKeyOptional(t *testing.T) {
	T := config.SomePubKey("/not/a/real/path")
	var OT config.PubKeyOptional = &T

	assert.Assert(t, OT.IsSome())
}
