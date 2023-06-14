package helpers

import (
	"fmt"
	"strconv"
)

func Itoa64(n int64) string {
	return strconv.FormatInt(n, 10)
}

func Atoi64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func RunOnBackground(function func()) {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()
		function()
	}()
}
