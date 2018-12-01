package safe

import (
	"encoding/json"
	"testing"

	"github.com/bndw/pick/backends"
	mockBackend "github.com/bndw/pick/backends/mock"
	"github.com/bndw/pick/config"
	"github.com/bndw/pick/crypto"
)

func TestSaveWithDefaultCrypto(t *testing.T) {
	// Backend client
	backendConfig := backends.NewDefaultConfig()
	backendClient := mockBackend.NewForTesting(t, &backendConfig, true)
	backendClient.Data = []byte(testSafeContent)

	// Crypto client
	cryptoConfig := crypto.NewDefaultConfig()
	cryptoClient, err := crypto.New(&cryptoConfig)
	if err != nil {
		t.Fatal(err)
	}

	// Safe config
	config := &config.Config{
		Encryption: cryptoConfig,
		Storage:    backendConfig,
		Version:    "1.2.3",
	}

	s, err := New([]byte(password), backendClient, cryptoClient, config, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.save(); err != nil {
		t.Fatal(err)
	}

	var dto safeDTO
	if err := json.Unmarshal(backendClient.Data, &dto); err != nil {
		t.Fatal(err)
	}

	if dto.Config.Type != crypto.DefaultConfigType {
		t.Errorf("Expected safe to be encrypted with %q client, got %q", crypto.DefaultConfigType, dto.Config.Type)
	}
}
