package crypto

import (
	"testing"
)

func TestNew(t *testing.T) {
	fixtures := []struct {
		t  string
		fn func(Client) bool
	}{
		{
			t: "",
			fn: func(c Client) bool {
				_, ok := c.(*ChaCha20Poly1305Client)
				return ok
			},
		},
		{
			t: "asdf",
			fn: func(c Client) bool {
				_, ok := c.(*ChaCha20Poly1305Client)
				return ok
			},
		},
		{
			t: ConfigTypeChaChaPoly,
			fn: func(c Client) bool {
				_, ok := c.(*ChaCha20Poly1305Client)
				return ok
			},
		},
		{
			t: ConfigTypeOpenPGP,
			fn: func(c Client) bool {
				_, ok := c.(*OpenPGPClient)
				return ok
			},
		},
		{
			t: ConfigTypeAESGCM,
			fn: func(c Client) bool {
				_, ok := c.(*AESGCMClient)
				return ok
			},
		},
	}

	for _, tt := range fixtures {
		t.Run(tt.t, func(t *testing.T) {
			cfg := NewDefaultConfig()
			cfg.Type = tt.t

			if client, err := New(&cfg); err != nil {
				t.Fatal(err)
			} else {
				if !tt.fn(client) {
					t.Errorf("Failed to assert type")
				}
			}
		})
	}
}
