package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Zecass/desktop-ps/image"
	"github.com/Zecass/desktop-ps/process"
	"github.com/Zecass/desktop-ps/settings"
	"github.com/Zecass/desktop-ps/wallpaper"
)

func main() {

	for {
		s, err := settings.GetSettings()
		if err != nil {
			fmt.Println(err)
			return
		}
		p, err := process.ListProcesses()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = image.GenerateImageFromProcess(p, s)
		if err != nil {
			fmt.Println(err)
			return
		}

		path, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}

		path = path + "\\wallpaper.png"
		err = wallpaper.SetWallpaper(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		time.Sleep(10 * time.Second)
	}
}
