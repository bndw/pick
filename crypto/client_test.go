package crypto

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	newDefaultCryptoClient := func() Client {
		c := NewDefaultConfig()
		dcc, err := New(&c)
		if err != nil {
			t.Fatal(err)
		}
		return dcc
	}

	fixtures := []struct {
		t  string
		fn func(Client) (Client, bool)
	}{
		{
			t: "",
			fn: func(c Client) (Client, bool) {
				dcc := newDefaultCryptoClient()
				return dcc, reflect.TypeOf(c) == reflect.TypeOf(dcc)
			},
		},
		{
			t: "asdf",
			fn: func(c Client) (Client, bool) {
				dcc := newDefaultCryptoClient()
				return dcc, reflect.TypeOf(c) == reflect.TypeOf(dcc)
			},
		},
		{
			t: ConfigTypeChaChaPoly,
			fn: func(c Client) (Client, bool) {
				c, ok := c.(*ChaCha20Poly1305Client)
				return c, ok
			},
		},
		{
			t: ConfigTypeOpenPGP,
			fn: func(c Client) (Client, bool) {
				c, ok := c.(*OpenPGPClient)
				return c, ok
			},
		},
		{
			t: ConfigTypeAESGCM,
			fn: func(c Client) (Client, bool) {
				c, ok := c.(*AESGCMClient)
				return c, ok
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
				if ac, ok := tt.fn(client); !ok {
					t.Errorf("Failed to assert type to %T", ac)
				}
			}
		})
	}
}
