package config

import (
	"ats/src/cfgtypts"
	"ats/src/slog"
	"flag"
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gopkg.in/yaml.v2"
)

type FlagArgs struct {
	CfgPath      string
	PrintVersion bool
	Plain        string // 接收命令行字符串，用于加密
}

var Cfg *cfgtypts.Config

func NewFlagArgs() *FlagArgs {
	fa := &FlagArgs{}
	flag.StringVar(&fa.CfgPath, "c", "ats.yaml", "Configuration file path")
	flag.BoolVar(&fa.PrintVersion, "version", false, "print program version")
	flag.StringVar(&fa.Plain, "encrypt", "", "Encrypted string.")
	flag.Parse()
	return fa
}

// InitConfig 初始化配置
func InitConfig() *cfgtypts.Config {
	klog := slog.FromContext(nil)
	var _cfg cfgtypts.Config
	fa := NewFlagArgs()
	if fa.PrintVersion {
		versions, _ := newVersions(Version, GoVersion, GitCommit)
		versions.Print(versions)
	}
	if fa.Plain != "" { // 加密命令行字符串
		encryption(fa.Plain)
	}
	klog.Infof("Read configuration file: %s", fa.CfgPath)
	configData, err := os.ReadFile(fa.CfgPath)
	if err != nil {
		hlog.Info("读取配置文件失败")
		os.Exit(1)
	}

	err = yaml.Unmarshal(configData, &_cfg)
	if err != nil {
		hlog.Info("配置解析失败")
		os.Exit(1)
	}
	decryptionDatabaseMysqlPwd(&_cfg)
	Cfg = &_cfg
	return Cfg
}
