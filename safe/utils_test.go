package safe

import (
	"testing"

	"github.com/bndw/pick/backends"
	mockBackend "github.com/bndw/pick/backends/mock"
	"github.com/bndw/pick/config"
	"github.com/bndw/pick/crypto"
)

const (
	testSafePassword = "seabreezes"
	// testSafeContent has one account in it, "foo".
	testSafeContent = `-----BEGIN PGP MESSAGE-----

wx4EBwMI/EyvqWA12cNgJBnoGRxYO1D0/F/w5Ro5uafS4AHkLjgl3wFVjIRB1vbo
GSX6FeE9q+Ap4JzhoTTgcOLB6iyW4HDmGZFzcVq+JgYYg0+7Q+4jlC/bBxyhtb1h
UHBuCvFGG4ENExdLliCsixI1bP8KB2TlLH459U859KWkg1aEJJ+1FeDR5E1GwV5y
Jn766KqjJFAUxwvguuNHI0fMMcIyfeA+4uNDsmXg+uRsGhwVdCP509FRtqes0EPh
4mqkkV7hFAgA=geI2
-----END PGP MESSAGE-----`
)

func createTestSafe(t *testing.T) (*Safe, error) {
	// t.Helper() // TOOD(leon): Go 1.9 only :(

	backendConfig := backends.NewDefaultConfig()
	backendClient := mockBackend.NewForTesting(t, &backendConfig)
	backendClient.Data = []byte(testSafeContent)

	cryptoConfig := crypto.Config{
		Type: crypto.ConfigTypeOpenPGP,
		OpenPGPSettings: &crypto.OpenPGPSettings{
			Cipher:   "aes128",
			S2KCount: 1024,
		},
	}

	cryptoClient, err := crypto.New(&cryptoConfig)
	if err != nil {
		return nil, err
	}

	config := &config.Config{
		Encryption: cryptoConfig,
		Storage:    backendConfig,
		Version:    "1.2.3",
	}

	return Load(
		[]byte(testSafePassword),
		backendClient,
		cryptoClient,
		config,
	)
}
