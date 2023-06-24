package fotos

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
)

type Corners struct {
	nibbles [12]uint64
	count   uint64
}

func (c *Corners) Multiply(times uint64) {
	c.count *= times
	for i := range c.nibbles {
		c.nibbles[i] *= times
	}
}

func (c Corners) String() string {
	if c.count == 0 {
		return "000000000000000000000000"
	}
	var data [12]byte
	for i, v := range c.nibbles {
		data[i] = byte(v / c.count)
	}
	return hex.EncodeToString(data[:])
}

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

func NewCorners(in string) (Corners, error) {
	data, err := hex.DecodeString(in)
	if err != nil {
		return Corners{}, fmt.Errorf("invalid hex string \"%v\" in ImageCorners: %w", in, err)
	}
	corners := Corners{count: 1}
	for i, v := range data {
		corners.nibbles[i] = uint64(v)
	}
	return corners, nil
}

func SumCorners(in []Corners) Corners {
	corners := Corners{}
	for _, c := range in {
		for i, v := range c.nibbles {
			corners.nibbles[i] += v
		}
		corners.count += c.count
	}
	return corners
}
