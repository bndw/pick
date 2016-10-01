package crypto

import (
	"testing"
)

func TestOpenPGPDefaultSettings(t *testing.T) {
	defaultSettings := DefaultOpenPGPSettings()
	if defaultSettings == nil {
		t.Errorf("Expected default settings, got nil")
	}

	if defaultSettings.Cipher != openpgpDefaultCipher {
		t.Errorf("DefaultCipher is unexpected. Expected %s, actual %s",
			openpgpDefaultCipher, defaultSettings.Cipher)
	}

	if defaultSettings.S2KCount != openpgpDefaultS2KCount {
		t.Errorf("DefaultS2KCount is unexpected. Expected %d, actual %d",
			openpgpDefaultS2KCount, defaultSettings.S2KCount)
	}
}

func TestOpenPGPNewClient(t *testing.T) {
	if _, err := NewOpenPGPClient(DefaultOpenPGPSettings()); err != nil {
		t.Error(err)
	}
}

func TestOpenPGPS2kCount(t *testing.T) {
	c, err := NewOpenPGPClient(DefaultOpenPGPSettings())
	if err != nil {
		t.Error(err)
	}

	// Currently supported values: 1024-65011712 (inclusive).
	if c.s2kCount() < 1024 {
		t.Errorf("Expected s2kCount greater than 1024")
	}
	if c.s2kCount() > 65011712 {
		t.Errorf("Expected s2kCount less than 65011712")
	}
}

func TestOpenPGPDecrypt(t *testing.T) {
	const ciphertext = `-----BEGIN PGP MESSAGE-----

wy4ECQMIlYyEIihH6mD/sYAV1dj5VBXNexlblCUJQK/Tgabg0b9hO+e+Rc3hzdrI
0uAB5As80DvzzhrOdyewfysw0ErhOXvgvOBJ4SBW4K7ixJJVLuC15FXMfpwqfJwZ
NchaPB+OnGbguON33ZY8LhArV+At4k3JkB3gleCt4JTkoDuNJlgpe7cfHo33NhIv
K+I87kKs4YINAA==
=DaHB
-----END PGP MESSAGE-----`

	c, err := NewOpenPGPClient(DefaultOpenPGPSettings())
	if err != nil {
		t.Error(err)
	}

	_, err = c.Decrypt([]byte(ciphertext), []byte("d3adb33f"))
	if err != nil {
		t.Error(err)
	}
}

func TestOpenPGPEncrypt(t *testing.T) {
	c, err := NewOpenPGPClient(DefaultOpenPGPSettings())
	if err != nil {
		t.Error(err)
	}

	_, err = c.Encrypt([]byte("Salad, it's what's for dinner"), []byte("d3adb33f"))
	if err != nil {
		t.Error(err)
	}
}
