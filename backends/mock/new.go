package mock

import (
	"testing"

	"github.com/bndw/pick/backends"
)

func _new(config *backends.Config) (backends.Client, error) {
	safeData := []byte(`-----BEGIN PGP MESSAGE-----

wx4EBwMI/EyvqWA12cNgJBnoGRxYO1D0/F/w5Ro5uafS4AHkLjgl3wFVjIRB1vbo
GSX6FeE9q+Ap4JzhoTTgcOLB6iyW4HDmGZFzcVq+JgYYg0+7Q+4jlC/bBxyhtb1h
UHBuCvFGG4ENExdLliCsixI1bP8KB2TlLH459U859KWkg1aEJJ+1FeDR5E1GwV5y
Jn766KqjJFAUxwvguuNHI0fMMcIyfeA+4uNDsmXg+uRsGhwVdCP509FRtqes0EPh
4mqkkV7hFAgA=geI2
-----END PGP MESSAGE-----`)

	return &client{Data: safeData}, nil
}

func NewForTesting(t *testing.T, config *backends.Config, writable bool) *client {
	// t.Helper() // TOOD(leon): Go 1.9 only :(
	if config == nil {
		tmp := backends.NewDefaultConfig()
		config = &tmp
	}
	// c, err := backends.NewWithType(clientName, config)
	// TODO(leon): This fails for whatever reason:
	// interface conversion: backends.Client is *mock.client, not *mock.client
	c, err := _new(config)
	if err != nil {
		t.Fatalf("Failed to create mock backend: %v", err)
	}
	c.SetWritable(writable)
	return c.(*client)
}
