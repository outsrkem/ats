package config

import (
	"ats/src/pkg/crypto"
	"fmt"
	"os"
)

func encryption(plain string) {
	fmt.Println(crypto.Encryption(plain))
	os.Exit(0)
}
