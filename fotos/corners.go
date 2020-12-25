package fotos

import (
	"fmt"
	"image"
	"image/color"
)

func ImageCorner(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("%02x%02x%02x", r>>8, g>>8, b>>8)
}

func ImageCorners(m image.Image) string {
	r := m.Bounds()
	return ImageCorner(m.At(r.Min.X, r.Min.Y)) +
		ImageCorner(m.At(r.Max.X-1, r.Min.Y)) +
		ImageCorner(m.At(r.Min.X, r.Max.Y-1)) +
		ImageCorner(m.At(r.Max.X-1, r.Max.Y-1))
}
