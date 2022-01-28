package monitor

import (
	"errors"

	"github.com/vcraescu/go-xrandr"
)

// ConnectedMonitors returns the list of connected (active) monitors for a given
// XRandR screen identifier.
func ConnectedMonitors(xRandRScreenNo int) ([]xrandr.Monitor, error) {
	screens, err := xrandr.GetScreens()
	if err != nil {
		return []xrandr.Monitor{}, err
	}

	for _, screen := range screens {
		if screen.No != xRandRScreenNo {
			continue
		}

		var monitors []xrandr.Monitor

		for _, monitor := range screen.Monitors {
			if !monitor.Connected {
				continue
			}

			monitors = append(monitors, monitor)
		}

		return monitors, nil
	}

	return []xrandr.Monitor{}, errors.New("Screen not found")
}

// MaxHeight returns the maximum height (in pixels) of a list of XRandR
// monitors.
func MaxHeight(monitors []xrandr.Monitor) float32 {
	var maxHeight float32

	for _, monitor := range monitors {
		if monitor.Resolution.Height > maxHeight {
			maxHeight = monitor.Resolution.Height
		}
	}

	return maxHeight
}

// MaxWidth returns the maximum width (in pixels) of a list of XRandR
// monitors.
func MaxWidth(monitors []xrandr.Monitor) float32 {
	var maxWidth float32

	for _, monitor := range monitors {
		if monitor.Resolution.Width > maxWidth {
			maxWidth = monitor.Resolution.Width
		}
	}

	return maxWidth
}
