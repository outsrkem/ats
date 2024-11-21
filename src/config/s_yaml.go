package config

import (
	"ats/src/pkg/crypto"
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type Config struct {
	Ats Ats `yaml:"ats"`
}
type Ats struct {
	App      App      `yaml:"app"`
	Database Database `yaml:"database"`
	Log      Log      `yaml:"log"`
	Uias     Uias     `yaml:"uias"`
}

type App struct {
	Bind string `yaml:"bind"`
}

type Database struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Name   string `yaml:"name"`
	User   string `yaml:"user"`
	Passwd string `yaml:"passwd"`
}

type Uias struct {
	Endpoint string `yaml:"endpoint"`
}

type Log struct {
	Level string `yaml:"level"`
}

func (c *Config) decryptionDatabaseMysqlPwd() {
	if c.Ats.Database.Passwd != "" {
		if plain, err := crypto.Decryption(c.Ats.Database.Passwd); err != nil {
			hlog.Fatal("Decryption of database password failed. uias.yaml:uias.database.passwd ", c.Ats.Database.Passwd)
			os.Exit(100)
		} else {
			c.Ats.Database.Passwd = plain
		}
	}
}
