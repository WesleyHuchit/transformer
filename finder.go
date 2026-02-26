package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
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
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

func setFinderSortByName() error {
	script := `tell application "Finder"
        if (count of windows) > 0 then
            set current view of front window to list view
            set sort column of list view options of front window to name column
        end if
    end tell`
	return exec.Command("osascript", "-e", script).Run()
}
