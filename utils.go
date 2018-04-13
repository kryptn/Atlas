package main

import (
	"strings"
)

func Split2(str, separator string) (string, string) {
	sep := strings.SplitN(str, separator, 2)
	return sep[0], sep[1]
}
