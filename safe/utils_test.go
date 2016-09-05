package safe

import (
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
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

func createTestSafe() string {
	err := ioutil.WriteFile(testSafeName, []byte(testSafeContent), 0600)
	if err != nil {
		panic(err)
	}

	return testSafeName
}

func removeTestSafe(testSafe string) {
	_ = os.Remove(testSafe)
	return
}

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
