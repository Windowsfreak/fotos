package images

import (
	"bytes"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/disintegration/gift"
	"github.com/jdeng/goheif"
	"github.com/karmdip-mi/go-fitz"
	"github.com/lmittmann/ppm"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"strconv"
)

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
		".xmp", ".golf",
		".bmp", ".lnk", ".doc", ".7z", ".mp3", ".vcf", ".cpi", ".mpl", ".vpl", ".ai", ".dxf", ".emf", ".eps",
		".svg", ".stl":
		return nil, nil
	case ".pdf":
		var doc *fitz.Document
		doc, err = fitz.New(filename)
		if err != nil {
			return nil, fmt.Errorf("open document \"%v\" failed: %w", filename, err)
		}
		defer doc.Close()
		m, err = doc.Image(0)
		if err != nil {
			return nil, fmt.Errorf("rastering first page of document \"%v\" failed: %w", filename, err)
		}
		return
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
	if err := os.WriteFile(filename, buf.Bytes(), os.ModePerm); err != nil {
		return fmt.Errorf("writing file from \"webp.Encode\" with \"%v\" failed: %w", filename, err)
	}
	return nil
}

func EncodePNG(m image.Image, filename string) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, m); err != nil {
		return fmt.Errorf("calling \"png.Encode\" with \"%v\" failed: %w", filename, err)
	}
	if err := os.WriteFile(filename, buf.Bytes(), os.ModePerm); err != nil {
		return fmt.Errorf("writing file from \"png.Encode\" with \"%v\" failed: %w", filename, err)
	}
	return nil
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
