package redwall

import "github.com/vcraescu/go-xrandr"

type Resolution struct {
	HeightPx int
	WidthPx  int
}

func WallpaperResolution(monitors []xrandr.Monitor) *Resolution {
	return &Resolution{
		HeightPx: int(maxMonitorHeight(monitors)),
		WidthPx:  int(maxMonitorWidth(monitors)),
	}
}
