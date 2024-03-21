package conf

var (
	config *Config
)

func C() *Config {
	return config
}

type Config struct {
	UDP *UDP `toml:"udp"`
	TCP *TCP `toml:"tcp"`
}

type UDP struct {
	Port string `toml:"port"`
}

func NewDefaultUDP() *UDP {
	return &UDP{
		Port: "",
	}
}

type TCP struct {
	Port string `toml:"port"`
}

func NewDefaultTCP() *TCP {
	return &TCP{
		Port: "",
	}
}

func NewDefaultConfig() *Config {
	return &Config{
		UDP: NewDefaultUDP(),
		TCP: NewDefaultTCP(),
	}
}
