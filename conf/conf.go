package conf

import "fmt"

var (
	config *Config
)

func C() *Config {
	return config
}

type Config struct {
	UDP *UDP `toml:"udp"`
	TCP *TCP `toml:"tcp"`
	Log *Log `toml:"log"`
}

type UDP struct {
	Host string `toml:"host" env:"APP_HOST"`
	Port string `toml:"port" env:"APP_PORT"`
}

func NewDefaultUDP() *UDP {
	return &UDP{
		Port: "",
	}
}

type TCP struct {
	Host string `toml:"host" env:"APP_HOST"`
	Port string `toml:"port" env:"APP_PORT"`
}

func NewDefaultTCP() *TCP {
	return &TCP{
		Port: "",
	}
}

func (t *TCP) HttpAddr() string {
	return fmt.Sprintf("%s:%s", t.Host, t.Port)
}

// Log todo
type Log struct {
	Level   string    `toml:"level" env:"LOG_LEVEL"`
	PathDir string    `toml:"path_dir" env:"LOG_PATH_DIR"`
	Format  LogFormat `toml:"format" env:"LOG_FORMAT"`
	To      LogTo     `toml:"to" env:"LOG_TO"`
}

func NewDefaultLog() *Log {
	return &Log{
		// debug, info, error, warn
		Level:  "info",
		Format: TextFormat,
		To:     ToStdout,
	}
}

func NewDefaultConfig() *Config {
	return &Config{
		UDP: NewDefaultUDP(),
		TCP: NewDefaultTCP(),
		Log: NewDefaultLog(),
	}
}
