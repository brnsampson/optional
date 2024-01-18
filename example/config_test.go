package main_test

import (
	"github.com/brnsampson/optional/confopt"
	main "github.com/brnsampson/optional/example"
	"gotest.tools/v3/assert"
	"testing"
)

// Just a few basic tests to make sure that I'm consistent with which side of these functions wins

func TestSubConfigBuilderPieces(t *testing.T) {
	port := 1234
	portOption := confopt.SomeInt(port)
	port2 := 4000
	port2Option := confopt.SomeInt(port2)
	none := confopt.NoInt()

	// Check OrPort works as expected
	scl := main.NewSubConfigLoader()
	assert.Assert(t, scl.Port.IsNone())

	scl = scl.OrPort(portOption)
	assert.Assert(t, scl.Port.Match(port))

	// Or should only use the new value if the old value was None
	scl = scl.OrPort(port2Option)
	assert.Assert(t, scl.Port.Match(port))

	// Test WithPort
	scl = main.NewSubConfigLoader()
	assert.Assert(t, scl.Port.IsNone())

	scl = scl.WithPort(portOption)
	assert.Assert(t, scl.Port.Match(port))
	// WithPort should always prefer the new value (unless it is None)
	scl = scl.WithPort(port2Option)
	assert.Assert(t, scl.Port.Match(port2))

	scl = scl.WithPort(none)
	assert.Assert(t, scl.Port.Match(port2))
}

func TestSubConfigLoaderMerge(t *testing.T) {
	port1 := 1234
	port1Option := confopt.SomeInt(port1)
	scl1 := main.NewSubConfigLoader().WithPort(port1Option)
	port2 := 4000
	port2Option := confopt.SomeInt(port2)
	scl2 := main.NewSubConfigLoader().WithPort(port2Option)

	expected := main.SubConfigLoader{port1Option}
	res := scl1.Merged(scl2)
	assert.Assert(t, expected.Port.Eq(res.Port))
}

func TestConfigBuilderPieces(t *testing.T) {
	host := "myhost"
	hostOption := confopt.SomeStr(host)
	host2 := "otherhost"
	host2Option := confopt.SomeStr(host2)
	none := confopt.NoStr()

	// Check OrHost works as expected
	scl := main.NewConfigLoader()
	assert.Assert(t, scl.Host.IsNone())

	scl = scl.OrHost(hostOption)
	assert.Assert(t, scl.Host.Match(host))

	// Or should only use the new value if the old value was None
	scl = scl.OrHost(host2Option)
	assert.Assert(t, scl.Host.Match(host))

	// Test WithHost
	scl = main.NewConfigLoader()
	assert.Assert(t, scl.Host.IsNone())

	scl = scl.WithHost(hostOption)
	assert.Assert(t, scl.Host.Match(host))
	// WithHost should always prefer the new value (unless it is None)
	scl = scl.WithHost(host2Option)
	assert.Assert(t, scl.Host.Match(host2))

	scl = scl.WithHost(none)
	assert.Assert(t, scl.Host.Match(host2))
}

func TestConfigLoaderMerge(t *testing.T) {
	port1 := 1234
	port1Option := confopt.SomeInt(port1)
	scl1 := main.NewSubConfigLoader().WithPort(port1Option)

	host1 := "myhost"
	host1Option := confopt.SomeStr(host1)
	sc1 := main.NewConfigLoader().WithHost(host1Option)
	host2 := "otherHost"
	host2Option := confopt.SomeStr(host2)
	sc2 := main.NewConfigLoader().WithHost(host2Option)

	expected := main.ConfigLoader{confopt.NoStr(), host1Option, scl1}
	res := sc1.Merged(sc2)
	assert.Assert(t, expected.Host.Eq(res.Host))
}
