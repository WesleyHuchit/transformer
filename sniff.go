package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

func sniff() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add("/Users/wesleyhuchit/Downloads")
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
