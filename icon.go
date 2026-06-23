package main

import (
	"bytes"
	"image"
	"image/png"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/getlantern/systray"
)

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
				systray.SetIcon(icon)
				return
			}
		}

	}()
	systray.SetIcon(icon)
	return func() { once.Do(func() { close(stopCh) }) }
}
