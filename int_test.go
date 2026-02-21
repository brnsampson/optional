package optional_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
)

func TestIntType(t *testing.T) {
	o := optional.SomeInt(42)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestIntString(t *testing.T) {
	var i int = 42
	iStr := strconv.Itoa(i)
	o := optional.SomeInt(i)
	assert.Equal(t, iStr, o.String())
}

func TestIntMarshalText(t *testing.T) {
	var i int = 42
	iStr := strconv.Itoa(i)
	o := optional.SomeInt(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestIntUnmarshalText(t *testing.T) {
	var i int = 42
	iStr := strconv.Itoa(i)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := optional.NoInt()
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

	// Test unmarshaling non-int
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestInt8Type(t *testing.T) {
	o := optional.SomeInt8(42)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestInt8String(t *testing.T) {
	var i int8 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := optional.SomeInt8(i)
	assert.Equal(t, iStr, o.String())
}

func TestInt8MarshalText(t *testing.T) {
	var i int8 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := optional.SomeInt8(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestInt8UnmarshalText(t *testing.T) {
	var i int8 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := optional.NoInt8()
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

	// Test unmarshaling non-int
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestInt8Sql(t *testing.T) {
	ins := "INSERT INTO optionTest (val) VALUES (?)"
	q := `SELECT * FROM optionTest WHERE val = ?`
	test1Val := int8(8)
	test1 := optional.SomeInt8(test1Val)
	test2 := optional.NoInt8()
	out1 := optional.NoInt8()
	out2 := optional.SomeInt8(8)

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

func TestInt16Type(t *testing.T) {
	o := optional.SomeInt16(42)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestInt16String(t *testing.T) {
	var i int16 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := optional.SomeInt16(i)
	assert.Equal(t, iStr, o.String())
}

func TestInt16MarshalText(t *testing.T) {
	var i int16 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := optional.SomeInt16(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestInt16UnmarshalText(t *testing.T) {
	var i int16 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := optional.NoInt16()
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

	// Test unmarshaling non-int
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestInt16Sql(t *testing.T) {
	ins := "INSERT INTO optionTest (val) VALUES (?)"
	q := `SELECT * FROM optionTest WHERE val = ?`
	test1Val := int16(16)
	test1 := optional.SomeInt16(test1Val)
	test2 := optional.NoInt16()
	out1 := optional.NoInt16()
	out2 := optional.SomeInt16(16)

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

func TestInt32Type(t *testing.T) {
	o := optional.SomeInt32(42)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestInt32String(t *testing.T) {
	var i int32 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := optional.SomeInt32(i)
	assert.Equal(t, iStr, o.String())
}

func TestInt32MarshalText(t *testing.T) {
	var i int32 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := optional.SomeInt32(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestInt32UnmarshalText(t *testing.T) {
	var i int32 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := optional.NoInt32()
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

	// Test unmarshaling non-int
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestInt32Sql(t *testing.T) {
	ins := "INSERT INTO optionTest (val) VALUES (?)"
	q := `SELECT * FROM optionTest WHERE val = ?`
	test1Val := int32(32)
	test1 := optional.SomeInt32(test1Val)
	test2 := optional.NoInt32()
	out1 := optional.NoInt32()
	out2 := optional.SomeInt32(32)

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

func TestInt64Type(t *testing.T) {
	o := optional.SomeInt64(42)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestInt64String(t *testing.T) {
	var i int64 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := optional.SomeInt64(i)
	assert.Equal(t, iStr, o.String())
}

func TestInt64MarshalText(t *testing.T) {
	var i int64 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	o := optional.SomeInt64(i)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, iStr, string(s))
}

func TestInt64UnmarshalText(t *testing.T) {
	var i int64 = 42
	iStr := strconv.FormatInt(int64(i), 10)
	nullStr := "null"
	s := "this is not a number"

	// Text sucessful unmarshaling
	o := optional.NoInt64()
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

	// Test unmarshaling non-int
	err = o.UnmarshalText([]byte(s))
	assert.Assert(t, err != nil)
}

func TestInt64Sql(t *testing.T) {
	ins := "INSERT INTO optionTest (val) VALUES (?)"
	q := `SELECT * FROM optionTest WHERE val = ?`
	test1Val := int64(64)
	test1 := optional.SomeInt64(test1Val)
	test2 := optional.NoInt64()
	out1 := optional.NoInt64()
	out2 := optional.SomeInt64(64)

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
