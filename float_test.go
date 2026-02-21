package optional_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
)

func TestFloat32Type(t *testing.T) {
	o := optional.SomeFloat32(42.0)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestFloat32String(t *testing.T) {
	var f float32 = 42.1
	fStr := strconv.FormatFloat(float64(f), 'g', 3, 32)
	o := optional.SomeFloat32(f)
	assert.Equal(t, fStr, o.String())
}

func TestFloat32MarshalText(t *testing.T) {
	var f float32 = 42.1
	fStr := strconv.FormatFloat(float64(f), 'g', 3, 32)
	o := optional.SomeFloat32(f)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, fStr, string(s))
}

func TestFloat32UnmarshalText(t *testing.T) {
	var f float32 = 42.1
	fStr := strconv.FormatFloat(float64(f), 'g', 3, 32)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := optional.NoFloat32()
	err := o.UnmarshalText([]byte(fStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, f, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-float
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestFloat64Type(t *testing.T) {
	o := optional.SomeFloat64(42.0)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestFloat64String(t *testing.T) {
	var f float64 = 42.1
	fStr := strconv.FormatFloat(float64(f), 'g', 3, 64)
	o := optional.SomeFloat64(f)
	assert.Equal(t, fStr, o.String())
}

func TestFloat64MarshalText(t *testing.T) {
	var f float64 = 42.1
	fStr := strconv.FormatFloat(float64(f), 'g', 3, 64)
	o := optional.SomeFloat64(f)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, fStr, string(s))
}

func TestFloat64UnmarshalText(t *testing.T) {
	var f float64 = 42.1
	fStr := strconv.FormatFloat(float64(f), 'g', 3, 64)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := optional.NoFloat64()
	err := o.UnmarshalText([]byte(fStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, f, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-float
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestFloat64Sql(t *testing.T) {
	ins := "INSERT INTO optionTest (val) VALUES (?)"
	q := `SELECT * FROM optionTest WHERE val = ?`
	test1Val := float64(64.64)
	test1 := optional.SomeFloat64(test1Val)
	test2 := optional.NoFloat64()
	out1 := optional.NoFloat64()
	out2 := optional.SomeFloat64(float64('n'))

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	assert.NilError(t, err, "failed to open mock database connection")
	defer db.Close()

	// mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO optionTest").WithArgs(test1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO optionTest").WithArgs(test2).WillReturnResult(sqlmock.NewResult(2, 1))
	r1 := sqlmock.NewRows([]string{"val"})
	r1.AddRow(test1Val)
	select1 := mock.ExpectQuery(`SELECT (.+) FROM optionTest`)
	select1.WithArgs(test1).WillReturnRows(r1)
	r2 := sqlmock.NewRows([]string{"val"})
	r2.AddRow(nil)
	select2 := mock.ExpectQuery(`SELECT (.+) FROM optionTest`)
	select2.WithArgs(test2).WillReturnRows(r2)

	_, err = db.Exec(ins, test1)
	assert.NilError(t, err, "error using mock to insert optional Some type")
	_, err = db.Exec(ins, test2)
	assert.NilError(t, err, "error using mock to insert optional None type")

	rows1, err := db.Query(q, test1)
	assert.NilError(t, err, "error using mock to query with optional Some type")
	defer rows1.Close()

	for rows1.Next() {
		err = rows1.Scan(&out1)
		assert.NilError(t, err, "error using Scan to convert sql row to optional type")
		assert.Equal(t, test1, out1, "Some Scan case failed")
	}

	rows2, err := db.Query(q, test2)
	assert.NilError(t, err, "error using mock to query with optional None type")
	defer rows2.Close()

	for rows2.Next() {
		err = rows2.Scan(&out2)
		assert.NilError(t, err, "error using Scan to convert sql row to optional None type")
		assert.Assert(t, out2.IsNone(), "None case failed")
	}
}
