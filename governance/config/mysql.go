package config

import "go-micro.dev/v4/config"

type Mysql struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Pwd      string `json:"pwd"`
	Database string `json:"database"`
	Port     int    `json:"port"`
}

func GetMysqlFromConfig(config config.Config, path ...string) *Mysql {
	mysqlConfig := &Mysql{}

	config.Get(path...).Scan(mysqlConfig)
	return mysqlConfig
}
