package backends

// mockBackend is used for tests
type mockBackend struct {
	data []byte
}

func NewMockBackend() *mockBackend {
	safeData := []byte(`-----BEGIN PGP SIGNATURE-----

wx4EBwMIEp9YTGqqetBgwXj5+80vLdBMeIvLS49/GPDS4AHkVDeGlKy/bXalpl9M
q7YwLOEBo+DO4CTh9ujgVeKriVie4HjmoykOiH3REZLAJO7U3ejMIP4onlD0u6SJ
vD4U2ipsQIWhWaWfASmpl6T0Qq39tDqE8XZjFcTt/Btujfb5zoPpbeA84s53sEXg
suDC4K7k5IYOU66CG3XbOaafkrrfmeKCjJGG4aZyAA==
=Wzqv
-----END PGP SIGNATURE-----`)

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
