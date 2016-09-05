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
	testSafeContent = `-----BEGIN PGP SIGNATURE-----

wx4EBwMIEp9YTGqqetBgwXj5+80vLdBMeIvLS49/GPDS4AHkVDeGlKy/bXalpl9M
q7YwLOEBo+DO4CTh9ujgVeKriVie4HjmoykOiH3REZLAJO7U3ejMIP4onlD0u6SJ
vD4U2ipsQIWhWaWfASmpl6T0Qq39tDqE8XZjFcTt/Btujfb5zoPpbeA84s53sEXg
suDC4K7k5IYOU66CG3XbOaafkrrfmeKCjJGG4aZyAA==
=Wzqv
-----END PGP SIGNATURE-----`
	testSafeName = "test.safe"
)

var (
	testSafePassword = []byte("seabreezes\n")
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

	cryptoClient, err := crypto.New(crypto.Config{
		Type:     "aes-openpgp",
		Settings: map[string]interface{}{},
	})
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
