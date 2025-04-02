package optional

import (
	"fmt"
	"log/slog"
)

// Secret wraps another Optional which may hold sensitive
// information. The value may still be marshaled into messages,
// but when printing to the console or logs it should be redacted.
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

// Override the Type() method from the inner string. Part of the flag.Value interface.
func (s Secret) Type() string {
	return "Secret"
}

// ALWAYS redact secrets
func (s Secret) String() string {
	return "***REDACTED***"
}

// ALWAYS ALWAYS redact secrets no matter what formatting verb or flag is set
func (s Secret) Format(f fmt.State, verb rune) {
	f.Write([]byte(s.String()))
}

// Even redact secrets when logging. What a surprise!
func (s Secret) LogValue() slog.Value {
	return slog.StringValue(s.String())
}
