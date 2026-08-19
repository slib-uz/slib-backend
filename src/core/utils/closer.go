package utils

import (
	"io"
)

func Closer(c io.Closer) func() {
	return func() {
		_ = c.Close()
	}
}
