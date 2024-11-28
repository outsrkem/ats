package common

import "regexp"

func CheckUuId(uuid string) bool {
	if uuid == "" {
		return false
	}
	if !regexp.MustCompile(`^[0-9a-fA-F]{32}$`).MatchString(uuid) {
		return false
	}
	return true
}
