package config

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type FlagArgs struct {
	CfgPath      string
	PrintVersion bool
}

func NewFlagArgs() *FlagArgs {
	fa := &FlagArgs{}
	flag.StringVar(&fa.CfgPath, "c", "ats.yaml", "Configuration file path")
	flag.BoolVar(&fa.PrintVersion, "version", false, "print program version")
	flag.Parse()
	return fa
}

// InitConfig 初始化配置
func InitConfig() *Config {
	var _cfg Config
	fa := NewFlagArgs()


	log.Println("Read configuration file:", fa.CfgPath)

	configData, err := os.ReadFile(fa.CfgPath)
	if err != nil {
		log.Println("读取配置文件失败:", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(configData, &_cfg)
	if err != nil {
		log.Println("解析配置文件失败:", err)
		os.Exit(1)
	}

	return &_cfg
}