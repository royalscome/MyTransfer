package conf

import "fmt"

var (
	config   *Config
	FilePath string
)

func C() *Config {
	return config
}

type Config struct {
	UDP       *UDP       `toml:"udp"`
	TCP       *TCP       `toml:"tcp"`
	Log       *Log       `toml:"log"`
	HTTP      *HTTP      `toml:"http"`
	WEBSOCKET *WEBSOCKET `toml:"websocket"`
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

func (t *TCP) TcpAddr() string {
	return fmt.Sprintf("%s:%s", t.Host, t.Port)

}

func NewDefaultTCP() *TCP {
	return &TCP{
		Port: "",
	}
}

type HTTP struct {
	Host string `toml:"host" env:"APP_HOST"`
	Port string `toml:"port" env:"APP_PORT"`
}

func (t *HTTP) HttpAddr() string {
	return fmt.Sprintf("%s:%s", t.Host, t.Port)
}

func NewDefaultHTTP() *HTTP {
	return &HTTP{
		Port: "",
	}
}

type WEBSOCKET struct {
	Host string `toml:"host" env:"APP_HOST"`
	Port string `toml:"port" env:"APP_PORT"`
}

func (w *WEBSOCKET) SocketAddr() string {
	return fmt.Sprintf("%s:%s", w.Host, w.Port)
}

func NewDefaultWEBSOCKET() *WEBSOCKET {
	return &WEBSOCKET{
		Port: "",
	}
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
		UDP:       NewDefaultUDP(),
		TCP:       NewDefaultTCP(),
		Log:       NewDefaultLog(),
		HTTP:      NewDefaultHTTP(),
		WEBSOCKET: NewDefaultWEBSOCKET(),
	}
}
