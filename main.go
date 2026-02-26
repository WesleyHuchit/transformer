package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "embed"

	"github.com/disintegration/imaging"
	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func rotateIcon(icon []byte) (stop func()) {
	stopCh := make(chan struct{})
	var once sync.Once
	img, _, _ := image.Decode(bytes.NewReader(icon))

	go func() {
		ticker := time.NewTicker(120 * time.Millisecond)
		defer ticker.Stop()
		angle := 0.0

		for {
			select {
			case <-ticker.C:
				const iconSize = 20
				rotated := imaging.Rotate(img, angle, image.Transparent)
				rotated = imaging.Resize(rotated, iconSize, iconSize, imaging.Lanczos)

				var buf bytes.Buffer
				png.Encode(&buf, rotated)

				systray.SetIcon(buf.Bytes())

				angle += 30
				if angle >= 360 {
					angle = 0
				}

			case <-stopCh:
				systray.SetIcon(icon) // volta para o original
				return
			}
		}

	}()
	systray.SetIcon(icon)
	// notify("Sucesso", "Arquivos convertidos com sucesso")
	return func() { once.Do(func() { close(stopCh) }) }
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

	mPath := systray.AddMenuItem("Copy", "Copy path")
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
				fmt.Println("Erro:", err) // ou mostrar no menu/notificação
				stop()
				continue
			}

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

			// fmt.Println("Arquivos HEIC encontrados:", len(heicFiles))

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
