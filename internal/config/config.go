package config

const (
	ConfigFilePath = "CONFIG_FILE_PATH"
)

type Config struct {
	Redis struct {
		DNS string
	}
}
