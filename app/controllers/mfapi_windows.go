// +build windows

package controllers

import (
	"os"
)

func getServerDir() string {

	dir := os.Getenv("PROGRAMDATA")
	if dir == "" {
		panic("Could not expand PROGRAMDATA var")
	}
	return dir + "/" + serverDir
}
