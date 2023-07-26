package images

import (
	"bytes"
	"fmt"
	"github.com/disintegration/gift"
	"github.com/lmittmann/ppm"
	"image"
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

var PyVipsFiles = []string{
	".gif", ".jpg", ".jpe", ".jpeg", ".jfif",
	".jxl", ".png", ".webp", ".pdf", ".svg", ".ppm",
	".tif", ".tiff", ".heic", ".heif", ".avif",
	".mat", ".v", ".vips", ".img", ".hdr",
	".pbm", ".pgm", ".ppm", ".pfm", ".pnm",
	".svg", ".svgz", ".svg.gz",
	".j2k", ".jp2", ".jpt", ".j2c", ".jpc",
	".fits", ".fit", ".fts",
	".exr", ".svs", ".vms", ".vmu", ".ndpi", ".scn", ".mrxs", ".svslide", ".bif",
	".bpg", ".bmp", ".dib", ".dcm", ".emf",
}

func Decode(filename string, format string) (m image.Image, err error) {
	var f *os.File
	defer f.Close()
	switch format {
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
	case ".3gp", ".flv", ".mov", ".qt",
		".m2ts", ".mts", ".divx", ".vob",
		".webm", ".mkv", ".mka", ".wmv", ".avi", ".mp4",
		".mpg", ".mpeg", ".ps", ".ts", ".rm", ".ogv", ".dv":
		return DecodeVideo(filename)
	case ".ini", ".pano", ".html", ".log", ".db", ".zip", ".thumbnail", ".exe", ".tif", ".info", ".tlv", ".map", ".psd",
		".tar", ".rar", ".txt", ".ivp", ".rzs", ".dat", ".tmp", ".mrk", ".acv", ".atn", ".shh", ".bdm", ".tdt", ".tid",
		".xmp", ".golf", ".pto",
		".bmp", ".lnk", ".doc", ".7z", ".mp3", ".vcf", ".cpi", ".mpl", ".vpl", ".ai", ".dxf", ".emf", ".eps",
		".svg", ".stl",
		".json", ".prproj", ".m4a", ".cfa", ".pek", ".srt":
		return nil, nil
	case ".pdf":
		return DecodePDF(filename)
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

func DecodePDF(filename string) (image.Image, error) {
	// Run mutool draw command to convert the first page of the PDF to PNG format
	cmd := exec.Command("mutool", "draw", "-o", "-", filename, "1")
	var outbuf bytes.Buffer
	cmd.Stdout = &outbuf
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("executing \"mutool\" with \"%v\" failed: %w", filename, err)
	}

	// Decode the PNG image from the output of mutool
	return png.Decode(&outbuf)
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
