package config

type Config struct {
	Ats Ats `yaml:"ats"`
}
type Ats struct {
	App      App      `yaml:"app"`
	Database Database `yaml:"database"`
	Log      Log      `yaml:"log"`
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

type Log struct {
	Level string `yaml:"level"`
}