package randomstring_test

import (
	"testing"

	"github.com/leonklingele/randomstring"
)

func TestRandomStringLength(t *testing.T) {
	ls := []int{1, 5, 10, 50, 100}
	dict := randomstring.CharsASCII

	for _, l := range ls {
		s, err := randomstring.Generate(l, dict)
		if err != nil {
			t.Error(err)
		}

		if len(s) != l {
			t.Errorf("unexpected length of generated string: want %d, got %d", l, len(s))
		}
	}
}

func TestRandomStringInvalidLength(t *testing.T) {
	ls := []int{-1, 0}

	for _, l := range ls {
		if _, err := randomstring.Generate(l, randomstring.CharsAlpha); err != randomstring.ErrInvalidLengthSpecified {
			t.Errorf("unexpected error for length %d: %v", l, err)
		}
	}
}
