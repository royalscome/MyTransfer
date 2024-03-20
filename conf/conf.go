package conf

var (
	Config *PortConfig
)

type PortConfig struct {
	Port string `toml:"port"`
}

func NewDefaultPortConfig() *PortConfig {
	return &PortConfig{
		Port: "",
	}
}
