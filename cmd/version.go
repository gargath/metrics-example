package main

import (
	"fmt"
	"strings"
)

// VERSION holds the version of the code, injected at compile time
var VERSION = "undefined"

func version() string {
	if strings.Contains(VERSION, ".") {
		return fmt.Sprintf("version %s", VERSION)
	}
	return fmt.Sprintf("build sha:%s", VERSION)
}
