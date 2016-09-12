package backends

type Config struct {
	Type              string                 `toml:"type"`
	AutoBackupEnabled bool                   `toml:"auto_backup"`
	Settings          map[string]interface{} `toml:"settings"`
}

const (
	ConfigTypeFile = "file"
	ConfigTypeMock = "mock"
)

func NewDefaultConfig() Config {
	return Config{
		Type:              ConfigTypeFile,
		AutoBackupEnabled: true,
	}
}
