# Cryptographically secure random strings in Go

[![Build Status](https://travis-ci.org/leonklingele/randomstring.svg?branch=master)](https://travis-ci.org/leonklingele/randomstring)

I was tired of so many Go apps and libraries being [modulo-biased](https://stackoverflow.com/a/10984975) when generating random strings.
Here's another library to generate cryptographically secure and unbiased strings.
Internally this library uses Go's [crypto/rand](https://golang.org/pkg/crypto/rand/) package which reads from the best random generator your system provides.

## tl;dr

```sh
# Install library
go get -u github.com/leonklingele/randomstring
```

```go
// .. and use it. Generate a 30 character long string
s, err := randomstring.Generate(30, randomstring.CharsASCII)
```

```go
// You can even use your own alphabet "abc". Be careful though!
s, err := randomstring.Generate(30, "abc")
```

Supported alphabets:

- `randomstring.CharsNum`: Contains numbers from 0-9
- `randomstring.CharsAlpha`: Contains the full English alphabet: letters a-z and A-Z
- `randomstring.CharsAlphaNum`: Is a combination of CharsNum and CharsAlpha
- `randomstring.CharsASCII`: Contains all printable ASCII characters in code range [32, 126]
