package conf

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

func LoadConfigFromToml(filePath string) error {
	config = NewDefaultConfig()
	// 读取toml格式文件
	_, err := toml.DecodeFile(filePath, config)
	if err != nil {
		return fmt.Errorf("load config from file error, path: %s, %s", filePath, err)
	}

	return nil
}
