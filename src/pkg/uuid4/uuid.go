package uuid4

import (
	"strings"

	"github.com/google/uuid"
)

func Uuid4Str() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
