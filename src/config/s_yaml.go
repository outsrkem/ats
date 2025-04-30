package config

import (
	"ats/src/cfgtypts"
	"ats/src/pkg/crypto"
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func decryptionDatabaseMysqlPwd(c *cfgtypts.Config) {
	if c.Ats.Database.Passwd == "" {
		return
	}

	plain, err := crypto.Decryption(c.Ats.Database.Passwd)
	if err != nil {
		hlog.Fatalf("Decryption of database password failed. uias.yaml:uias.database.passwd: %s", c.Ats.Database.Passwd)
		os.Exit(100)
	}

	c.Ats.Database.Passwd = plain
}
