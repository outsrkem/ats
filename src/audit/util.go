package audit

import (
	"strconv"
)

func strToInt64(str string) (int64, error) {
	intValue, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}
