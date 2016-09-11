package backends

// mockBackend is used for tests
type mockBackend struct {
	data []byte
}

func NewMockBackend() *mockBackend {
	safeData := []byte(`-----BEGIN PGP MESSAGE-----
wx4EBwMI/EyvqWA12cNgJBnoGRxYO1D0/F/w5Ro5uafS4AHkLjgl3wFVjIRB1vbo
GSX6FeE9q+Ap4JzhoTTgcOLB6iyW4HDmGZFzcVq+JgYYg0+7Q+4jlC/bBxyhtb1h
UHBuCvFGG4ENExdLliCsixI1bP8KB2TlLH459U859KWkg1aEJJ+1FeDR5E1GwV5y
Jn766KqjJFAUxwvguuNHI0fMMcIyfeA+4uNDsmXg+uRsGhwVdCP509FRtqes0EPh
4mqkkV7hFAgA=geI2
-----END PGP MESSAGE-----`)

	return &mockBackend{data: safeData}
}

func (b *mockBackend) Backup() error {
	return nil
}

func (b *mockBackend) Load() ([]byte, error) {
	return b.data, nil
}

func (b *mockBackend) Save(ciphertext []byte) error {
	b.data = ciphertext
	return nil
}
