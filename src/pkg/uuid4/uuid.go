package uuid4

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

func Uuid4Str() string {
	u4 := uuid.NewV4().String()
	return strings.ReplaceAll(u4, "-", "")
}
