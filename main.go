package main

import (
	"fmt"
	"os"
	"path/filepath"

	"strings"

	_ "embed"
	"time"

	"bytes"
	"image"
	"image/png"

	"github.com/disintegration/imaging"
	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func rotateIcon(icon []byte) {
	img, _, _ := image.Decode(bytes.NewReader(icon))

	go func() {
		angle := 0.0
		ticker := time.NewTicker(120 * time.Millisecond)

		for range ticker.C {
			rotated := imaging.Rotate(img, angle, image.Transparent)

			var buf bytes.Buffer
			png.Encode(&buf, rotated)

			systray.SetIcon(buf.Bytes())

			angle += 30
			if angle >= 360 {
				angle = 0
			}
		}
	}()
}

func onReady() {
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

	go func() {
		for range mPath.ClickedCh {
			rotateIcon(iconData)

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

			fmt.Println("Arquivos HEIC encontrados:", len(heicFiles))

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
