package database

import "fmt"

type Config struct {
	Debug    bool   `yaml:"debug"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

func (c Config) Connection() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		c.User,
		c.Password,
		c.Dbname,
		c.Host,
		c.Port,
	)
}

func (c Config) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Dbname,
	)
}
