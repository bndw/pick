# Secure temporary directories and files in Go

[![Build Status](https://travis-ci.org/leonklingele/securetemp.svg?branch=master)](https://travis-ci.org/leonklingele/securetemp)

This library can be used to create temporary directories and files inside RAM.
This is useful if you temporarily need a file to store sensitive data in it.

## tl;dr

```sh
# Install library
go get -u github.com/leonklingele/securetemp
```

```go
// Create a temporary file inside a RAM disk which is securetemp.DefaultSize big
tmpFile, cleanupFunc, err := securetemp.TempFile(securetemp.DefaultSize)
if err != nil {
	// TODO: Properly handle error
	log.Fatal(err)
}
// `tmpFile` is an *os.File
if _, err := tmpFile.WriteString("Hello, World!"); err != nil {
	// TODO: Properly handle error
	log.Fatal(err)
}
// Call the cleanup func as soon as you no longer need the file.
// The file is deleted and the RAM is freed again.
// The cleanup function also calls `tmpFile.Close()` for you.
cleanupFunc()
```

```go
// Create a temporary directory inside a RAM disk which is 20MB big
tmpDir, cleanupFunc, err := securetemp.TempDir(20 * securetemp.SizeMB)
if err != nil {
	// TODO: Properly handle error
	log.Fatal(err)
}
// Create one / multiple file/s inside `tmpDir`
file, err := os.Create(path.Join(tmpDir, "myfile.txt"))
if err != nil {
	// TODO: Properly handle error
	log.Fatal(err)
}
// Do something with `file`
if _, err := tmpFile.WriteString("Hello, World!"); err != nil {
	// TODO: Properly handle error
	log.Fatal(err)
}
// Call the cleanup func as soon as you no longer need the directory.
// The directory is deleted and the RAM is freed again.
cleanupFunc()
```

## Idea

This project was inspired by [pass](https://www.passwordstore.org/)
