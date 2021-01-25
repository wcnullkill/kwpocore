package conf

import (
	"encoding/json"
)

// MyConfig 配置
type MyConfig struct {
	DataBases []DBConfig `json:"databases"`
}

// DBConfig 数据库详细配置
type DBConfig struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	Pwd    string `json:"pwd"`
	DBName string `json:"dbname"`
	Driver string `json:"driver"`
	Params string `json:"params"` //合并合法参数
}

// ReadConfig 从file []byte 读取配置文件
func ReadConfig(file []byte) (*MyConfig, error) {

	var cfg MyConfig
	err := json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
