package safe

import (
	"io/ioutil"
	"os"

	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/crypto"
)

func init() {
	removeTestSafe()
}

const (
	// testSafeContent has one account in it, "foo".
	testSafeContent = `-----BEGIN PGP MESSAGE-----

wx4EBwMI/EyvqWA12cNgJBnoGRxYO1D0/F/w5Ro5uafS4AHkLjgl3wFVjIRB1vbo
GSX6FeE9q+Ap4JzhoTTgcOLB6iyW4HDmGZFzcVq+JgYYg0+7Q+4jlC/bBxyhtb1h
UHBuCvFGG4ENExdLliCsixI1bP8KB2TlLH459U859KWkg1aEJJ+1FeDR5E1GwV5y
Jn766KqjJFAUxwvguuNHI0fMMcIyfeA+4uNDsmXg+uRsGhwVdCP509FRtqes0EPh
4mqkkV7hFAgA=geI2
-----END PGP MESSAGE-----`
	testSafeName = "test.safe"
)

var (
	testSafePassword = []byte("seabreezes")
)

func createTestSafe() (*Safe, error) {
	err := ioutil.WriteFile(testSafeName, []byte(testSafeContent), 0600)
	if err != nil {
		return nil, err
	}

	backendClient, err := backends.New(backends.Config{
		Type: "file",
		Settings: map[string]interface{}{
			"path": testSafeName,
		},
	})
	if err != nil {
		return nil, err
	}

	cryptoClient, err := crypto.New(crypto.NewDefaultConfig())
	if err != nil {
		return nil, err
	}

	return Load(
		testSafePassword,
		backendClient,
		cryptoClient,
	)
}

func removeTestSafe() {
	_ = os.Remove(testSafeName)
	return
}
