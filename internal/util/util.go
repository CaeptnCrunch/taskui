package util

import (
	"os/user"
	"path/filepath"
)

func ExpandUnixPath(path string) string {
	if len(path) < 2 {
		return path
	}
	fixedPath := path
	if fixedPath[:2] == "~/" {
		userDir, _ := user.Current()
		homeDir := userDir.HomeDir
		fixedPath = filepath.Join(homeDir, fixedPath[2:])
	}
	return fixedPath
}
