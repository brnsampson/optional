package main

import (
	"fmt"
	"os"

	"github.com/brnsampson/optional/file"
)

func LoadingAndReadingFiles(path *string) {
	// However we got it, we either have or do not have a path. For our example, let's assume we loaded this from a
	// flag so we end up with a *string which could be nil

	f := file.NoFile()
	if path != nil {
		f = file.SomeFile(*path)
	}

	// Just read the contents of a file. Acts like os.ReadFile(path)
	data, err := f.ReadFile()
	if err != nil {
		fmt.Println("Failed to read from file:")
		fmt.Println(err)
	} else {
		fmt.Println("Got file contents:")
		fmt.Println(string(data))
	}

	// Open a file for reading. File.Open() works just like os.Open(path),
	// so the file is opend in ReadOnly mode.
	var opened *os.File
	opened, err = f.Open()
	if err != nil {
		fmt.Println("Failed to open file for reading:")
		fmt.Println(err)
		return
	}
	defer opened.Close()

	// Now use the file handle exactly as you would if you called os.Open()
}

func WritingAndDeletingFiles(path *string) {
	f := file.NoFile()
	if path != nil {
		f = file.SomeFile(*path)
	}

	// Delete a file. Works like os.Remove, but also returns an error if the path is still None
	err := f.Remove()
	if err != nil {
		fmt.Println("Failed to remove file:")
		fmt.Println(err)
	}

	// Write the contents of a file. Acts like os.WriteFile(path)
	data := []byte("Hello, World!")
	err = f.WriteFile(data, 0644)
	if err != nil {
		fmt.Println("Failed to write file:")
		fmt.Println(err)
	}

	// Open a file for read/write. File.Create() works like like os.Create(path), which means
	// calling this will either create a file or truncate an existing file. If you want to
	// append to a file, you must use File.OpenFile(os.O_RDWR|os.O_CREATE, 0644) in the same way
	// that would need to when calling os.OpenFile. See https://pkg.go.dev/os#OpenFile for details.
	var opened *os.File
	opened, err = f.Create()
	if err != nil {
		fmt.Println("Failed to open/create file:")
		fmt.Println(err)
		return
	}
	defer opened.Close()

	// Now use the file handle exactly as you would if you called os.Create(path)
	opened.Write(data)
}

func AdditionalFileTools(path *string) {
	f := file.NoFile()
	if path != nil {
		f = file.SomeFile(*path)
	}

	// Read back the path
	p, ok := f.Get()
	if ok {
		fmt.Println("Got path:")
		fmt.Println(p)
	} else {
		fmt.Println("No path given!")
		os.Exit(1)
	}

	// Check if the given path is the same as some other path, matching all equivalent absolute and relative paths.
	// In this case, check if the given path is equivalent to our working directory.
	if f.Match(".") {
		fmt.Println("We are operating on our working directory. Be careful!")
	} else {
		fmt.Println("We are not in our working directory. Go nuts!")
	}

	// Get a new optional with any relative path converted to absolute path (also ensuring it is a valid path)
	abs, err := f.Abs()
	if err != nil {
		fmt.Println("Could not convert path into absolute path. Is it a valid path?")
		os.Exit(1)
	}

	// Stat the file, or just check if it exists if you don't care about other file info
	if abs.Exists() {
		fmt.Println("The file exists!")
	}

	info, err := abs.Stat() // I don't care about the info
	if err != nil {
		fmt.Println("Could not stat the file")
	} else {
		fmt.Println("Got file info:")
		fmt.Println(info)
	}

	// Check that the file has permissions 700 and modify it if it does not
	valid, err := abs.FilePermsValid(0644)
	if err != nil {
		fmt.Println("Could not read file permissions!")
		os.Exit(1)
	}

	if !valid {
		err = abs.SetFilePerms(0644)
		if err != nil {
			fmt.Println("Failed to set file perms to 700")
			os.Exit(1)
		}
	}
}
