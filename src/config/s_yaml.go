package config

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
