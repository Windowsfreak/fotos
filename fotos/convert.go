package fotos

import (
	"fmt"
	"fotos/fotos/images"
	"github.com/nfnt/resize"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const plausibility = true

func Convert(inFile string, outFile string, fileInfo fs.FileInfo) (Img, error) {
	img, err := Info(inFile, fileInfo)
	if img.Name == "" {
		return img, fmt.Errorf("empty image information: %w", err)
	}
	if plausibility && img.Orientation > 4 && img.ExifH > img.ExifW {
		img.Orientation = 1
	}
	format := strings.ToLower(filepath.Ext(fileInfo.Name()))
	m, err := images.Decode(inFile, format)
	if err != nil {
		return img, fmt.Errorf("decode image \"%v\" failed: %w", inFile, err)
	}
	if m == nil {
		// ignored file format
		return Img{}, nil
	}
	if plausibility && ((img.Orientation > 4 && img.ExifH > img.ExifW) || format == ".cr2") {
		img.Orientation = 1
	}
	/*if img.Orientation > 4 {
		img.Width, img.Height = img.Height, img.Width
	}*/
	if img.Orientation > 1 {
		m = images.Rotate(m, img.Orientation)
	}
	if err := images.EncodeJxl(m, outFile+".o.jxl", 75); err != nil {
		return img, fmt.Errorf("encode image \"%v.o.jxl\" failed: %w", outFile, err)
	}
	_ = os.Chtimes(outFile+".o.jxl", time.Unix(img.ModTime, 0), time.Unix(img.ModTime, 0))
	bounds := m.Bounds()
	img.Width = bounds.Dx()
	img.Height = bounds.Dy()
	large := Thumb(m, 2048, 2048)
	m = nil
	if err := images.EncodeJxl(large, outFile+".h.jxl", 60); err != nil {
		return img, fmt.Errorf("encode image \"%v.h.jxl\" failed: %w", outFile, err)
	}
	_ = os.Chtimes(outFile+".h.jxl", time.Unix(img.ModTime, 0), time.Unix(img.ModTime, 0))
	small := Thumb(large, 400, 200)
	large = nil
	if err := images.EncodeJxl(small, outFile+".s.jxl", 20); err != nil {
		return img, fmt.Errorf("encode image \"%v.s.jxl\" failed: %w", outFile, err)
	}
	_ = os.Chtimes(outFile+".s.jxl", time.Unix(img.ModTime, 0), time.Unix(img.ModTime, 0))
	pico := resize.Resize(4, 4, small, resize.Lanczos3)
	img.Color = MeshGradient(pico)
	return img, nil
}

func RemoveThumbnails(outFile string) {
	// println("Here I would delete the file " + outFile + " with .o.jxl, .h.jxl and .s.jxl")
	if err := os.Remove(outFile + ".o.jxl"); err != nil {
		println(fmt.Errorf("deleting image \"%v.o.jxl\" failed: %w", outFile, err).Error())
	}
	if err := os.Remove(outFile + ".h.jxl"); err != nil {
		println(fmt.Errorf("deleting image \"%v.h.jxl\" failed: %w", outFile, err).Error())
	}
	if err := os.Remove(outFile + ".s.jxl"); err != nil {
		println(fmt.Errorf("deleting image \"%v.s.jxl\" failed: %w", outFile, err).Error())
	}
}
