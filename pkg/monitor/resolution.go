package monitor

// Resolution represents a monitor's resolution, in pixels.
type Resolution struct {
	HeightPx int
	WidthPx  int
}

// Validate ensures this Resolution is valid.
func (r *Resolution) Validate() error {
	if r.HeightPx < 1 || r.WidthPx < 1 {
		return ErrResolutionInvalid
	}

	return nil
}
