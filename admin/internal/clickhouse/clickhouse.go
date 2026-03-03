package clickhouse

import (
	"fmt"
	"time"

	"github.com/uptrace/go-clickhouse/ch"
)

type Clickhouse struct {
	DB *ch.DB
}

func New(config Config) (*Clickhouse, error) {
	db := ch.Connect(
		ch.WithDSN(fmt.Sprintf("clickhouse://%s:%s/%s?sslmode=disable", config.Host, config.Port, config.Database)),
		// ch.WithInsecure(true),
		// ch.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		ch.WithUser(config.User),
		ch.WithPassword(config.Password),
		ch.WithTimeout(30*time.Second),
		ch.WithDialTimeout(30*time.Second),
		// ch.WithReadTimeout(5*time.Second),
		// ch.WithWriteTimeout(5*time.Second),
		ch.WithQuerySettings(map[string]interface{}{
			"prefer_column_name_to_alias": 1,
		}),
	)

	return &Clickhouse{
		DB: db,
	}, nil
}

func (c *Clickhouse) Close() error {
	err := c.DB.Close()
	if err != nil {
		return err
	}
	return nil
}
