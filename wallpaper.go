package walric

import (
	"github.com/vcraescu/go-xrandr"
	"github.com/virtualtam/walric/monitor"
)

// WallpaperResolution returns the resolution constraints for an image to be
// suitable for a list of monitors.
func WallpaperResolution(monitors []xrandr.Monitor) *monitor.Resolution {
	return &monitor.Resolution{
		HeightPx: int(monitor.MaxHeight(monitors)),
		WidthPx:  int(monitor.MaxWidth(monitors)),
	}
}
