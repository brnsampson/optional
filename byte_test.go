package optional_test

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
)

func TestByteType(t *testing.T) {
	o := optional.SomeByte('r')
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestByteString(t *testing.T) {
	byteString := "6e"
	o := optional.SomeByte('n')
	assert.Equal(t, byteString, o.String())
}

func TestByteMarshalText(t *testing.T) {
	byteString := "6e"
	o := optional.SomeByte('n')

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, byteString, string(s))
}

func TestByteUnmarshalText(t *testing.T) {
	byteString := "6e"
	nullString := "null"
	longString := "this is more than one byte"

	// Text sucessful unmarshaling
	o := optional.NoByte()
	err := o.UnmarshalText([]byte(byteString))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, byte('n'), ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullString))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-byte
	err = o.UnmarshalText([]byte(longString))
	assert.Assert(t, err != nil)
}

func TestByteSql(t *testing.T) {
	ins := "INSERT INTO optionTest (val) VALUES (?)"
	q := `SELECT * FROM optionTest WHERE val = ?`
	test1Val := byte('n')
	test1 := optional.SomeByte(test1Val)
	test2 := optional.NoByte()
	out1 := optional.NoByte()
	out2 := optional.SomeByte(byte('n'))

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
