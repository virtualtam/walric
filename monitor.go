package redwall

import (
	"errors"

	"github.com/vcraescu/go-xrandr"
)

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

func maxMonitorHeight(monitors []xrandr.Monitor) float32 {
	var maxHeight float32

	for _, monitor := range monitors {
		if monitor.Resolution.Height > maxHeight {
			maxHeight = monitor.Resolution.Height
		}
	}

	return maxHeight
}

func maxMonitorWidth(monitors []xrandr.Monitor) float32 {
	var maxWidth float32

	for _, monitor := range monitors {
		if monitor.Resolution.Width > maxWidth {
			maxWidth = monitor.Resolution.Width
		}
	}

	return maxWidth
}
