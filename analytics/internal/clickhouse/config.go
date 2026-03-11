package clickhouse

import (
	"time"
)

type Config struct {
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	Database        string        `yaml:"database"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
}
