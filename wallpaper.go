package redwall

import (
	"github.com/vcraescu/go-xrandr"
	"github.com/virtualtam/redwall2/monitor"
)

func WallpaperResolution(monitors []xrandr.Monitor) *monitor.Resolution {
	return &monitor.Resolution{
		HeightPx: int(monitor.MaxHeight(monitors)),
		WidthPx:  int(monitor.MaxWidth(monitors)),
	}
}
