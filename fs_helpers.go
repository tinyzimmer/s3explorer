package main

import "os"

func FileExists(path string) bool {
	if _, err := os.Stat(logFile); err == nil {
		return true
	}
	return false
}
