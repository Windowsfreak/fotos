package main

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"github.com/chai2010/webp"
	"github.com/jdeng/goheif"
	"github.com/lmittmann/ppm"

	"github.com/nfnt/resize"
)

var cl = collate.New(language.German)

func Decode(filename string, format string) (m image.Image, err error) {
	var f *os.File
	defer f.Close()
	switch format {
	case ".gif":
		if f, err = os.Open(filename); err != nil {
			return nil, fmt.Errorf("open image \"%v\" failed: %w", filename, err)
		}
		return gif.Decode(f)
	case ".jpg", ".jpe", ".jpeg", ".jfif":
		if f, err = os.Open(filename); err != nil {
			return nil, fmt.Errorf("open image \"%v\" failed: %w", filename, err)
		}
		return jpeg.Decode(f)
	case ".png":
		if f, err = os.Open(filename); err != nil {
			return nil, fmt.Errorf("open image \"%v\" failed: %w", filename, err)
		}
		return png.Decode(f)
	case ".ppm":
		if f, err = os.Open(filename); err != nil {
			return nil, fmt.Errorf("open image \"%v\" failed: %w", filename, err)
		}
		return ppm.Decode(f)
	case ".cr2", ".dng", ".nef", ".raw", ".arw", ".crw", ".mrw", ".raf":
		if f, err = os.Open(filename); err != nil {
			return nil, fmt.Errorf("open image \"%v\" failed: %w", filename, err)
		}
		return DecodeDcraw(filename)
	case ".webp":
		if f, err = os.Open(filename); err != nil {
			return nil, fmt.Errorf("open image \"%v\" failed: %w", filename, err)
		}
		return webp.Decode(f)
	case ".heic", ".heif":
		if f, err = os.Open(filename); err != nil {
			return nil, fmt.Errorf("open image \"%v\" failed: %w", filename, err)
		}
		return goheif.Decode(f)
	case ".3gp", ".flv", ".mov", ".qt",
		".m2ts", ".mts", ".divx", ".vob",
		".webm", ".mkv", ".mka", ".wmv", ".avi", ".mp4",
		".mpg", ".mpeg", ".ps", ".ts", ".rm", ".ogv", ".dv":
		return DecodeVideo(filename)
	case ".ini", ".pano", ".html", ".log", ".db", ".zip", ".thumbnail", ".exe", ".tif", ".info", ".tlv", ".map", ".tar":
		return nil, nil
	default:
		if f, err = os.Open(filename); err != nil {
			return nil, fmt.Errorf("open image \"%v\" failed: %w", filename, err)
		}
		m, _, err = image.Decode(f)
		return
	}
}
func DecodeDcraw(filename string) (image.Image, error) {
	cmd := exec.Command("dcraw", "-c", filename)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("creating pipe for \"dcraw -c %v\" failed: %w", filename, err)
	}
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("executing \"dcraw -c %v\" failed: %w", filename, err)
	}
	defer out.Close()
	defer cmd.Wait()
	return ppm.Decode(out)
}
func DecodeVideo(filename string) (image.Image, error) {
	cmd := exec.Command("ffmpeg", "-i", filename, "-ss", "00:00:00.000", "-vframes", "1", "-f", "image2pipe", "-vcodec", "png", "-")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("creating pipe for \"ffmpeg\" with \"%v\" failed: %w", filename, err)
	}
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("executing \"ffmpeg\" with \"%v\" failed: %w", filename, err)
	}
	defer out.Close()
	defer cmd.Wait()
	return png.Decode(out)
}
func EncodeWebP(m image.Image, filename string, quality float32) error {
	var buf bytes.Buffer
	if err := webp.Encode(&buf, m, &webp.Options{Lossless: false, Quality: quality}); err != nil {
		return fmt.Errorf("calling \"webp.Encode\" with \"%v\" failed: %w", filename, err)
	}
	if err := ioutil.WriteFile(filename, buf.Bytes(), os.ModePerm); err != nil {
		return fmt.Errorf("writing file from \"webp.Encode\" with \"%v\" failed: %w", filename, err)
	}
	return nil
}

func Thumb(m image.Image, w uint, h uint) image.Image {
	return resize.Thumbnail(w, h, m, resize.Lanczos3)
}

func CheckFileAge(in os.FileInfo, outFile string) (bool, error) {
	out, err := os.Stat(outFile + ".s.webp")
	if err != nil {
		return false, fmt.Errorf("stat \"%v.s.webp\" failed: %w", outFile, err)
	}
	return math.Abs(in.ModTime().Sub(out.ModTime()).Seconds()) < float64(5*time.Second), nil
}

func Run(inFile string, outFile string, fileInfo os.FileInfo) (Img, error) {
	img, mod, err := Info(inFile, fileInfo)
	if img.N == "" {
		return img, fmt.Errorf("empty image information: %w", err)
	}
	stats.LastImageName = inFile
	stats.LastImageSize = fileInfo.Size()
	m, err := Decode(inFile, strings.ToLower(filepath.Ext(fileInfo.Name())))
	if err != nil {
		return img, fmt.Errorf("decode image \"%v\" failed: %w", inFile, err)
	}
	if m == nil {
		// ignored file format
		return Img{}, nil
	}
	bounds := m.Bounds()
	img.W = bounds.Dx()
	img.H = bounds.Dy()
	large := Thumb(m, 2048, 2048)
	m = nil
	if err := EncodeWebP(large, outFile+".h.webp", 60); err != nil {
		return img, fmt.Errorf("encode image \"%v.h.webp\" failed: %w", outFile, err)
	}
	_ = os.Chtimes(outFile+".h.webp", mod, mod)
	small := Thumb(large, 400, 200)
	large = nil
	if err := EncodeWebP(small, outFile+".s.webp", 20); err != nil {
		return img, fmt.Errorf("encode image \"%v.s.webp\" failed: %w", outFile, err)
	}
	_ = os.Chtimes(outFile+".s.webp", mod, mod)
	pico := Thumb(small, 4, 4)
	img.C = ImageCorners(pico)
	return img, nil
}
