package backends

type Config struct {
	Type     string                 `toml:"type"`
	Backup   backupConfig           `toml:"backup"`
	Settings map[string]interface{} `toml:"settings"`
}

type backupConfig struct {
	DirPath     string
	AutoEnabled bool `toml:"auto"`
	MaxFiles    int  `toml:"max"`
}

const (
	ConfigTypeFile = "file"
	ConfigTypeS3   = "s3"
	ConfigTypeMock = "mock"
)

func NewDefaultConfig() Config {
	return Config{
		Type: ConfigTypeFile,
		Backup: backupConfig{
			AutoEnabled: true,
			MaxFiles:    100,
		},
	}
}
