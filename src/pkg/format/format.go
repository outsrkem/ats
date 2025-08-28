package format

import (
	"strconv"
)

func Int2str(i int) string {
	str := strconv.Itoa(i)
	return str
}

// Str2Int converts a string to an integer.
func Str2Int(str string) (int, error) {
	n, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// Str2Int64 converts a string to an int64 type.
func Str2Int64(str string) int64 {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return -1
	}
	return n
}

// Int642Str converts an int64 type to a string.
func Int642Str(n int64) string {
	str := strconv.FormatInt(n, 10)
	return str
}
