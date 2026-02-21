package optional

import (
	"database/sql/driver"
	"fmt"
	"log/slog"
)

// Secret wraps another Optional which may hold sensitive
// information. The value may still be marshaled into messages,
// but when printing to the console or logs it should be redacted.
//
// TODO: We should make this automatically encrypt secrets at rest.
type Secret struct {
	Str
}

// MakeSecret creates a Secret from a pointer to Str optional then clears
// the original. This is to ensure the original optional is never
// logged, thereby leaking the secret.
func MakeSecret(str *Str) Secret {
	var tmp Str
	val, ok := str.Get()
	if ok {
		tmp = SomeStr(val)
		str.Clear()
	}

	return Secret{tmp}
}

func SomeSecret(value string) Secret {
	return Secret{SomeStr(value)}
}

func NoSecret() Secret {
	return Secret{NoStr()}
}

// Type overrides the Type() method from the inner string. Part of the flag.Value interface.
func (s Secret) Type() string {
	return "Secret"
}

// String ALWAYS redact secrets
func (s Secret) String() string {
	return "***REDACTED***"
}

// Format ALWAYS ALWAYS redact secrets no matter what formatting verb or flag is set
func (s Secret) Format(f fmt.State, verb rune) {
	// We don't really care if there is an error or not I think...
	_, _ = f.Write([]byte(s.String()))
}

// LogValue even redacts secrets when logging. What a surprise!
func (s Secret) LogValue() slog.Value {
	return slog.StringValue(s.String())
}

// Scan implements database/sql.Scanner interface.
func (s *Secret) Scan(src any) error {
	// TODO: encode with a secret key before writing to DB
	if src == nil {
		// NULL value row
		s.Clear()
		return nil
	}
	switch t := src.(type) {
	case string:
		_ = s.Replace(t)
	case []byte:
		_ = s.Replace(string(t))
	default:
		return fmt.Errorf("converting driver.Value type %T to %s", src, s.Type())
	}
	return nil
}

// Value implements the database/sql/driver.Valuer interface
func (s Secret) Value() (driver.Value, error) {
	// TODO: decode with secret key after receiving from the DB
	val, ok := s.Get()
	if ok {
		return val, nil
	}
	return nil, nil
}
