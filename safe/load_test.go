package safe_test

import (
	"testing"

	"github.com/bndw/pick/backends"
	mockBackend "github.com/bndw/pick/backends/mock"
	"github.com/bndw/pick/config"
	"github.com/bndw/pick/crypto"
	"github.com/bndw/pick/crypto/pbkdf2"
	"github.com/bndw/pick/crypto/scrypt"
	"github.com/bndw/pick/safe"
)

const (
	password = "seabreezes"
)

func TestLoadWithUpdatedConfig(t *testing.T) {
	// The safe defined in data was created with 100000 iterations. This test
	// asserts that when providing an updated config the safe will be updated
	// accordingly.
	const expectedIterations = 1024
	data := []byte(`{"config":{"type":"chachapoly","chachapoly":{"keyderivation":"pbkdf2","pbkdf2":{"hash":"sha512","iterations":100000,"saltlen":16}}},"ciphertext":"eyJzYWx0IjoiUmRMb1lwZ2RzbzNkMjFLWVBvc2ZVZz09Iiwibm9uY2UiOiJjNVVJTTgwZ2NUUEFXNytmIiwiY2lwaGVydGV4dCI6IndmY21Td1kwOWZ0UDNhVDhubjRLRlEzL3ZreXR3ZElHYWFZUThqa0s5NjhoUHhMWjViU21KTmd6YmVPTGRMTXIwdC9GaVVWK0gxdlpWY3ZSZzZnczJmWEljZUhVZXg0QmxBNjN3TkNaMVRxNVNZaEpvMG9nUXVua0xRZ0xVMEc0RFBML3V2bC9sVFRGcWtGd2tUUGhacitFYWp6ZDlaN2Z5QXg5TmpUQm1uM1B5SjhMb1R1Z3VWKzBiY3FJbzcxcmVXcmpUVmc3WWFEVlhlUXRycC92NTFyT0tSYzFtTmhCUkRrU202Q2ZMb2N4b0hEUnBWY3dDNVZzaDVKMTVrWWJXUzF2QlJvWVp4TGM4S1g0Nm1BbHJDZ3ZRYU51R0FCczVFbVd2WFg1ZVlQR0hwMWl5RTByWGVoRXV2YXFmVlg4OUF0T3Q3RUp2aTJMenBmR0tOSmFBN1hiWFlvZ1IyL0FBTTVvOVBxV3p1aUlqalJXYVRsTVJZVnJHUGc3UENRRWNteUZhWVRkdnZyZnRHbHpLdXlFdmk5VlJZY2NPS1A1cGFZSURwbTlEd0hTOTh1dlByNDdHeHJxdU9OOFhrREx2RXUzTGFrMk1pUFZjUGF0Y09LaTArRnFqV3QzZC8vc2NkMXNZMkpSVGFHTnhsVFUvZFhLNXJEV0xSTTlFT0FITlF2dE9oeHhTQXVSdmlKOUdwSm5FTXJram5ZdHFYZVhpcU4ybXdiTmx0b2VMdER1RjdOQUhxKzlCM3NUaC81ajVRMVhuWlRPN1BrTjFIQ1J5b0xKdHhIUnVIaWpUbHU2NG1JVFdEbTludDNQcDNmRzJBV3d2cHFqUGorV3dRUlZsanIyU3JtbGxsZ0pzZlVNL0xhQUNoV2dzTnV0M3dOTmJQRC9wVVFDL3k1d2hqcFlzZEIzQ3Y4QlZQY3dOR21ZK2hFMXEyZEZkRGs3SklqaS9JeG0yUTkwbnBjUnBDa29XS3lwUzNPQmtweW9QWGZMa1J5WWJjNkFNNGdSZWM2MGhxMC8ifQ=="}
	`)
	cryptoConfig := crypto.Config{
		Type: crypto.ConfigTypeChaChaPoly,
		ChaCha20Poly1305Settings: &crypto.ChaCha20Poly1305Settings{
			KeyDerivation: "pbkdf2",
			PBKDF2: &pbkdf2.PBKDF2{
				Hash:       "sha512",
				Iterations: expectedIterations,
				SaltLen:    16,
			},
		},
	}
	cryptoClient, err := crypto.New(&cryptoConfig)
	if err != nil {
		t.Fatal(err)
	}

	conf := &config.Config{
		Encryption: cryptoConfig,
		Storage:    backends.NewDefaultConfig(),
		Version:    "0.6.0",
	}
	backendClient := mockBackend.NewForTesting(t, &conf.Storage)
	backendClient.Data = data

	s, err := safe.Load([]byte(password), backendClient, cryptoClient, conf)
	if err != nil {
		t.Fatalf("Failed to load safe: %s", err.Error())
	}
	if _, err := s.Get("test"); err != nil {
		t.Fatalf("Failed to get 'test' safe account: %s", err.Error())
	}

	// Make sure the safe's iterations have been updated
	dto := safe.NewSafeDTO(backendClient.Data)
	if dto.Config.ChaCha20Poly1305Settings.PBKDF2.Iterations != expectedIterations {
		t.Fatalf("Expected safe to be upgraded when config changed")
	}
}

func TestLoad(t *testing.T) {
	backendConfig := backends.NewDefaultConfig()
	backendClient := mockBackend.NewForTesting(t, &backendConfig)
	for i, cryptoConfig := range safeCryptoConfigs {
		cryptoClient, err := crypto.New(&cryptoConfig)
		if err != nil {
			t.Fatalf("#%d: failed to create crypto client: %s", i, err.Error())
		}
		conf := &config.Config{
			Encryption: cryptoConfig,
			Storage:    backendConfig,
			Version:    "1.2.3",
		}
		for j, data := range safeData {
			backendClient.Data = []byte(data)
			s, err := safe.Load([]byte(password), backendClient, cryptoClient, conf)
			if err != nil {
				t.Fatalf("#%d.#%d: failed to load safe: %s", i, j, err.Error())
			}
			if _, err := s.Get("foo"); err != nil {
				t.Fatalf("#%d.#%d: failed to get 'foo' safe account: %s", i, j, err.Error())
			}
		}
	}
}

