package config

import (
	"github.com/spf13/viper"
)

const (
	DiverNamePostgreSQL = "postgres"
	DiverNameMySQL      = "mysql"
	DiverNameClickHouse = "clickhouse"
)

var (
	Cfg   *Config
	Viper *viper.Viper
)
