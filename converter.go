package main

import (
	"os/exec"
	"path/filepath"
	"strings"
)

func heicToJPG(heicPath string) (string, error) {
	ext := filepath.Ext(heicPath)
	jpgPath := strings.TrimSuffix(heicPath, ext) + ".jpg"
	cmd := exec.Command("sips", "-s", "format", "jpeg", heicPath, "--out", jpgPath)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return jpgPath, nil
}
