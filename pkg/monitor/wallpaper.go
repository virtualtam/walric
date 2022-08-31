package monitor

import (
	"github.com/vcraescu/go-xrandr"
)

// WallpaperResolution returns the resolution constraints for an image to be
// suitable for a list of monitors.
func WallpaperResolution(monitors []xrandr.Monitor) *Resolution {
	return &Resolution{
		HeightPx: int(MaxHeight(monitors)),
		WidthPx:  int(MaxWidth(monitors)),
	}
}
