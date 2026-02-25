package main

import (
	"fmt"
	"os"
	"path/filepath"

	"os/exec"
	"strings"

	"github.com/getlantern/systray"
)

func main() {
	fmt.Println("Hello, World!")

	systray.Run(onReady, onExit)
}

func onReady() {
	fmt.Println("Open")
	iconData, err := os.ReadFile("icon.png")

	if err != nil {
		fmt.Println("Error")
		return
	}

	systray.SetIcon(iconData)
	// systray.SetTitle("Transformer")
	systray.SetTooltip("Converta HEIC para JPG")

	mPath := systray.AddMenuItem("Copy", "Copy path")
	mQuit := systray.AddMenuItem("Sair", "Encerrar o app")

	// go func() {
	// 	<-mPath.ClickedCh
	// 	fmt.Println("Path")
	// }()

	go func() {
		for range mPath.ClickedCh {

			path, err := getFrontmostFinderPath()

			if err != nil {
				fmt.Println("Erro:", err) // ou mostrar no menu/notificação
				continue
			}

			entries, err := os.ReadDir(path)

			if err != nil {
				fmt.Println("Erro ao listar pasta:", err)
				continue
			}

			var heicFiles []string

			for _, e := range entries {
				if e.IsDir() {
					continue
				}

				if strings.EqualFold(filepath.Ext(e.Name()), ".heic") {
					fullPath := filepath.Join(path, e.Name())
					heicFiles = append(heicFiles, fullPath)
				}

			}

			// fmt.Println("Arquivos HEIC encontrados:", len(entries))
			fmt.Println("Arquivos HEIC encontrados:", len(heicFiles))
			// for _, f := range heicFiles {
			// 	fmt.Println(f)
			// }

			// fmt.Println("Pasta atual no Finder:", path)
			heicFolderName := "heic"
			heicFolderPath := filepath.Join(path, heicFolderName)

			if err := os.MkdirAll(heicFolderPath, 0755); err != nil {
				fmt.Println("Erro ao criar pasta:", err)
				continue
			}

			for _, heicPath := range heicFiles {
				_, err := heicToJPG(heicPath)
				if err != nil {
					fmt.Println("Erro ao converter", heicPath, err)
					continue
				}
				// fmt.Println("Convertido:", heicPath, "->", jpgPath)
			}

			for _, heicPath := range entries {
				fileName := filepath.Base(filepath.Join(path, heicPath.Name()))
				destPath := filepath.Join(heicFolderPath, fileName)
				if err := os.Rename(filepath.Join(path, heicPath.Name()), destPath); err != nil {
					fmt.Println("Erro ao mover", fileName, err)
					continue
				}
				// fmt.Println("Movido:", fileName, "->", heicFolderPath)
			}

		}
	}()

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

}

func onExit() {
	fmt.Println("Exit")
}

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

func heicToJPG(heicPath string) (string, error) {
	ext := filepath.Ext(heicPath)
	jpgPath := strings.TrimSuffix(heicPath, ext) + ".jpg"
	cmd := exec.Command("sips", "-s", "format", "jpeg", heicPath, "--out", jpgPath)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return jpgPath, nil
}
