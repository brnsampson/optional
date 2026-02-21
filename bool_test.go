package optional_test

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
)

func TestBoolType(t *testing.T) {
	o := optional.SomeBool(true)
	assert.Equal(t, reflect.TypeOf(o).Name(), o.Type())
}

func TestBoolString(t *testing.T) {
	trueString := "true"
	o := optional.SomeBool(true)
	assert.Equal(t, trueString, o.String())
}

func TestBoolMarshalText(t *testing.T) {
	trueString := "true"
	o := optional.SomeBool(true)

	s, err := o.MarshalText()
	assert.NilError(t, err)
	assert.Equal(t, trueString, string(s))
}

func TestBoolUnmarshalText(t *testing.T) {
	trueString := "true"
	nullString := "null"
	intString := "42"

	// Text sucessful unmarshaling
	o := optional.NoBool()
	err := o.UnmarshalText([]byte(trueString))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, true, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullString))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// Test unmarshaling non-bool
	err = o.UnmarshalText([]byte(intString))
	assert.Assert(t, err != nil)
}

func TestBoolTrue(t *testing.T) {
	o := optional.NoBool()
	assert.Assert(t, !o.True())

	o = optional.SomeBool(false)
	assert.Assert(t, !o.True())

	o = optional.SomeBool(true)
	assert.Assert(t, o.True())
}

func TestBoolSql(t *testing.T) {
	ins := "INSERT INTO boolTest (state) VALUES (?)"
	q := `SELECT * FROM boolTest WHERE state = ?`
	test1Val := true
	test1 := optional.SomeBool(test1Val)
	test2 := optional.SomeBool(false)
	test3 := optional.NoBool()
	out1 := optional.NoBool()
	out2 := optional.NoBool()
	out3 := optional.SomeBool(true)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	assert.NilError(t, err, "failed to open mock database connection")
	defer db.Close()

	// mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO boolTest").WithArgs(test1).WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectExec("INSERT INTO boolTest").WithArgs(nil).WillReturnResult(sqlmock.NewResult(2, 1))
	r1 := sqlmock.NewRows([]string{"status"})
	r1.AddRow(true)
	select1 := mock.ExpectQuery(`SELECT (.+) FROM boolTest`)
	select1.WithArgs(test1).WillReturnRows(r1)
	r2 := sqlmock.NewRows([]string{"status"})
	r2.AddRow(false)
	select2 := mock.ExpectQuery(`SELECT (.+) FROM boolTest`)
	select2.WithArgs(test2).WillReturnRows(r2)
	r3 := sqlmock.NewRows([]string{"status"})
	r3.AddRow(nil)
	select3 := mock.ExpectQuery(`SELECT (.+) FROM boolTest`)
	select3.WithArgs(test3).WillReturnRows(r3)

	_, err = db.Exec(ins, test1)
	assert.NilError(t, err, "error using mock to insert optional type")

	// true case
	rows1, err := db.Query(q, test1)
	assert.NilError(t, err, "error using mock to query with optional type")
	defer rows1.Close()

	for rows1.Next() {
		err = rows1.Scan(&out1)
		assert.NilError(t, err, "error using Scan to convert sql row to optional type")
		assert.Equal(t, test1, out1, "Some(true) case failed")
	}

	// false case
	rows2, err := db.Query(q, test2)
	assert.NilError(t, err, "error using mock to query with optional type")
	defer rows2.Close()

	for rows2.Next() {
		err = rows2.Scan(&out2)
		assert.NilError(t, err, "error using Scan to convert sql row to optional type")
		assert.Equal(t, test2, out2, "Some(false) case failed")
	}

	// None case
	rows3, err := db.Query(q, test3)
	assert.NilError(t, err, "error using mock to query with optional type")
	defer rows3.Close()

	for rows3.Next() {
		err = rows3.Scan(&out3)
		assert.NilError(t, err, "error using Scan to convert sql row to optional type")
		assert.Equal(t, test3.IsNone(), out3.IsNone(), "None / NULL case failed")
	}
}