var safeCryptoConfigs = []crypto.Config{
	crypto.NewDefaultConfigWithType(crypto.ConfigTypeOpenPGP),
	{
		Type: crypto.ConfigTypeOpenPGP,
		OpenPGPSettings: &crypto.OpenPGPSettings{
			Cipher:   "aes128",
			S2KCount: 1024,
		},
	},
	{
		Type: crypto.ConfigTypeOpenPGP,
		OpenPGPSettings: &crypto.OpenPGPSettings{
			Cipher:   "aes256",
			S2KCount: 1024,
		},
	},
	{
		Type: crypto.ConfigTypeOpenPGP,
		OpenPGPSettings: &crypto.OpenPGPSettings{
			Cipher:   "does-not-exist",
			S2KCount: 1024,
		},
	},
	{
		Type: crypto.ConfigTypeOpenPGP,
		OpenPGPSettings: &crypto.OpenPGPSettings{
			Cipher:   "aes128",
			S2KCount: 4096,
		},
	},
	{
		Type: crypto.ConfigTypeOpenPGP,
		OpenPGPSettings: &crypto.OpenPGPSettings{
			Cipher:   "aes256",
			S2KCount: 4096,
		},
	},
	{
		Type: crypto.ConfigTypeOpenPGP,
		OpenPGPSettings: &crypto.OpenPGPSettings{
			Cipher:   "does-not-exist",
			S2KCount: 4096,
		},
	},
	crypto.NewDefaultConfigWithType(crypto.ConfigTypeAESGCM),
	{
		Type: crypto.ConfigTypeAESGCM,
		AESGCMSettings: &crypto.AESGCMSettings{
			KeyLen:           16,
			KeyDerivation:    "pbkdf2",
			Pbkdf2Hash:       "sha256",
			Pbkdf2Iterations: 1024,
			Pbkdf2SaltLen:    16,
		},
	},
	{
		Type: crypto.ConfigTypeAESGCM,
		AESGCMSettings: &crypto.AESGCMSettings{
			KeyLen:           32,
			KeyDerivation:    "pbkdf2",
			Pbkdf2Hash:       "sha512",
			Pbkdf2Iterations: 4096,
			Pbkdf2SaltLen:    32,
		},
	},
	{
		Type: crypto.ConfigTypeAESGCM,
		AESGCMSettings: &crypto.AESGCMSettings{
			KeyLen:        16,
			KeyDerivation: "pbkdf2",
			PBKDF2: &pbkdf2.PBKDF2{
				Hash:       "sha256",
				Iterations: 1024,
				SaltLen:    16,
			},
		},
	},
	{
		Type: crypto.ConfigTypeAESGCM,
		AESGCMSettings: &crypto.AESGCMSettings{
			KeyLen:        32,
			KeyDerivation: "pbkdf2",
			PBKDF2: &pbkdf2.PBKDF2{
				Hash:       "sha512",
				Iterations: 4096,
				SaltLen:    32,
			},
		},
	},
	{
		Type: crypto.ConfigTypeAESGCM,
		AESGCMSettings: &crypto.AESGCMSettings{
			KeyLen:        16,
			KeyDerivation: "scrypt",
			Scrypt: &scrypt.Scrypt{
				SaltLen: 16,
				N:       4096,
				R:       1,
				P:       1,
			},
		},
	},
	{
		Type: crypto.ConfigTypeAESGCM,
		AESGCMSettings: &crypto.AESGCMSettings{
			KeyLen:        32,
			KeyDerivation: "scrypt",
			Scrypt: &scrypt.Scrypt{
				SaltLen: 32,
				N:       8192,
				R:       2,
				P:       2,
			},
		},
	},
	crypto.NewDefaultConfigWithType(crypto.ConfigTypeChaChaPoly),
	{
		Type: crypto.ConfigTypeChaChaPoly,
		ChaCha20Poly1305Settings: &crypto.ChaCha20Poly1305Settings{
			KeyDerivation: "pbkdf2",
			PBKDF2: &pbkdf2.PBKDF2{
				Hash:       "sha256",
				Iterations: 1024,
				SaltLen:    16,
			},
		},
	},
	{
		Type: crypto.ConfigTypeChaChaPoly,
		ChaCha20Poly1305Settings: &crypto.ChaCha20Poly1305Settings{
			KeyDerivation: "pbkdf2",
			PBKDF2: &pbkdf2.PBKDF2{
				Hash:       "sha512",
				Iterations: 4096,
				SaltLen:    32,
			},
		},
	},
	{
		Type: crypto.ConfigTypeChaChaPoly,
		ChaCha20Poly1305Settings: &crypto.ChaCha20Poly1305Settings{
			KeyDerivation: "scrypt",
			Scrypt: &scrypt.Scrypt{
				SaltLen: 16,
				N:       4096,
				R:       1,
				P:       1,
			},
		},
	},
	{
		Type: crypto.ConfigTypeChaChaPoly,
		ChaCha20Poly1305Settings: &crypto.ChaCha20Poly1305Settings{
			KeyDerivation: "scrypt",
			Scrypt: &scrypt.Scrypt{
				SaltLen: 32,
				N:       8192,
				R:       2,
				P:       2,
			},
		},
	},
}

