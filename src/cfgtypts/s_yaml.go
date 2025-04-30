package cfgtypts

type Config struct {
	Ats Ats `yaml:"ats"`
}

type Ats struct {
	App      App      `yaml:"app"`
	Database Database `yaml:"database"`
	Uias     Uias     `yaml:"uias"`
	Cron     Cron     `yaml:"cron"`
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

type Uias struct {
	Endpoint string `yaml:"endpoint"`
}

type Cron struct {
	Cleanlog struct {
		Time string `yaml:"time"`
		Days int    `yaml:"days"`
	} `yaml:"cleanlog"`
}

type Log struct {
	Level  string `yaml:"level"`
	Output Output `yaml:"output"`
}

type Output struct {
	File   File   `yaml:"file"`
	Stdout string `yaml:"stdout"`
}

type File struct {
	Name       string `yaml:"name"`
	MaxSize    int    `yaml:"maxsize"`
	MaxBackups int    `yaml:"maxbackups"`
	MaxAge     int    `yaml:"maxage"`
	Compress   bool   `yaml:"compress"`
}
