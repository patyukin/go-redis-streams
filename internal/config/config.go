package config

const (
	ConfigFilePath = "CONFIG_FILE_PATH"
)

type Config struct {
	Redis struct {
		DNS      string
		Password string
	}
	Logger struct {
		Path string
	}
	MySQL struct {
		DSN string
	}
}
