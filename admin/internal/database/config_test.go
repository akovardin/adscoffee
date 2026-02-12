package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Connection(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   string
	}{
		{
			name: "Standard configuration",
			config: Config{
				User:     "testuser",
				Password: "testpass",
				Dbname:   "testdb",
				Host:     "localhost",
				Port:     "5432",
			},
			want: "user=testuser password=testpass dbname=testdb sslmode=disable host=localhost port=5432",
		},
		{
			name: "Empty values",
			config: Config{
				User:     "",
				Password: "",
				Dbname:   "",
				Host:     "",
				Port:     "",
			},
			want: "user= password= dbname= sslmode=disable host= port=",
		},
		{
			name: "Special characters in password",
			config: Config{
				User:     "testuser",
				Password: "p@ssw0rd!",
				Dbname:   "testdb",
				Host:     "localhost",
				Port:     "5432",
			},
			want: "user=testuser password=p@ssw0rd! dbname=testdb sslmode=disable host=localhost port=5432",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.Connection()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfig_URL(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   string
	}{
		{
			name: "Standard configuration",
			config: Config{
				User:     "testuser",
				Password: "testpass",
				Dbname:   "testdb",
				Host:     "localhost",
				Port:     "5432",
			},
			want: "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable",
		},
		{
			name: "Empty values",
			config: Config{
				User:     "",
				Password: "",
				Dbname:   "",
				Host:     "",
				Port:     "",
			},
			want: "postgres://:@:/?sslmode=disable",
		},
		{
			name: "Special characters in password",
			config: Config{
				User:     "testuser",
				Password: "p@ssw0rd!",
				Dbname:   "testdb",
				Host:     "localhost",
				Port:     "5432",
			},
			want: "postgres://testuser:p@ssw0rd!@localhost:5432/testdb?sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.URL()
			assert.Equal(t, tt.want, got)
		})
	}
}
