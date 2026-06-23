package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func resourcePath(filename string) string {
	exe, err := os.Executable()
	if err != nil {
		return filepath.Join("images", filename)
	}
	exe, _ = filepath.EvalSymlinks(exe)
	dir := filepath.Dir(exe)
	if strings.Contains(dir, ".app/Contents/MacOS") {
		return filepath.Join(dir, "..", "Resources", filename)
	}
	return filepath.Join("images", filename)
}

func onReady() {
	iconPath := resourcePath("icon.png")
	iconData, err := os.ReadFile(iconPath)

	if err != nil {
		fmt.Println("Error loading icon:", iconPath, err)
		return
	}

	systray.SetIcon(iconData)
	systray.SetTooltip("Converta HEIC para JPG")

	go sniff()

	mPath := systray.AddMenuItem("Converter", "Copy path")
	mQuit := systray.AddMenuItem("Sair", "Encerrar o app")

	loadingPath := resourcePath("loading.png")
	loadingIcon, err := os.ReadFile(loadingPath)

	if err != nil {
		fmt.Println("Error reading loading icon:", loadingPath, err)
		return
	}

	go func() {
		for range mPath.ClickedCh {

			stop := rotateIcon(loadingIcon)

			path, err := getFrontmostFinderPath()

			if err != nil {
				fmt.Println("Erro:", err)
				stop()
				continue
			}

			// systray.SetTitle(path)

			_ = setFinderSortByName()

			entries, err := os.ReadDir(path)

			if err != nil {
				fmt.Println("Erro ao listar pasta:", err)
				stop()
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

			if len(heicFiles) == 0 {
				notify("Erro", "Nenhum arquivo HEIC encontrado")
				stop()
				systray.SetIcon(iconData)
				return
			}

			heicFolderName := "heic"
			heicFolderPath := filepath.Join(path, heicFolderName)

			if err := os.MkdirAll(heicFolderPath, 0755); err != nil {
				fmt.Println("Erro ao criar pasta:", err)
				stop()
				continue
			}

			for _, heicPath := range heicFiles {
				_, err := heicToJPG(heicPath)
				if err != nil {
					fmt.Println("Erro ao converter", heicPath, err)
					stop()
					continue
				}
				// fmt.Println("Convertido:", heicPath, "->", jpgPath)
			}

			for _, heicPath := range heicFiles {
				fileName := filepath.Base(heicPath)
				destPath := filepath.Join(heicFolderPath, fileName)
				if err := os.Rename(heicPath, destPath); err != nil {
					fmt.Println("Erro ao mover", fileName, err)
					stop()
					continue
				}
				// fmt.Println("Movido:", fileName, "->", heicFolderPath)
			}

			// systray.SetTooltip("Conversão concluída!")
			stop()
			systray.SetIcon(iconData)
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

func notify(title, message string) {
	cmd := exec.Command("osascript", "-e",
		fmt.Sprintf(`display notification "%s" with title "%s"`,
			message, title))
	cmd.Run()
}
