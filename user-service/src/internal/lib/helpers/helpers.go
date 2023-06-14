package helpers

import (
	"strconv"
)

func Itoa64(n int64) string {
	return strconv.FormatInt(n, 10)
}

func Atoi64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
