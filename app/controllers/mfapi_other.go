// +build !windows

package controllers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func getServerDir() string {

	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	// If that fails, try getent
	var stdout bytes.Buffer
	cmd := exec.Command("getent", "passwd", strconv.Itoa(os.Getuid()))
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		// If the error is ErrNotFound, we ignore it. Otherwise, return it.
		if err != exec.ErrNotFound {
			fmt.Println("getServerDir Error! Failed running getent")
			panic(err)
		}
	} else {
		if passwd := strings.TrimSpace(stdout.String()); passwd != "" {
			// username:password:uid:gid:gecos:home:shell
			passwdParts := strings.SplitN(passwd, ":", 7)
			if len(passwdParts) > 5 {
				return passwdParts[5] + "/" + serverDir
			}
		}
	}

	// If all else fails, try the shell
	stdout.Reset()
	cmd = exec.Command("sh", "-c", "cd && pwd")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("getServerDir Error! Failed running sh -c pwd")
		panic(err)
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		panic("getServerDir Error! Blank output for home directory")
	}

	return result + "/" + serverDir
}
