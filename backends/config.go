package backends

type Config struct {
	Type     string                 `toml:"type"`
	Backup   BackupConfig           `toml:"backup"`
	Settings map[string]interface{} `toml:"settings"`
}

type BackupConfig struct {
	DirPath     string
	AutoEnabled bool `toml:"auto"`
	MaxFiles    int  `toml:"max"`
}

func NewDefaultConfig() Config {
	return Config{
		Type: defaultClient().name,
		Backup: BackupConfig{
			AutoEnabled: true,
			MaxFiles:    100,
		},
		Settings: make(map[string]interface{}),
	}
}
