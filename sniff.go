package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func sniff() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	downloads := filepath.Join(home, "Downloads")
	err = watcher.Add(downloads)

	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create {
				fmt.Println("Novo arquivo recebido:", event.Name)
			}
		case err := <-watcher.Errors:
			log.Println("erro:", err)
		}
	}
}
