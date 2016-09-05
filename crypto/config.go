package crypto

type Config struct {
	Type     string                 `toml:"type"`
	Settings map[string]interface{} `toml:"settings"`
}
