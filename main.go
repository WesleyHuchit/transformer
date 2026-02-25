package main

import (
	"fmt"
	"os"

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
			fmt.Println("Pasta atual no Finder:", path)
			// usar path (ex.: listar HEICs e converter)

			folderPath := path

			entries, err := os.ReadDir(folderPath)
			if err != nil {
				fmt.Println("Erro ao listar pasta:", err)
				continue
			}

			for _, e := range entries {
				info, _ := e.Info()
				var size int64
				if info != nil {
					size = info.Size()
				}
				if e.IsDir() {
					fmt.Println("[DIR]", e.Name(), size)
				} else {
					fmt.Println("[FILE]", e.Name(), size)
				}
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
