package pswdgen

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/leonklingele/randomstring"
	"github.com/pkg/term"
)

type passwordGenerator struct {
	config Config
}

func (p *passwordGenerator) IncreaseStrengh() {
	if p.config.Strength < 3 {
		p.config.Strength++
	}
}

func (p *passwordGenerator) DecreaseStrengh() {
	if p.config.Strength > 0 {
		p.config.Strength--
	}
}

func (p *passwordGenerator) IncreaseLength() {
	p.config.Length++
}

func (p *passwordGenerator) DecreaseLength() {
	if p.config.Length > 1 {
		p.config.Length--
	}
}

func (p *passwordGenerator) Generate() (string, error) {
	if p.config.Mode != passwordModeInteractive {
		return p.New()
	}
	fmt.Println("Entering interactive password generation mode")
	fmt.Println("- Increase character set via: Up-arrow key / \"k\"")
	fmt.Println("- Decrease character set via: Down-arrow key / \"j\"")
	fmt.Println("- Increase length via: Right-arrow key / \"l\"")
	fmt.Println("- Decrease length via: Left-arrow key / \"h\"")
	fmt.Println("- Use current password: Enter key")
	return p.NewInteractively()
}

func (p *passwordGenerator) New() (string, error) {
	var chars string
	switch p.config.Strength {
	default:
		fallthrough
	case 3:
		chars = randomstring.CharsASCII
	case 2:
		chars = randomstring.CharsAlphaNum
	case 1:
		chars = randomstring.CharsAlpha
	case 0:
		chars = randomstring.CharsNum
	}
	password, err := randomstring.Generate(p.config.Length, chars)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (p *passwordGenerator) NewInteractively() (password string, err error) {
	print := func(s string) {
		// Clear password
		fmt.Print("\r", strings.Repeat(" ", p.config.Length+1))
		fmt.Print("\r", s)
	}
	genNew := func() {
		password, err = p.New()
		if err != nil {
			fmt.Println("Failed to generate password", err)
			return
		}
		print(password)
	}
	genNew()
	for {
		c, readErr := readTermChar()
		if readErr != nil {
			continue
		}
		switch {
		case bytes.Equal(c, []byte{3}), bytes.Equal(c, []byte{13}): // exit
			print("Password generated\n")
			return
		case bytes.Equal(c, []byte{27, 91, 65}), bytes.Equal(c, []byte{107}): // up
			p.IncreaseStrengh()
		case bytes.Equal(c, []byte{27, 91, 66}), bytes.Equal(c, []byte{106}): // down
			p.DecreaseStrengh()
		case bytes.Equal(c, []byte{27, 91, 67}), bytes.Equal(c, []byte{108}): // right
			p.IncreaseLength()
		case bytes.Equal(c, []byte{27, 91, 68}), bytes.Equal(c, []byte{104}): // left
			p.DecreaseLength()
		}
		genNew()
	}
}

func newPasswordGenerator(config Config) passwordGenerator {
	return passwordGenerator{
		config: config,
	}
}

// Generate generates a password.
func Generate(config Config) (string, error) {
	pswdGen := newPasswordGenerator(config)
	return pswdGen.Generate()
}

func readTermChar() ([]byte, error) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)
	numRead, err := t.Read(bytes)
	t.Restore()
	t.Close()
	if err != nil {
		return nil, err
	}
	return bytes[0:numRead], nil
}