var safeData = []string{
	`-----BEGIN PGP MESSAGE-----

wx4EBwMI/EyvqWA12cNgJBnoGRxYO1D0/F/w5Ro5uafS4AHkLjgl3wFVjIRB1vbo
GSX6FeE9q+Ap4JzhoTTgcOLB6iyW4HDmGZFzcVq+JgYYg0+7Q+4jlC/bBxyhtb1h
UHBuCvFGG4ENExdLliCsixI1bP8KB2TlLH459U859KWkg1aEJJ+1FeDR5E1GwV5y
Jn766KqjJFAUxwvguuNHI0fMMcIyfeA+4uNDsmXg+uRsGhwVdCP509FRtqes0EPh
4mqkkV7hFAgA=geI2
-----END PGP MESSAGE-----`,
	`{"config":{"type":"openpgp","openpgp":{"cipher":"aes256","s2kcount":65011712}},"ciphertext":"LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eTRFQ1FNSVdlWTY2VzNNRUozL2FYM29qTmt3ZjVqR3F4UHB4R0J3STV0bUVJcmJQc2k0QXA3RVVPc1J1TUpUCjB1QUI1QWFsNVMwYiswUXk0b1huRzZSUmZxTGgva1RnZ2VDWDRVd1E0RHppZFkrWkErQ1g2QnRScWZPWFdyMy8KejNEUVVWM0NqRzRSVkc5U0Y0bWlpb0NnaGV2VU9KditOcGxLUHkzMGszNm5rZzJQZ253ZlpwZmlITEYyUHdvMQpuYmwwRGhvT2NaU0JUS281dVptZ3JFWXdzbXhmKzZxVVVZL1pnaG1zZU5tYzhsSGJTaHJVUE5PekoxZ3lZdWMzCkk5eUtjQ21qNjFKSVVyQzBiVU9GaGcwR3RvMmVTck1TRzNXYlVvVlVHOTBmWUVKZ3VQTjBzMFB0aFBuSllFbmQKaERwUWpUR3gzQVFtUWpWM3I1djJaOWorTTFmeHh5MDd1T1dxYUhQak5ZS3d1VGJhY0xHQkVXVFQvSmxQcXpwagpiSTlmb1IybVBiZkpuV0t1VmNBN0E5TUsrQi9wbWpiZXJhVmpRaWYvNHpHRXp0OG1pRFpEV0NqcDMvRnp5NGpqCldrT1dFUS90aTh6Z2pPYUErMFBvdXJ4WjlKSkt0YVcvck5zT0Ftb0lKU2J3ZkFablhzcU9KVU5iQjVHMk01dnEKdzJyN3djdUNjT1RlTm44YVQ3V1FIcGo5cVhVbkgyOGNpdWlRNEQzaWpQUVV0T0JyNElyZy91UkNhWCtWRDBWbQo0TEZNWDNHUDhXY3Q0c1FLVG5IaHZjQUEKPXo2cEIKLS0tLS1FTkQgUEdQIE1FU1NBR0UtLS0tLQ=="}`,
	`{"config":{"type":"openpgp","openpgp":{"cipher":"aes128","s2kcount":1024}},"ciphertext":"LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eDRFQndNSUxwM2RYQ1lzSzZnQVdWMmF6bWlISVlKRTcrTFpNSlRVR1lYUzRBSGs2YmNEM2hzblZwS041ejhrClAxWS92dUU4SXVDazRKSGg0cERnVU9JMVZUV0o0RW5vbjB5Um53dlFtdytMbG5vUllzWVpGWGZMdEJGWDMwZVEKUUNMWVlITi9US0RZSytiS2R3V3hxNE02NkVOMEo5SmRrd2ozczJKNHlJUXE4TE50MnRhSjhXcVNteTlVNVBpUgpNQ3l3VDNDT3gvcURxTXF5UVpxU1hadHBaSFBUUEtPRGJISjdOekpmSm5OMDdDVCt3eWdFcjQwUjE5cXQ1TXNCCnF4NjkyTWorcDlFZnFIMStuZ3A5TVdUalVYMmRXVWl4SG5OaHdDVkovSTJsT2FKbFI3dUJCR2lCd0JBakFadGMKVVdPZVF5d0FWa0xSM0YrSWFkZkxjSFhhR3RCMFVkUzU2WGRvQng2S3Y0aEdiWWNJbjMweXhQQTd1YVJ5czNQTwpBeloxRlZrL0R3azlxd0dJanNiQ0UyWlJVOHo1VGJTLzA0K3Blc0pnTW1uWWRaYXN1RGFoZE9EVTVyVHQ0YXZBCmxBZVRwakdlSklwY1VGaWhmS2lxaG1rUFVYTXFLdDIvY0ZvK09tanF0VGZBZ2E2TVZwaUVuaW1CQWwvK0pHN1YKZW1sdUtPY0k4MWdHeXgvZ09lQk80Q3JrQytYVlAvOVQvWnpwZ3VZdkFlUC91T0l5THFUazRjb05BQT09Cj1ua2ZsCi0tLS0tRU5EIFBHUCBNRVNTQUdFLS0tLS0="}`,
	`{"config":{"type":"openpgp","openpgp":{"cipher":"aes256","s2kcount":1024}},"ciphertext":"LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eTRFQ1FNSThNazl4SVhUNzE4QW9Mb3oyM25VS0ZYWk9IWk81eW1yMjlXTXhsV25iaXhaY0dyUDkwODREWS8vCjB1QUI1SHUzZzRhT01WeXVCUTE2N3pPQ0hQcmhTb2pnbHVCYjRYdXE0QWJpL0lhc3d1RDk2THRJQWN3L2VjblUKZG1FVGlXd05YNDNINk14SkxqVi9TNitTcEZTT053R2RNZzU2azQ0Q3YwRlA4b3JvQ1NDQWZFbEpaaG5qekU3dAp1dTBKRWFZdStOMnRRK0YxeWhKeFhkQ1c2R2dzRlJpQ1NUYmRSbThJZUFqSko3bjBhdE5VN2kvOHJOZEZoY2NQClU4UldvWFNiUGl1ZHg3MVNvYTZLV2NMZWR1Z3hqQmFKMXBiZUsxSDdNVEhVTGg5SUtVTHdRdjhXeVNiL1lTb0wKL3ZDdWFScTVPcjFzNWxCQVV2QjNIdmtnQnlBYjNPTnpGb3JxZi9qQ2p3NU9uLzFJNjBKQmR0WnhIOW5LVlp0VwpJZk9JeG40WjhvbTh4YUVRVmZGenlrdU5jMDJ2UEtvNDVZdU51RDlxY2ZxQVRVK3FyUGwycnhoNlU3YkxCamNJCjR0eXVETnJoYmJIZysrYkFYbWV4QThmYWk5LzF2Z2Y3Z0g2SUlUNXd5MEFPMDlCSkxqS1hBQ09ITkJrT1V6cGwKTXg4Nk9ISmt5Rk5nOWMzMFpsdmEwcnVqU3N1YnUvUStPRkZQNE52Z0V1Q0E1SjB6dkFsMHYwelZZVHQ1SzBsZwp0T1BpS29YMTUrRm9WQUE9Cj1jc1VvCi0tLS0tRU5EIFBHUCBNRVNTQUdFLS0tLS0="}`,
	`{"config":{"type":"openpgp","openpgp":{"cipher":"does-not-exist","s2kcount":1024}},"ciphertext":"LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eTRFQ1FNSVpLY202YVJGZ0JBQWR5bW1laUZZOW5MV0hYK1luL3NlRm5DSWJjc095NGV6TzRuSS9NOElYbDhECjB1QUI1STFudi9rMUdYRjVNUXZVbmxWVXZBbmhodm5naE9ERTRYZHI0TzNpVzBSNkVlRFo2Tklnb29DVGREbVoKR1lQeExnK0oySDFjaTlCaUQ3WjhRMnJWdFJra3QvTDR3WkY5bVBLNE9uTzhLMUkxdEQxVkhxRGdVMDA3M1F0awplZzlvYXdJTTFpcnViWncxYnNMZUFnNitlZXZhQXdnY3FtaHVqTGpJV3pkL3ZuS2xaRTdoc2lQQ1JvN0pKQzNoCjdxZVpSaHFuWkFMZUxpdG1zb25aVUtmRjcyejF2Z1lBbEljRnlScmp3b3NiQmIveS9nZXA2cFd0Z2dHQUNNbnoKU2FhZFYxM0w4eFV4SjZta0t1RHJTdS9BeGp6SldWWEZESzRUSzZOZ1p3Z1ZuaVlRYitxdXBHbU82WUtEckRRNQp2c0tNczFmSWM3TGFkYk5CbEMyUmxCYUZkd2lQOVhUVld2bnVDUHFZOC9MaUIyOTNDVmR4dkxsdE1razBPNjhWCmtGWG5oaG16dmZQZ3RlWUVPaSt5cDNiWWRqSk16N3dYSmxjQWEzRlVNT3NoREVaM2ZKanFxd3FUUUZLeXdiNHEKR3NxZThvRTQ1RGFsc25FMGFiSFAwaEY4SjJqTzkwcVZxMnpoNElqaldvbnJQQ3pzdVEvZ2R1RHE0RmprSTE1bApXRExhdGlKeG9vMVJiUTNId3VJOHBkSDg0YkpBQUE9PQo9NkRCbgotLS0tLUVORCBQR1AgTUVTU0FHRS0tLS0t"}`,
	`{"config":{"type":"openpgp","openpgp":{"cipher":"aes128","s2kcount":4096}},"ciphertext":"LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eDRFQndNSTdDNktiZnBINC9BZ1NrYW9FM3pWemtFOEtBcytBWkV1SzdUUzRBSGt6aWhHaGJZSEszY0NIRFdHCnd0R1NSZUhKU09CSjRDWGhMYy9ndU9KRmRwc3c0SUxvUkRiTnN4YUMrYzBRaFgvOGo2ejVmaE1LTWF1RTVrakwKQkg3QlNCRjB5aXJXbTVpUG9vbkhsSHFPWFlEMjgyRjF5Q2JIRWswOEJ2MmtRSG5QdWRPQThCM2VUSVFIamdndQpyUUxaSU01b1MzQUtZWTg0MzZxcnA2Q1lPK2ZkSHZOM29iSFNzdGhnRDZESWZIRUtvMkxSTDdJeTR5bkdURU1vCk1YbFhmTEJJWDBKZlFOYnBGVE1YMzRCSklKeDFGeVdTNlZwWEMyb2JHZHBJMmxoTlFURXVzc0RBTTl0VW9XdTkKQk12Y1BlUC8yVUVleTR2MkYvWEgyK1VaK2Zmc0ZMV05hdW9QT3VLUk85WU9ZcVpxN1M0N0JsWG8rV0tDdk1saQpjT2dHdU5kTTRkUkh1ZDdPd29LS1dJUEtZa3N4aXY5UEJ0UFpHQ1BPSkJqM1FFSS9Od29BZXVDVDV1MjlBbGE5CjZ0MnFPaFY1SnFrT3dOdXd1MGJ0K3A3Q0gveEYzaURDdWRTZCt4Yk9XOUVYRndicjltUXFYYkQ3VDgyZUhhMHMKNXR6MVFETVYrY2lXNUh6Zzd1QU00SWprZVEyQkdBU3ZFNExhZEUvQk5mRXdJK0tmb2xRdjRWK1dBQT09Cj1rV1RsCi0tLS0tRU5EIFBHUCBNRVNTQUdFLS0tLS0="}`,
	`{"config":{"type":"openpgp","openpgp":{"cipher":"aes256","s2kcount":4096}},"ciphertext":"LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eTRFQ1FNSVZySDBLTVVxNnA4Z3BrM1VKYmZmVndVdlhkRUo4S0xVRjBiSExDSTVvc2NqajZOZkJVYjNwTVlVCjB1QUI1UEJ2VkxjOXk5T3VOcFNNODN3amdXZmhXelhnSmVCMjRVeFA0UGppT28yTTFlQkE2QXBsMmRCSWg5dXIKTElWOGxIZEdEamtaWktydXdiZFN0dXg0RDZ3TGZDZUpTMGpPUitUaGxHQVA4am1paDEwYlYxdVpTWUJ0N1JSZgpsWjFEWkFWZlR0dmJ4ZkdHWWFGV0tMWGxYZnB4cDVpdEtLOTN1RnFnbUhhVkNRdmlWdXIrR01jRHNoYnczOWlqCmVkc2ZVcGNPVTRQL1FQbGU4aS9Sa05POGF5cm1MVWVpaTVuOWZ3cVNtQ1dYd1Q1TTBFZGhuS3NwR3ZlT3g0MmEKcUtLeE1vNWtHZWwyMkRRYnZGWTVyWGo5WFQvYWhMWndpd2hrV01aU1A4S1RmSmE0OHdFNm8zU2M2cTFPa1NVQgppREhzS0hSZkV4aFpNNXlUN2M0YVRUa3JKOTFqNjFBeXhEeThlMi9JR01hVDRSRWdMd0NmdlBKMGVaQ1pkdEplCnI3b2dmTWYwSkJiZ1dPWUpFUTdhM0FZdVhRZjFtT0VZQmVoUmJmSHEvMkhQak4zUVEyS3VnMjlXYW5vNjJnSEsKcjU3cGk3Q3dQSlVMcG9ULzdIdFRib2VxbXNOV1V3K3NHN24xNE43Z0dPRHE1REgyWEpYZ3N2ZlBOYmlGUDRtRgptNkxpUmVMOFh1RXdzd0E9Cj1lMVZyCi0tLS0tRU5EIFBHUCBNRVNTQUdFLS0tLS0="}`,
	`{"config":{"type":"openpgp","openpgp":{"cipher":"does-not-exist","s2kcount":4096}},"ciphertext":"LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eTRFQ1FNSXN1LzJQTElXK1BJZytSc2NiQ3Z3SS9Kc2VmZWVLV25aNHZ3eGgyRDczVDdBb0dna29qUDg2OXVHCjB1QUI1SG9DOXVsN1hlYm9MT2N3TWFCbHM0N2hkaVRndU9DMzRiWm40SG5pVnB6ZzllQXY2Q0FPSkhlSmFWS1EKazNIYURZYUY0dEtlZEVqdy93d0tDbnlFUUUyVlJURlFTYk4zWnRIK0FONXdzOWhCYmptVHVJQWlUeUUzdUJxaApzSHY3WFhJcEs2VnNrU2U3alRuSnBUMy83NDV2cHNYTXc0SjR5ckhYU0dOeDdaUy9rUTFuWWZwc2JQSndkeGdmCnZaOGNmcW54MzlsS1plOWswYVVoSS83a1FmanVhOW5xSklkVUxVY1k5WkkzV2RFRTY5N3BhMFd0TlhyaDRiRW8KMWpId0lTdzY3eStiTzZIZmtIaVZLdTlEQkNIZHovNnd3ZFVSRFc1UmlRUzRLNnpMd3A2UDJ0VEdXMjdIeE91YQowZ3Jnd00ySmxYNTljV1QvVXVTYlNnWTk0R2JNVjRLNjl1WjJDcEU0UkF1Z3VmWXNNcmJ6SEJTVi9QNDNJbDhSClBucUdnbkJmZmxiZ0x1YU8rY0M5SXU1SGR0ZzVUb1ZkUHd4OXVNaGkzL2E4TE1MTVR5RHd6OStQdktrZGZMMGQKYmF1ZGdnNGRFM3BObGVhQWdsUHUyTDBLYmpldjFsM0p6RGlKNEliam1mZ2RaZU43MUl6Zzl1Qis0QWJrbWp4cApXeEdsNnpWZWhqcDFQMi9mbCtKRDZCTU40YlY0QUE9PQo9SHAvUQotLS0tLUVORCBQR1AgTUVTU0FHRS0tLS0t"}`,
	`{"config":{"type":"openpgp","openpgp":{"cipher":"aes128"}},"ciphertext":"LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eTRFQ1FNSXN1LzJQTElXK1BJZytSc2NiQ3Z3SS9Kc2VmZWVLV25aNHZ3eGgyRDczVDdBb0dna29qUDg2OXVHCjB1QUI1SG9DOXVsN1hlYm9MT2N3TWFCbHM0N2hkaVRndU9DMzRiWm40SG5pVnB6ZzllQXY2Q0FPSkhlSmFWS1EKazNIYURZYUY0dEtlZEVqdy93d0tDbnlFUUUyVlJURlFTYk4zWnRIK0FONXdzOWhCYmptVHVJQWlUeUUzdUJxaApzSHY3WFhJcEs2VnNrU2U3alRuSnBUMy83NDV2cHNYTXc0SjR5ckhYU0dOeDdaUy9rUTFuWWZwc2JQSndkeGdmCnZaOGNmcW54MzlsS1plOWswYVVoSS83a1FmanVhOW5xSklkVUxVY1k5WkkzV2RFRTY5N3BhMFd0TlhyaDRiRW8KMWpId0lTdzY3eStiTzZIZmtIaVZLdTlEQkNIZHovNnd3ZFVSRFc1UmlRUzRLNnpMd3A2UDJ0VEdXMjdIeE91YQowZ3Jnd00ySmxYNTljV1QvVXVTYlNnWTk0R2JNVjRLNjl1WjJDcEU0UkF1Z3VmWXNNcmJ6SEJTVi9QNDNJbDhSClBucUdnbkJmZmxiZ0x1YU8rY0M5SXU1SGR0ZzVUb1ZkUHd4OXVNaGkzL2E4TE1MTVR5RHd6OStQdktrZGZMMGQKYmF1ZGdnNGRFM3BObGVhQWdsUHUyTDBLYmpldjFsM0p6RGlKNEliam1mZ2RaZU43MUl6Zzl1Qis0QWJrbWp4cApXeEdsNnpWZWhqcDFQMi9mbCtKRDZCTU40YlY0QUE9PQo9SHAvUQotLS0tLUVORCBQR1AgTUVTU0FHRS0tLS0t"}`,
	`{"config":{"type":"aes_gcm","aes_gcm":{"keylen":32,"keyderivation":"pbkdf2","pbkdf2":{"hash":"sha512","iterations":100000,"saltlen":16}}},"ciphertext":"eyJzYWx0IjoiRkJETXg5TSt6dXAvSWpkOWkrQkh3dz09Iiwibm9uY2UiOiI3ZUsxcmVOc2lPNFRscVZZIiwiY2lwaGVydGV4dCI6IkFNbTNxYm5Pbk4xM2xsY3V3UnM0cEpTU0dMVUVlQ2VCbUlETkVFRmFwUUFyNDdxUnJtNzArek5hQmJ6K1o1TFFVSVJSd2RYR1pNM1p3UHYzMTF0TW1JR0VOODA5NkVGOFBnc1JZLzhsTTl3dUNJN0EzS3JKazRzOVRjRStwbFZNZUNzQ2FtQU0rS29UVGVxREd1Umo0QkY2OFNPV3E3YVB5Q3hXTWhmR05OZFFhem5VQVYyUDM4ajZYZ1gzdjh1aEZ0VDBVWHl0dHNRUVRHVXRZakdEMUE0K0FKWmxNYzlxZ3NabWhqYzg2d2tSb2I3bHYvUVFseWRrbWRoWXU5SmdyOG1JV1lrUlRSVng0RFBOUE4wMkV0Q2ZoK1BZb0QxdzVMb3hKNUdNdmhhUE5CUEdvOHlycnBKeXpjK3lnTmFTSSt6ZndBbVliYlhwS0syNGdReEJPTFkvWmFaa1lETzhjTVFiN2xrV0dVQTU4RHhaY3BjdzMyTytRNEdGK08rdUVHc0lKcWN1SzlLeXNQU1ZHdlJoaGVhUFE0T0FWdUl1RG5GUjFnMGtGOXdtQ0NOd01QMkZlSjRLOVl4YmRDMUtlOEFLanE5bmJkaHFlQnBKZm0yM3p4OUl2aVdmSm00bUM5QzY0NDNJMkUxcHZnNFB5OVNNSDBKNVNsY1hBK09WQVB1cGJ3NUIyYXh3OTRFZW1CZHFMUT09In0="}`,
	`{"config":{"type":"aes_gcm","aes_gcm":{"keylen":16,"keyderivation":"pbkdf2","pbkdf2":{"hash":"sha256","iterations":1024,"saltlen":16},"pbkdf2hash":"sha256","pbkdf2iterations":1024,"pbkdf2saltlen":16}},"ciphertext":"eyJzYWx0IjoiakR4ZzFXWTFDbWJEeVVpeEgzS0RnZz09Iiwibm9uY2UiOiJpZEYwdER6czNjOWtuNmdvIiwiY2lwaGVydGV4dCI6IldabVpnbVhwTGx5RzJnTitVRFRVcFBzcWdyTWM3Ui8wajl3NW1zYjF5K1ZYNGZjOHlMVGgyV2t2UG5yVmxVNzcweDlHOTRqT0VRaFZGa1F3ZVMzVGpMa040Ukxxbm9YNEhCNWdPNzBtOHBENEV0MG1uUEVPSjJqTzJnTmNWcXBrK2lHaDJ1TVBBSW5KaHlDam8yK3BIT2FwVW1JeUR3ZS9FRCtlNUNPRXVGSVdsUERxaEZXd0dtNCtqUDErUmhUaE1qbzVNdTVPZWJybXF1ZE50NmNlVDd5R1dVNkFGdFhsYU9PNk1XL2FEdC9hREl5d2lOVXNuRVl3aDdqUy96MTJaNlVOSmNEei96MWQxb0Q0ZE10SEVuZGhyenpJTE9teC9rKytmWENaNE9GMExOVmpiNVZkWDFJUjgzNGUxM2JVbFpyWGFRcE01MmhwcG0yenhNbVBvRSs4UG54dlNCS2lzOWlxV2JRUldLMDhrZHFaNlJUZzFxYWtaWWJnclBwYmpYb1ByV2pvMVhpVXhkYTdsT3I1ZDFKV3hpQlRaNGRMZ1F5b0xndGtsSmpDWi9GNzBDMUxZdXRZYXNuMW5TTlN5L3BRajBDWTF4dTBUc2ZkS2JrMmhhbi9mS1crd05RMGFOTlZRQzd5MVVwMThYRXN4ZC9GWlZubTZuUUJMQ0EzUFYrdDViZ1MyVWdOeUd1N3ljdVRsQmpkbmdweENwZ0Jjc3FCKy8vUU94cTk3NnRnVFFWOTdnc1RjK0wveG5jeDhJd2U2bGErQVJTVThYN3ZxWUFiWi93cE1tTEhUajNGTDgxQUFEWGFNZz09In0="}`,
	`{"config":{"type":"aes_gcm","aes_gcm":{"keylen":32,"keyderivation":"pbkdf2","pbkdf2":{"hash":"sha512","iterations":4096,"saltlen":32},"pbkdf2hash":"sha512","pbkdf2iterations":4096,"pbkdf2saltlen":32}},"ciphertext":"eyJzYWx0IjoiR1BNT2E4YldCNTVOUHNwTDBQUWtlTTVNOVE3WVlYeGVHTTB1dTlnbUtiND0iLCJub25jZSI6ImlsbXUvRWYzZFE3QXBUOVciLCJjaXBoZXJ0ZXh0IjoiVFBkTHd1ZkNyWjlZVUNUOEZNQTVTK2N1WlVhZmEzQ25tMXVYV0JtOFA0WElwMG1JWjFkNCtPU2N0TUVJbkM3Y3lqcXc0b2JkVk9pMFpUZkRITXQwSm1DVzlxeGFIZGtkTzcrK2lQc2d1cDEyTi9KQ3V2aWNPeC90cHVtS2x5U1UrMXMyWTFQQ2ZyT1RUempwWHdON1hSTzdNMjVua2R4Y1BVcXVib2dEcTd5bXM3YzJnL250UG1tdkg5M0dHRGZYRDZQOFVseHdaOVJJZ0RVemtVUDFXMzc0Mmg4WnJMU1NzRGV0Mm9vblQ0Nmdadit3NlQvUEtHUVQ1TXVER09maExiN3Q0Z01ibmViSEg1STR5SXQ2cHRKYnFnU0hyVytYRWtLQ2tYVkxUditGRjlkK1Q2Z3ZwR3UxMDhDNWRiQklnU1IyZTR4a1FDZGd6MTlyc3ZPTHh5VStWMFRhN1pmZ3lIMllOZFZhNllYTW56NHZuZVBvZ2pEVFhFdGd1WXRtSDZDZm96UTVpdUx1SHE5SDhxbVZNQk8xWlExQ00xTHFnZDQ1TVJqUzFpQUZJa3BIUFhtS3ROY3lKZnNzUGd4WWh0eTFyMzZMR1NmVU9iREJ4cG9lcjhKdkw3M1hNc2psRUVvN01mTEhLSWtUYXh6MzZ6cU1yNS8wdHU3SXhXSUFEckltOWpKKzd3cGR3WnFWVUFPUVN4RzZISWtDVmxHL0tTbFgyMytNdm1uVHIyWnlHSmFGc1hpVkFpeDBsSzRSY1dYMC9NZ3ZVSWJlMGI4NGRaQUlsYmoydXdtVDM5ekNnNmpHUGtyU3F3PT0ifQ=="}`,
	`{"config":{"type":"aes_gcm","aes_gcm":{"keylen":16,"keyderivation":"pbkdf2","pbkdf2":{"hash":"sha256","iterations":1024,"saltlen":16}}},"ciphertext":"eyJzYWx0IjoidUxnQmI5eVAvT1NHQVMxekJGSDZOZz09Iiwibm9uY2UiOiJleFpTS1NmbjJpTm4zWWlVIiwiY2lwaGVydGV4dCI6IlJDZ3JIMGw3ajdnQ254eUhoWTVYaVJMOHRrOW94dmVuRkcrNWx5Z01ORDgzZDJ2R3pxb2x5bFhsbGxwM1NOSWlaT1g1UHY4dnRhdkVZMnpYblJ4eUdVWVZiRTFTb0RrNDBmUi9ZendFckh6UHI5MkdjSjNKRHB2T2lBZGFpbzdsVTh1MHUvZjlNQkYzL1JzOWlQZFNLTVZKeTQxSlpUL1BQOTBoMklwMXpMUjNzTkN4OStxT2J5TTJlSXd1RXF3UDcvVTZpanRoeXFPWFBIdmNHTEkyQTM2ejdUamw1S1Fsd1JIeWpHQVJ5Y0paK3FxWHlKRkdHdHV4TkJ1cU1wMFdBWm4xUUQvQ3VHZDRLby9FMFdqbUZHRUFaRFZPNWFmdjZxeVBGazVLeE1BbG0xM05zQS9GVzdyZ0dlSGtkVlV2NnNadGhIcUhvVVhmbmM0OW1YZHNvbUpuWm5BTWJHYkdQc0RGT2tYcU1uTVNaS3BBVzN0c29WTEdveU9HWUw0M3VEKytJWXJ6ZzhJSXRnUmdpVjQ3M1hTR1pzN0IyRmJjNnhKSzJucHV4b2hObjV3THdNZU4xVDdBTXV0NTdiR2ZISXF0bW1kSFh2b1k2WlR2MDJtVWZ3dEM4dDhXWDdYcHBtSVE5Qjc1c251YzZ6Mmo5RUhqdnd1OHlkU3YxaGY1RmdzcHhrN1JJQXRuOUM1S1Y3az0ifQ=="}`,
	`{"config":{"type":"aes_gcm","aes_gcm":{"keylen":32,"keyderivation":"pbkdf2","pbkdf2":{"hash":"sha512","iterations":4096,"saltlen":32}}},"ciphertext":"eyJzYWx0IjoiMFIwaE5JOGxTSERrMHVOZWQwd01zNEtqL2JPUmZKVlp0WE9IQUVrVnQzTT0iLCJub25jZSI6IjRzSG5SeXhRNjFkdXBtN2kiLCJjaXBoZXJ0ZXh0IjoiRDZhbEEyM2MzcVdybDdsU3RyUGlKdmc4QkdyMDZUZ0o5cjBxUnBBZmJ3TTFsbFoyZDliYnl4SGppSXY1TExPWXUrNGppMnBZTXp1Y1NYRUVxU1NEcHlXTlp4YThzc3NuVXhCR25PaElPR3FlN0N6RmlFbG1HSnpRWVlmMXROODE0akRhMi9EU1phalF1WmNKdGdXb1B1bEJmaFozUkFqQjdOYVRGMWFwb3VLSkdEYTFvTmwveFRDNE1CWnFiZFNtSWsvQldQT242ZmU5enhEeXl3OWd0L0diakVLeThhSnQvSEErdUlRTHp2WityTG83SjltQ0ZpYWpTRTFUOFl0c1F2cnppOHZjK2hXSDRmWmpKV29TRnlRMG9xR1JnS3lMTG5wYUpuenM1YllMSWlHbkw1azh5cEpocHZOdVNoLzVNRUZ0Yzlob0dzWlhSaDBBSDhOYitsbjZGL0lVS3pwTWpOQlo3Rld0RTZUM1B5VWVUZUFhOXpaS3pZRG5rR0dCbkxLOVlWU0ZrOG9aZ0IreVNMUWlHVGFNTEtncC92ZnRKZ1YrYkwzM3E3RWZwc0ZYbkQvMmdRd1F6QVhzZWVSY2x1Skt0Z2JROFZJWUMyR0JrSlNPL29NV25VbUxGNktoUk5VRE10NmNvakk3UDYrQTlmUzZmUzVPSmpYWlBJRlpEa1l6U3REc2IyOFNCV0VvbVlrPSJ9"}`,
	`{"config":{"type":"aes_gcm","aes_gcm":{"keylen":16,"keyderivation":"scrypt","scrypt":{"saltlen":16,"n":4096,"r":1,"p":1}}},"ciphertext":"eyJzYWx0IjoiRUh3d3BuVXY3VnlwYStGcGtuQ3d5UT09Iiwibm9uY2UiOiI1SGphZTA2ZTdtcW8zUGZPIiwiY2lwaGVydGV4dCI6Ii9sVElHUHdacWpXeTY0VGd6TzBOSVdJNEZFVHVGVVRhQmhWQVBvTlkyaFhMVVh2RW52VVRCU0JXRFJnVW11RGFVNnBYZ25uSVZmZHE3V0xvR3ZDUkthRHlBUFM1YlpvQUFWNGYzOXFGR1E0dVdzRzdwc0tvUGhzTmpHamNWbnYyUE9YZmRJUi8yV2FXVDA3eGlXM0IzSVJBcEZmNkloaTVzRUMyeUErY3hxWTVMdTMrajZFY3M2SjZRMDdvY2ZPVGd6dW92c0kyYVg4a0ZPRHExaVIxSWVWNUdWdTB0ZVNvdHBXSWpCQTNsV29XZko4RS9xT1luR2Zibm90UFI4K3pxOUgyM1Q4QW84UE9HVmRYT2VoZmo2aWpnOUltd01qeFF0aFZrUFAxa3dvcDg0QXBCVSt4akM5MzNYenhFOUgzcTBHbHRTVnZzdlhrcVlMMy94TnptTzFncmg5QmxnbXlETzg0S1kyUHFxL0lNTzBaYjZKdkpTTk9HcjV0dGdIN2Z5QmtQeWs2ajFBc1Z1Qy9XNllFekxrM1VRb2lYVVhiQmttV2JvQXpEZjlVVjlBdmswZ1pwMXRHTzZ2dzAyMGs1Ti83R2dIQmhVc09HSEh0OHVoMklhZlpheGF1UldHa1h5Wm05YkJYd0xYMEdsNVhub0c5OHlXWTFZM3hmU3BqYmc9PSJ9"}`,
	`{"config":{"type":"aes_gcm","aes_gcm":{"keylen":32,"keyderivation":"scrypt","scrypt":{"saltlen":32,"n":8192,"r":2,"p":2}}},"ciphertext":"eyJzYWx0IjoiOEFhR0lZN3EyRjBlaUQ5bjE0TVhER2NCUlB4ekNNTmwyalJVZUtJYUxXTT0iLCJub25jZSI6Im8wVVhhV3RqKzArVzZIUlMiLCJjaXBoZXJ0ZXh0IjoiTU8rWG1sL0xSMWMvV0NKZjlqcFgyVHQzSStzamg0TStkayszSW9TQ011RlBsSHl6TjZXbFJnN2pDeFZQWStvSFV1a1ZGMFlXOXhzWTUxVnpIYStoTFhINGhCdmN5UTNxRW9RR3pCenZrd29hZGN5NkFRSkxOSXJCcWgwcTYzQWZSOWJteW9Fdk8xOWcrSW91OU5xQnpBQmxDRm1GVUFVWmJhN1JOeGY2WUZoVm1WTW54aFlIWjJIQllrYUxaT3FDN3VURGVoRGJPY09XdFNoT25aakxTZnNFTXZrSURhS3NqeEVlNXdUby9MMDJsMlJ3SW0xRENTRm9odnZxYUxmcjFOTVpMeHpON0QzcVJPUjl3UUk1VjNRVVg0VWc2cmFlVUJVUVNQQndIVzc5b3c4cHpLa1lFOExrSSt1dWFUdWwrQk9yVFk5U2kyYXRmUExsRGJTdTZ4WEQwd1FkL0p0L1dWaG1FMlkzOFZoTlR0TTV4ekt1VEUwc3JYaUxJVFdScWVXbytnbmZvTXZZa2hiTlViUkZzeW55Uk9XRTE1dlZWNmtLcEMwNDNVMWt6V1FxaWJOYUYySXlqeWE4SURUMy9mVUppcVBGOWJQYWF1VE1PUlFqQmFVbi9KV1h5Mm1lV1dab0FocnBqdnJJSVF5RFdKU2E5UTFNZDFadlBxUUcwQT09In0="}`,
	`{"config":{"type":"chachapoly","chachapoly":{"keyderivation":"pbkdf2","pbkdf2":{"hash":"sha512","iterations":100000,"saltlen":16}}},"ciphertext":"eyJzYWx0IjoickpGR2VnQVdNS252UjRJWGR4bGNNQT09Iiwibm9uY2UiOiJBSkRKVDlOV2Q1SUpDVmtEIiwiY2lwaGVydGV4dCI6InpYMzZVV3J3T1BVT3pMcGlvdldReDZQcXlxTytZTFBSVERORmk0eXNlWUhPdFo5V001VGVVcXdFL0kzRUJKenVJdWEveFpxWDdtcmNCYXZMT3NYck80NDdzRGRvL05IMjhmdWhTd0hKRU5GMVVIaGU0NGdab2lsVi9mTER3bi9ncjVObkhzOXFrSXJ0Yk1ERTFTUVZIWVN2M01Oem5qOUF5OU5UZldvYkIxcUtqNitFc281QkVhMFUwUEZWYVE0Ni9jdVRhVFBncTh2Vk01U1hJZk5pblh5YkpzYjhJdkdJa2NPVFJrN0tPRCt5dVA5c29EcEN3M2k2dVl5Mk1FM2tWL3l3Y0NPeWlMZlI2MWdhNzNQZ3l4a3RhdlFwcTd1ZFRFVldvTWhad0paazNxMlBuWkZmaVNwUlV2V2luNG1rL2w0Q05vSlArekx5WEFiQkVISElKTzNvWjluZ1BmWE9IRlViN0x4UVJVR0RrcDI2Zmh5bzJhVTcvQlR5eXRWSDFkUCtNNXFXSXU4RktFdmNJbUtjOEJuVjBqcWUxd1N1cnlVZVlhZjZlZ0xnbXlxODhpZFQxQk0wK3k2ZkVLdkhreU9vTXpINUlVdEt4bnR5b3FuTGg2b2g5K29YamVNc3FPMmNRTnpWS2s2SFpmeWpmOVQvb3Z0N3cxenlScUUzcjNzYlNsOUdzdEpCZlE9PSJ9"}`,
	`{"config":{"type":"chachapoly","chachapoly":{"keyderivation":"pbkdf2","pbkdf2":{"hash":"sha256","iterations":1024,"saltlen":16}}},"ciphertext":"eyJzYWx0Ijoibm9KYzQvUFhMNWhpQ2tzdzgrMVcyZz09Iiwibm9uY2UiOiI1TFVvb0tYMGh0amdtSXJ5IiwiY2lwaGVydGV4dCI6IndGamVwWWErR01sRWc2TEFDenRrVlhaWVdNSmRjSkxmTmJlZzRNRG9QVEowSURlem1LZS9HMTZBSXZpMXl2dXJuMWVHMmZreHVHUS9LbDJqNGI1d21zZlVJY0NRa1pYOUszV3MyaC96U0xzUWNKM2tGcmR3MjJ3aUNsMHdycTNjdG9OeG5wd0k2UGppVzBmWmwrV0lRUk56d3BiLzBkVndzSHFaSnYyQWFGRlU4eFJrR1hDK1QyU0lJZ3JtTVV0ck9kajRsamtXODBGR0Y1WFh2YnhWUzlvRzlTVzNWTTIva1dpNTAvT2prSmlYbFd5ditTWU50eFBIKzdLbTUvMlNNTUlBMy9BeldwRGNia3QwRkFmNmozMk16Mm01Zk1SOHFzcGgyQ01RZUVGaGRzRVRBcXMyK2JKeTN1WVJ6RmpjVDFBL21tOUowdlQ3UDk4TnBPL2NOd3hiY3hVcWRXaytnSXc4eE5kZUFJbHRibnNGN0x4NlV5a0FYNmFRLzFBKzhTajUvRG1ac0NYMTBBWjV1Z3liUWJrMkZGbDNNM1J5Nm5wS0lSckJsSUVnYWVaNC9sYjlKK29peWcvc2t4eTNSNHlkRkNsMEt1SVM0dDdjWlQ1TGhkaFowL3IrR28wcUtrcXZ3R25kcGhzcXRpaWxDTlVvaEt1T243eWhTNnI3OGFiVGpCMzFFTGM9In0="}`,
	`{"config":{"type":"chachapoly","chachapoly":{"keyderivation":"pbkdf2","pbkdf2":{"hash":"sha512","iterations":4096,"saltlen":32}}},"ciphertext":"eyJzYWx0IjoiaEpiZDVzejNOeVRwNWMzYm5tMHN3ZWxzN0NNa2tiWmtOejFPMFJ0WSt1RT0iLCJub25jZSI6IkVMaXI4QUp6THdPeVo0WksiLCJjaXBoZXJ0ZXh0IjoiUGQ4YytybXVNRGhjQ0FVTW8rVW9pVjJYSWJHSXorUGF3UEJHZnVRRTJLTGcvOGp4U09vaUxLSW8vcmJIcDRVRWExYWI0YTRNL2pwYWJJTEZNck13UU1lQUIxMFJCS3NuSWtabWt2U3BnU2s0UVhNLzgzWHBTZTR4VDU0STViVVFsWnJqeUx5OHJ3RHBGdERLekZOODNHSzlxRS9jek5BckRmSm1mUExvRUJRQjQrb0xUR2U3NU96UXpiRVdLa0l6RHo5Unk4NDNFQ3k4eHlmbTZlMVlDd0krR2twRTZmNzdqeVpCWUZnYW9JNHZxUkhhcUJWR05DRHdCK2JUYnNVTzhGRkEvbGo3N3A4MDJCMzcvUEF4YmdZRUd4QzJHVWFwUVE4U1V0WlN1QzZ3L1N5ZVRiVzhzMjJTWi95MSsyeE9Vd3duTmd3K2lTckowUXUyRFdPSmZkaTVxdkNTRWttM1NOWGVRYnNjWk1hd0RLNFpOZTRSSm55aTIwaEJ1T3YvVTFQdFdKcko5M3I0TytCUEYvSW91WkQ3bDlLWmVmMjA4aEpZTmY5YnI2b3B2WDRGWmhvSnM1WGNZVml4LzB5ZWdDSUtvbW5BV3NybFNhWDFKdDM2TlE5cHhrL0ZYTkUxQ0lqTk45NjJ3eGZKeXJBSEl2SUdpSzNpN3dHLzU5OHM2VDg5K1o5VFRCYz0ifQ=="}`,
	`{"config":{"type":"chachapoly","chachapoly":{"keyderivation":"scrypt","scrypt":{"saltlen":16,"n":4096,"r":1,"p":1}}},"ciphertext":"eyJzYWx0IjoibTRyR2lZdnJ4NGF3Wk83Yzl2K09Ndz09Iiwibm9uY2UiOiJKeitwRTBWR0o0UlUvdEtQIiwiY2lwaGVydGV4dCI6IkxOcnZ4OTE5QklYMGtTbm83Nk9YOFBHWEVBM2ZJSWN4SHNWYS91TVh6MmRnZDZTTFVoTS83aThRN3d5aG9naStTZ3pKMEZzRXRES3ROWXkzTTZqdmhrczRKNDFWbkNWUlI4VUx3ci9URHJYQVFTaVlkYUdPQ2V2RHhnL3hCblBZaEhRa1RKdzF6amU2dFdEd2hENVBvUHBTNVI0UmU1Q2JxNjJaOERyT1RBVnhNMkRyVlNiN1FWdmhuTVRLblJrK1duYysybHZsVjNHak1pNnhZL2JES00ySWZwNWQzTlJuWDFRMXUxZVNqN3FUZFlEVnhEMDh3ZnZCNmlMZjRka1VlQzlGeVlycFpZYTlkZVFORjVqUnhnNjBXOE1OQUtzUWx3OEYybVpoaFJ2S1dWcldvTTczc1FOc0VlNnRkaWYva21jbVRIa1NaT1dhYkJuQkN4eHBpTk9TeXRvMGNaUzFSYlFwb2wxNG1QWFNtbGFpR2VyY2xrT2JqV2RydWkzUHMzVXRFNUl5STc2WnhvQUVhTWdzVTJ6VTMxcms1blVPdVVaemEySDQ3ZVUrdDFNbGxDajczSGgvQ1UyZGREK2M2L2VXeFBFT0RFSm5Rd2llOUJtWTlmTEhFUzJ6M2dpQXh3SmpEZkpnRU9kNHRSRlRURWZjWWU4V2J3PT0ifQ=="}`,
	`{"config":{"type":"chachapoly","chachapoly":{"keyderivation":"scrypt","scrypt":{"saltlen":32,"n":8192,"r":2,"p":2}}},"ciphertext":"eyJzYWx0IjoiVlFuRDZKWFdXRGsxVHZtc1RmY3FSK1BENWZKaDUreUtwWUZNSmNSb2JNcz0iLCJub25jZSI6IkhHZE9nOWsyazFBdGg2dEUiLCJjaXBoZXJ0ZXh0IjoibmpIaWJIaG05L3B0U2h6YUNMS2EzclVxN3VpQ3NKQ2xNSEpDVVQvQnJrV2tRUCtCdzY1SGJTWGhJQzZhMzBWUTJQN0VnTlViY21kR3FTVFEveFNnbFZMdldQVGtDaHJxUmdRcE5qQTN3UVVldEI1QkZtdzJBZkRUS1RFaURXUjlIYkxEdDlpdHI2aDgrclI5cmUybHhlcnlaUTd5ZkQ2TXNSbmJybVBKcnpmaTZNWmZPanNJVVFCVks0b1JJTkswT0dxUVJGWjdCbmxBN1JPNm9yYlRTOUNQT01xWVpuUkZaRm03MHVpUFBqdWMrNCtKTUNJSW5GR3AvRzRJVm1iNHQ0cllXcG4vZ011WnpuNGhiUEhDT1RqRURlYWt5Uk1jWVdPL1Y5enpRaVZTRFpCZmpQcXluOUt2ZStyTGliSEViL2JnWXM4RHZsWXErWGVNMSt5RU9IejlkQVZscHI3b01pMndSUmFvKytOS01ZZ0dpQ2g4ZEpCRy8rZEdYdzdNZ09OSjJDd0ZiTm1NL3hxOEU5ZWZGa3hzR0hvVGFNYXJhZ1poQ3BHUjl0ZEw3RTczRGhrWmFtbjM2czBPazhZbWRJclFGQjNJY3ZIdHFjNXFnb2hWUFBWSUhoNlphMlRYQlV4RWtiSFJXc2kvbTJ4VksxVHdZejdLWlE9PSJ9"}`,
}
