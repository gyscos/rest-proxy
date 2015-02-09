package main

import (
	"fmt"
	"strings"
)

func ensureHasPort(host string, defaultPort int) string {
	id := strings.Index(host, ":")
	if id < 0 {
		return fmt.Sprintf("%v:%v", host, defaultPort)
	} else if id == len(host)-1 {

		return fmt.Sprintf("%v%v", host, defaultPort)
	} else {
		return host
	}
}
