package crypto

import (
	"testing"
)

func TestChaCha20Poly1305DefaultSettings(t *testing.T) {
	defaultSettings := DefaultChaCha20Poly1305Settings()
	if defaultSettings == nil {
		t.Errorf("Expected default settings, got nil")
	}

	if defaultSettings.KeyDerivation != chaCha20Poly1305DefaultKeyDerivation {
		t.Errorf("Default Key Derivation is unexpected. Expected %s, actual %s",
			chaCha20Poly1305DefaultKeyDerivation, defaultSettings.KeyDerivation)
	}

	if defaultSettings.PBKDF2 == nil {
		t.Errorf("Expected a default PBKDF2")
	}

	if defaultSettings.Scrypt == nil {
		t.Errorf("Expected a default Scrypt")
	}
}

func TestChaCha20Poly1305NewClient(t *testing.T) {
	if _, err := NewChaCha20Poly1305Client(DefaultChaCha20Poly1305Settings()); err != nil {
		t.Error(err)
	}
}

func TestChaCha20Poly1305KeyLen(t *testing.T) {
	c, err := NewChaCha20Poly1305Client(DefaultChaCha20Poly1305Settings())
	if err != nil {
		t.Error(err)
	}

	if c.keyLen() == 0 {
		t.Errorf("Expected keylen > 0")
	}
}

func TestChaCha20Poly1305DeriveKey(t *testing.T) {
	c, err := NewChaCha20Poly1305Client(DefaultChaCha20Poly1305Settings())
	if err != nil {
		t.Error(err)
	}

	_, _, err = c.deriveKey([]byte("d3adb33f"), c.keyLen())
	if err != nil {
		t.Error(err)
	}
}

func TestChaCha20Poly1305DeriveKeyWithSalt(t *testing.T) {
	c, err := NewChaCha20Poly1305Client(DefaultChaCha20Poly1305Settings())
	if err != nil {
		t.Error(err)
	}

	_, err = c.deriveKeyWithSalt([]byte("d3adb33f"), []byte("pepper"), c.keyLen())
	if err != nil {
		t.Error(err)
	}
}

func TestChaCha20Poly1305Decrypt(t *testing.T) {
	const ciphertext = `{"salt":"bsmaQTOsrM1HdaeAYm0psQ==","nonce":"yGfM0R3WU6++O4IM","ciphertext":"wPMci7YAKaAipWacrcYlkfMYmeaqxw65UL7wM/6Fi010w1CxejdNkyU9zv61"}`

	c, err := NewChaCha20Poly1305Client(DefaultChaCha20Poly1305Settings())
	if err != nil {
		t.Error(err)
	}

	_, err = c.Decrypt([]byte(ciphertext), []byte("d3adb33f"))
	if err != nil {
		t.Error(err)
	}
}

func TestChaCha20Poly1305Encrypt(t *testing.T) {
	c, err := NewChaCha20Poly1305Client(DefaultChaCha20Poly1305Settings())
	if err != nil {
		t.Error(err)
	}

	_, err = c.Encrypt([]byte("Salad, it's what's for dinner"), []byte("d3adb33f"))
	if err != nil {
		t.Error(err)
	}
}
