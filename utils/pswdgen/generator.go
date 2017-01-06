package pswdgen

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

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
		chars = charsFull
	case 2:
		chars = charsAlphaNum
	case 1:
		chars = charsAlpha
	case 0:
		chars = charsNum
	}
	password, err := generateUsingAlphabet(chars, p.config.Length)
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

func generateUsingAlphabet(alphabet string, length int) (string, error) {
	buffer := make([]byte, length)
	max := big.NewInt(int64(len(alphabet)))

	var index int
	var err error
	for i := 0; i < length; i++ {
		index, err = randomInt(max)
		if err != nil {
			return "", err
		}

		buffer[i] = alphabet[index]
	}

	return string(buffer), nil
}

func randomInt(max *big.Int) (int, error) {
	rand, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}

	return int(rand.Int64()), nil
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
