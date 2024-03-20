package conf

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

func LoadConfigFromToml(filePath string) error {
	Config = NewDefaultPortConfig()
	// 读取toml格式文件
	_, err := toml.DecodeFile(filePath, Config)
	if err != nil {
		return fmt.Errorf("load config from file error, path: %s, %s", filePath, err)
	}

	return nil
}
