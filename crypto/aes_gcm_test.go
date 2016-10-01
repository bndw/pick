package crypto

import (
	"testing"
)

func TestAESGCMDefaultSettings(t *testing.T) {
	defaultSettings := DefaultAESGCMSettings()
	if defaultSettings == nil {
		t.Errorf("Expected default settings, got nil")
	}

	if defaultSettings.KeyLen != aesGCMDefaultKeyLen {
		t.Errorf("KeyLen is unexpected, expected %d, actual %d",
			aesGCMDefaultKeyLen, defaultSettings.KeyLen)
	}

	if defaultSettings.KeyDerivation != aesGCMDefaultKeyDerivation {
		t.Errorf("Default Key Derivation is unexpected. Expected %s, actual %s",
			aesGCMDefaultKeyDerivation, defaultSettings.KeyDerivation)
	}

	if defaultSettings.PBKDF2 == nil {
		t.Errorf("Expected a default PBKDF2")
	}

	if defaultSettings.Scrypt == nil {
		t.Errorf("Expected a default Scrypt")
	}
}

func TestAESGCMNewClient(t *testing.T) {
	if _, err := NewAESGCMClient(DefaultAESGCMSettings()); err != nil {
		t.Error(err)
	}
}

func TestAESGCMKeyLen(t *testing.T) {
	c, err := NewAESGCMClient(DefaultAESGCMSettings())
	if err != nil {
		t.Error(err)
	}

	if c.keyLen() == 0 {
		t.Errorf("Expected keylen > 0")
	}
}

func TestAESGCMDeriveKey(t *testing.T) {
	c, err := NewAESGCMClient(DefaultAESGCMSettings())
	if err != nil {
		t.Error(err)
	}

	_, _, err = c.deriveKey([]byte("d3adb33f"), c.keyLen())
	if err != nil {
		t.Error(err)
	}
}

func TestAESGCMDeriveKeyWithSalt(t *testing.T) {
	c, err := NewAESGCMClient(DefaultAESGCMSettings())
	if err != nil {
		t.Error(err)
	}

	_, err = c.deriveKeyWithSalt([]byte("d3adb33f"), []byte("pepper"), c.keyLen())
	if err != nil {
		t.Error(err)
	}
}

func TestAESGCMDecrypt(t *testing.T) {
	const ciphertext = `{"salt":"5hoAENJ0RMEa5Mu4Il8AsA==","nonce":"rgTzBsZfX/sBxDcW","ciphertext":"A+EU28usBw82UfagOQdlyPnmdX/CCR7dvoP1+Ff8NVXS7jStZvaLIJilYzI8"}`

	c, err := NewAESGCMClient(DefaultAESGCMSettings())
	if err != nil {
		t.Error(err)
	}

	_, err = c.Decrypt([]byte(ciphertext), []byte("d3adb33f"))
	if err != nil {
		t.Error(err)
	}
}

func TestAESGCMEncrypt(t *testing.T) {
	c, err := NewAESGCMClient(DefaultAESGCMSettings())
	if err != nil {
		t.Error(err)
	}

	_, err = c.Encrypt([]byte("Salad, it's what's for dinner"), []byte("d3adb33f"))
	if err != nil {
		t.Error(err)
	}
}
