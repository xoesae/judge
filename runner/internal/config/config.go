package config

type Config struct {
	Port      string
	UploadDir string
	RootFs    string
}

var instance *Config

func LoadConfig() {
	instance = &Config{
		Port:      "8080",
		UploadDir: "/tmp/scripts",
		RootFs:    "/rootfs",
	}
}

func GetConfig() *Config {
	return instance
}
