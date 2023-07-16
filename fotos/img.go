package fotos

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/unicode/norm"
)

type Img struct {
	Path        string  `json:"p"`           // basename of encrypted img file
	Name        string  `json:"n"`           // original filename
	Width       int     `json:"w"`           // width
	Height      int     `json:"h"`           // height
	Date        int64   `json:"d,omitempty"` // EXIF date
	Color       string  `json:"c,omitempty"` // string representing a 4x4 low-res mesh gradient
	Lat         float64 `json:"lat,omitempty"`
	Lon         float64 `json:"lon,omitempty"`
	Orientation int     `json:"-"`
	ExifW       int     `json:"-"`
	ExifH       int     `json:"-"`
	ModTime     int64   `json:"-"`
}

func FindImg(items []Img, name string) (Img, bool) {
	for i := range items {
		if items[i].Name == name {
			return items[i], true
		}
	}
	return Img{}, false
}

func Info(filename string, info os.FileInfo) (img Img, err error) {
	img.Name = norm.NFC.String(info.Name())
	img.Date = info.ModTime().Unix()
	img.ModTime = img.Date
	cmd := exec.Command("exiftool", "-T", "-datetimeoriginal", "-orientation", "-gps:GPSLatitude", "-gps:GPSLongitude", "-imagewidth", "-imageheight", "-n", filename)
	out, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("creating pipe for \"exiftool\" with \"%v\" failed: %w", filename, err)
		return
	}
	err = cmd.Start()
	if err != nil {
		err = fmt.Errorf("executing \"exiftool\" with \"%v\" failed: %w", filename, err)
		return
	}
	defer out.Close()
	b, err := io.ReadAll(out)
	defer cmd.Wait()
	if err != nil {
		err = fmt.Errorf("reading from pipe of \"exiftool\" with \"%v\" failed: %w", filename, err)
		return
	}
	data := strings.Split(strings.Trim(string(b), " \r\n"), "\t")
	layout := "2006:01:02 15:04:05"
	date, err := time.ParseInLocation(layout, data[0], time.Local)
	if err == nil {
		img.Date = date.Unix()
	}
	if len(data) < 2 { // Assume error "No matching files"
		err = fmt.Errorf("no matching files")
		return
	}

	img.Orientation = orientation(data[1])
	img.Lat = ddStr(data[2])
	img.Lon = ddStr(data[3])
	img.ExifW, _ = strconv.Atoi(data[4])
	img.ExifH, _ = strconv.Atoi(data[5])
	return
}

func orientation(in string) int {
	val, err := strconv.Atoi(in)
	if err != nil {
		return 1
	}
	return val
}
func ddStr(in string) float64 {
	val, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return 0
	}
	return val
}
