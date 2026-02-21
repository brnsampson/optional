package optional_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
)

func TestUintType(t *testing.T) {
	o := optional.SomeUint(42)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestUintString(t *testing.T) {
	var i uint = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint(i)
	assert.Equal(t, iStr, o.String())
}

func TestUintMarshalText(t *testing.T) {
	var i uint = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestUintUnmarshalText(t *testing.T) {
	var i uint = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := optional.NoUint()
	err := o.UnmarshalText([]byte(iStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, i, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-uint
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestUint8Type(t *testing.T) {
	o := optional.SomeUint8(42)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestUint8String(t *testing.T) {
	var i uint8 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint8(i)
	assert.Equal(t, iStr, o.String())
}

func TestUint8MarshalText(t *testing.T) {
	var i uint8 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint8(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestUint8UnmarshalText(t *testing.T) {
	var i uint8 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := optional.NoUint8()
	err := o.UnmarshalText([]byte(iStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, i, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-uint
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestUint8Sql(t *testing.T) {
	ins := "INSERT INTO optionTest (val) VALUES (?)"
	q := `SELECT * FROM optionTest WHERE val = ?`
	test1Val := uint8(8)
	test1 := optional.SomeUint8(test1Val)
	test2 := optional.NoUint8()
	out1 := optional.NoUint8()
	out2 := optional.SomeUint8(8)

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

func TestUint16Type(t *testing.T) {
	o := optional.SomeUint16(42)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestUint16String(t *testing.T) {
	var i uint16 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint16(i)
	assert.Equal(t, iStr, o.String())
}

func TestUint16MarshalText(t *testing.T) {
	var i uint16 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint16(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestUint16UnmarshalText(t *testing.T) {
	var i uint16 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := optional.NoUint16()
	err := o.UnmarshalText([]byte(iStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, i, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-uint
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestUint16Sql(t *testing.T) {
	ins := "INSERT INTO optionTest (val) VALUES (?)"
	q := `SELECT * FROM optionTest WHERE val = ?`
	test1Val := uint16(16)
	test1 := optional.SomeUint16(test1Val)
	test2 := optional.NoUint16()
	out1 := optional.NoUint16()
	out2 := optional.SomeUint16(16)

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

func TestUint32Type(t *testing.T) {
	o := optional.SomeUint32(42)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestUint32String(t *testing.T) {
	var f uint32 = 42
	fStr := strconv.FormatUint(uint64(f), 10)
	o := optional.SomeUint32(f)
	assert.Equal(t, fStr, o.String())
}

func TestUint32MarshalText(t *testing.T) {
	var f uint32 = 42
	fStr := strconv.FormatUint(uint64(f), 10)
	o := optional.SomeUint32(f)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, fStr, string(s))
}

func TestUint32UnmarshalText(t *testing.T) {
	var i uint32 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := optional.NoUint32()
	err := o.UnmarshalText([]byte(iStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, i, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-uint
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestUint32Sql(t *testing.T) {
	ins := "INSERT INTO optionTest (val) VALUES (?)"
	q := `SELECT * FROM optionTest WHERE val = ?`
	test1Val := uint32(32)
	test1 := optional.SomeUint32(test1Val)
	test2 := optional.NoUint32()
	out1 := optional.NoUint32()
	out2 := optional.SomeUint32(32)

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

func TestUint64Type(t *testing.T) {
	o := optional.SomeUint64(42)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestUint64String(t *testing.T) {
	var i uint64 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint64(i)
	assert.Equal(t, iStr, o.String())
}

func TestUint64MarshalText(t *testing.T) {
	var i uint64 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	o := optional.SomeUint64(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestUint64UnmarshalText(t *testing.T) {
	var i uint64 = 42
	iStr := strconv.FormatUint(uint64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := optional.NoUint64()
	err := o.UnmarshalText([]byte(iStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, i, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-uint
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestUint64Sql(t *testing.T) {
	ins := "INSERT INTO optionTest (val) VALUES (?)"
	q := `SELECT * FROM optionTest WHERE val = ?`
	test1Val := uint64(64)
	test1 := optional.SomeUint64(test1Val)
	test2 := optional.NoUint64()
	out1 := optional.NoUint64()
	out2 := optional.SomeUint64(64)

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
