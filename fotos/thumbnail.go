package fotos

import (
	"github.com/nfnt/resize"
	"image"
)

func MeshGradient(m image.Image) string {
	glyphs := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	var text = make([]byte, 32)
	n := 0
	for y := 0; y <= 3; y++ {
		for x := 0; x <= 3; x++ {
			c := m.At(x, y)
			r, g, b, _ := c.RGBA()
			colorVal := (r&0xF000)>>4 | (g&0xF000)>>8 | (b&0xF000)>>12
			text[n] = glyphs[colorVal>>6]
			n++
			text[n] = glyphs[colorVal&63]
			n++
		}
	}
	return string(text)
}

func Thumb(m image.Image, w uint, h uint) image.Image {
	return resize.Thumbnail(w, h, m, resize.Lanczos3)
}
