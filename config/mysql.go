package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"time"
)

type Mysql struct {
	User            string        `env:"DB_USER" envDefault:"root"`
	PassWord        string        `env:"DB_PASSWORD" envDefault:""`
	Host            string        `env:"DB_HOST" envDefault:"127.0.0.1"`
	Port            int           `env:"DB_PORT" envDefault:"13502"`
	Database        string        `env:"DB_DATABASE" envDefault:"game_server_example"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS" envDefault:"25"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS" envDefault:"25"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME" envDefault:"5m"`
}

func NewMysqlConfig() (*Mysql, error) {
	m := &Mysql{}
	if err := env.Parse(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m Mysql) DNS() string {
	if m.PassWord == "" {
		return fmt.Sprintf(
			"%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
			m.User, m.Host, m.Port, m.Database,
		)
	}
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		m.User, m.PassWord, m.Host, m.Port, m.Database,
	)
}
