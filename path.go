package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// lookupUser is a variable for user lookup function, so it can be mocked in tests.
var lookupUser = user.Lookup

func ExpandTilde(pth string) (string, error) {
	if pth == "" {
		return "", fmt.Errorf("No Path provided")
	}

	if strings.HasPrefix(pth, "~") {
		pth = strings.TrimPrefix(pth, "~")

		if len(pth) == 0 || strings.HasPrefix(pth, "/") {
			return os.ExpandEnv("$HOME" + pth), nil
		}

		splitPth := strings.Split(pth, "/")
		username := splitPth[0]

		usr, err := lookupUser(username)
		if err != nil {
			return "", err
		}

		pathInUsrHome := strings.Join(splitPth[1:], "/")

		return filepath.Join(usr.HomeDir, pathInUsrHome), nil
	}

	return pth, nil
}
