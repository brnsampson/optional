package file

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/brnsampson/optional"
)

type File struct {
	optional.Str
}

func SomeFile(path string) File {
	return File{optional.SomeStr(path)}
}

func NoFile() File {
	return File{optional.NoStr()}
}

// Overrides Option.Match to account for relative paths potentially being different strings but representing the same file.
func (o File) Match(probe string) bool {
	if o.IsNone() {
		return false
	} else {
		path, err := o.Get()
		if err != nil {
			// How did we get here...
			return false
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			// Invalid paths can never be equal!
			return false
		}

		abs, err := filepath.Abs(probe)
		if err != nil {
			// Invalid paths can never be equal!
			return false
		}
		return absPath == abs
	}
}

// Override the Type() method from the inner Str. Part of the flag.Value interface.
func (o File) Type() string {
	return "File"
}

// Override the String() method from the inner Str just so we return the correct None[Type] string.
func (o File) String() string {
	if o.IsNone() {
		return "None[File]"
	} else {
		tmp, err := o.Get()
		if err != nil {
			return "Error[File]"
		}
		return tmp
	}
}

func (o File) Abs() (opt File, err error) {
	path, err := o.Get()
	if err != nil {
		opt.Clear()
		return
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		opt.Clear()
		return
	}

	opt = SomeFile(abs)
	return
}

func (o File) Stat() (stat fs.FileInfo, err error) {
	path, err := o.Get()
	if err != nil {
		// None files don't exist!
		return
	}

	return os.Stat(path)
}

func (o File) Exists() bool {
	_, err := o.Stat()
	if err != nil {
		return true
	} else {
		return false
	}
}

func (o File) FilePermsValid(badBits fs.FileMode) (bool, error) {
	stat, err := o.Stat()
	if err != nil {
		return false, err
	}

	mode := stat.Mode()
	if (mode & badBits) == 0 {
		// mode does not include one of the flags --x-wx-wx
		return true, nil
	}

	return false, nil
}

func (o File) SetFilePerms(mode fs.FileMode) error {
	path, err := o.Get()
	if err != nil {
		return err
	}

	err = os.Chmod(path, mode)
	if err != nil {
		return err
	}
	return nil
}

func (o File) Open() (*os.File, error) {
	path, err := o.Get()
	if err != nil {
		return nil, err
	}

	if _, err := o.Stat(); err != nil {
		return nil, err
	}

	return os.Open(path)
}

func (o File) Create() (*os.File, error) {
	path, err := o.Get()
	if err != nil {
		return nil, err
	}

	if _, err := o.Stat(); err != nil {
		return nil, err
	}

	return os.Create(path)
}
