package utilities

import "os"

func GetHostName() string {
	if name, err := os.Hostname(); err != nil {
		return ""
	} else {
		return name
	}
}
