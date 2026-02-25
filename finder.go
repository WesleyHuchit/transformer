package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func getFrontmostFinderPath() (string, error) {
	script := `tell application "Finder"
        if (count of windows) > 0 then
            set folderPath to POSIX path of (target of front window as alias)
            return folderPath
        else
            return ""
        end if
    end tell`
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return "", err
	}
	path := strings.TrimSpace(string(out))
	if path == "" {
		return "", fmt.Errorf("nenhuma janela do Finder aberta")
	}
	return path, nil
}
