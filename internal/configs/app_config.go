package config

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	BaseUrl string `env:"BASE_URL,default=http://localhost:8080"`
	Host    string `env:"HOST,default=0.0.0.0"`
	Port     int    `env:"PORT,default=8080"`
	DB DBConfig
}

type DBConfig struct {
	DSN      string `env:"mongodb_dsn,default=mongodb://localhost:27017"`
	Database string `env:"mongodb_database,default=shrt"`
}

func (c Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func Get() Config {
	var once sync.Once
	var cfg Config
	once.Do(func() {
		lookuper := UpcaseLookuper(envconfig.OsLookuper())

		if err := envconfig.ProcessWith(context.Background(), &cfg, lookuper); err != nil {
			log.Fatal(err)
		}
	})

	return cfg
}

type upcaseLookuper struct {
	Next envconfig.Lookuper
}

func (l *upcaseLookuper) Lookup(key string) (string, bool) {
	return l.Next.Lookup(strings.ToUpper(key))
}

func UpcaseLookuper(next envconfig.Lookuper) *upcaseLookuper {
	return &upcaseLookuper{
		Next: next,
	}
}