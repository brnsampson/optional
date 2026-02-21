package optional_test

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
)

func TestSecretType(t *testing.T) {
	o := optional.SomeStr("A dumb test string")
	s := optional.MakeSecret(&o)
	assert.Equal(t, reflect.TypeOf(s).Name(), s.Type())
	assert.Assert(t, o.IsNone())
}

func TestSecretString(t *testing.T) {
	str := "***REDACTED***"
	o := optional.SomeStr(str)
	s := optional.MakeSecret(&o)
	assert.Equal(t, str, s.String())
	assert.Assert(t, o.IsNone())
}

func TestSecretMarshalText(t *testing.T) {
	str := "testing this tester with the testing module"
	o := optional.SomeSecret(str)
	s, err := o.MarshalText()

	assert.NilError(t, err)
	assert.Equal(t, str, string(s))
}

func TestSecretUnmarshalText(t *testing.T) {
	str := "testing this tester with the testing module"
	nullStr := "null"
	intStr := "42"

	// Text sucessful unmarshaling
	o := optional.NoSecret()
	err := o.UnmarshalText([]byte(str))
	assert.NilError(t, err)
	assert.Assert(t, o.IsSome())

	ret, ok := o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, str, ret)

	// Test unmarshaling null
	err = o.UnmarshalText([]byte(nullStr))
	assert.NilError(t, err)
	assert.Assert(t, o.IsNone())

	// There is no string that we cannot unmarshal into a string, but we should check that other types actually
	// end up as the string version as expected I guess...
	err = o.UnmarshalText([]byte(intStr))
	assert.NilError(t, err)

	ret, ok = o.Get()
	assert.Assert(t, ok)
	assert.Equal(t, intStr, ret)
}

func TestSecretSql(t *testing.T) {
	ins := "INSERT INTO optionTest (val) VALUES (?)"
	q := `SELECT * FROM optionTest WHERE val = ?`
	test1Val := "this is a secret!"
	test1 := optional.SomeSecret(test1Val)
	test2 := optional.NoSecret()
	out1 := optional.NoSecret()
	out2 := optional.SomeSecret("This is an old secret")

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
