package fotos

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
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"github.com/chai2010/webp"
	"github.com/disintegration/gift"
	"github.com/jdeng/goheif"
	"github.com/lmittmann/ppm"

	"github.com/nfnt/resize"
)

var cl = collate.New(language.German)

var MediaFiles = []string{
	".gif", ".jpg", ".jpe", ".jpeg", ".jfif", ".jxl", ".png", ".ppm",
	".cr2", ".dng", ".nef", ".raw", ".arw", ".crw", ".mrw", ".raf", ".webp", ".heic", ".heif",
	".3gp", ".flv", ".mov", ".qt", ".m2ts", ".mts", ".divx", ".vob", ".webm", ".mkv", ".mka", ".wmv", ".avi", ".mp4",
	".mpg", ".mpeg", ".ps", ".ts", ".rm", ".ogv", ".dv",
	".bmp", ".ico", ".tiff", ".bigtiff", ".tiff85", ".pbm", ".pgm", ".dds", ".pcx", ".bpg", ".xbm", ".mac", ".tga", ".lmp",
}

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
	case ".jxl":
		if f, err = os.Open(filename); err != nil {
			return nil, fmt.Errorf("open image \"%v\" failed: %w", filename, err)
		}
		return DecodeJxl(filename)
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
	case ".ini", ".pano", ".html", ".log", ".db", ".zip", ".thumbnail", ".exe", ".tif", ".info", ".tlv", ".map", ".psd",
		".tar", ".rar", ".txt", ".ivp", ".rzs", ".dat", ".tmp", ".mrk", ".acv", ".atn", ".shh", ".bdm", ".tdt", ".tid",
		".xmp", ".golf":
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
	var outbuf bytes.Buffer
	cmd.Stdout = &outbuf
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("executing \"dcraw -c %v\" failed: %w", filename, err)
	}
	return ppm.Decode(bytes.NewReader(outbuf.Bytes()))
}

func DecodeVideo(filename string) (image.Image, error) {
	cmd := exec.Command("ffmpeg", "-i", filename, "-ss", "00:00:00.000", "-vframes", "1", "-f", "image2pipe", "-vcodec", "png", "-")
	var outbuf bytes.Buffer
	cmd.Stdout = &outbuf
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("executing \"ffmpeg\" with \"%v\" failed: %w", filename, err)
	}
	return png.Decode(bytes.NewReader(outbuf.Bytes()))
}

func DecodeJxl(filename string) (image.Image, error) {
	// Create a temporary file to hold the decoded image in PNG format.
	f, err := os.CreateTemp("", "fotos-*.png")
	if err != nil {
		return nil, fmt.Errorf("creating temporary file failed: %w", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	// Run the "djxl" command to decode the image and write it to the temporary file in PNG format.
	cmd := exec.Command("djxl", filename, f.Name())
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("executing \"djxl %v %v\" failed: %w", filename, f.Name(), err)
	}

	// Open the temporary file and decode the image from it.
	f, err = os.Open(f.Name())
	if err != nil {
		return nil, fmt.Errorf("open temporary file failed: %w", err)
	}
	defer f.Close()
	return png.Decode(f)
}

func EncodeJxl(m image.Image, filename string, quality int) error {
	// Create a temporary file to hold the image in PNG format.
	f, err := os.CreateTemp("", "fotos-*.png")
	if err != nil {
		return fmt.Errorf("creating temporary file failed: %w", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	// Encode the image in PNG format and write it to the temporary file.
	err = png.Encode(f, m)
	if err != nil {
		return fmt.Errorf("encoding image in PNG format failed: %w", err)
	}

	// Run the "cjxl" command to encode the image in JPEG-XL format.
	cmd := exec.Command("cjxl", f.Name(), filename, "-q", strconv.Itoa(quality))
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("executing \"cjxl %v %v -q %v\" failed: %w", f.Name(), filename, quality, err)
	}
	return nil
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

func Rotate(in image.Image, o int) image.Image {
	dim := in.Bounds()
	var r gift.Filter
	switch o {
	case 2:
		r = gift.FlipHorizontal()
	case 3:
		r = gift.Rotate180()
	case 4:
		r = gift.FlipVertical()
	case 5:
		r = gift.Transpose()
	case 6:
		r = gift.Rotate270()
	case 7:
		r = gift.Transverse()
	case 8:
		r = gift.Rotate90()
	default:
		return in
	}
	out := image.NewRGBA(r.Bounds(dim))
	r.Draw(out, in, nil)
	return out
}

func CheckFileAge(in os.FileInfo, outFile string) (bool, time.Time, error) {
	out, err := os.Stat(outFile + ".s.jxl")
	if err != nil {
		return false, in.ModTime(), fmt.Errorf("stat \"%v.s.jxl\" failed: %w", outFile, err)
	}
	return math.Abs(in.ModTime().Sub(out.ModTime()).Seconds()) < float64(5*time.Second), in.ModTime(), nil
}

func Run(inFile string, outFile string, fileInfo os.FileInfo) (Img, error) {
	img, err := Info(inFile, fileInfo)
	if img.N == "" {
		return img, fmt.Errorf("empty image information: %w", err)
	}
	if plausibility && img.Orientation > 4 && img.ExifH > img.ExifW {
		img.Orientation = 1
	}
	stats.LastImageName = inFile
	stats.LastImageSize = fileInfo.Size()
	format := strings.ToLower(filepath.Ext(fileInfo.Name()))
	m, err := Decode(inFile, format)
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
	if plausibility && ((img.Orientation > 4 && img.ExifH > img.ExifW) || format == ".cr2") {
		img.Orientation = 1
	}
	if img.Orientation > 4 {
		img.W, img.H = img.H, img.W
	}
	if img.Orientation > 1 {
		large = Rotate(large, img.Orientation)
	}
	m = nil
	if err := EncodeJxl(large, outFile+".h.jxl", 60); err != nil {
		return img, fmt.Errorf("encode image \"%v.h.jxl\" failed: %w", outFile, err)
	}
	_ = os.Chtimes(outFile+".h.jxl", img.ModTime, img.ModTime)
	small := Thumb(large, 400, 200)
	large = nil
	if err := EncodeJxl(small, outFile+".s.jxl", 20); err != nil {
		return img, fmt.Errorf("encode image \"%v.s.jxl\" failed: %w", outFile, err)
	}
	_ = os.Chtimes(outFile+".s.jxl", img.ModTime, img.ModTime)
	pico := Thumb(small, 4, 4)
	img.C = ImageCorners(pico)
	return img, nil
}
