package conf

type PortConfig struct {
	Port string `toml:"port"`
}

func NewDefaultPortConfig() *PortConfig {
	return &PortConfig{}
}
