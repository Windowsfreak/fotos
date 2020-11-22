package fotos

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/unicode/norm"
)

type Img struct {
	N           string    `json:"n"`
	W           int       `json:"w"`
	H           int       `json:"h"`
	D           time.Time `json:"d,omitempty"`
	C           string    `json:"c,omitempty"`
	Lat         float64   `json:"lat,omitempty"`
	Lon         float64   `json:"lon,omitempty"`
	Orientation int       `json:"-"`
	ExifW       int       `json:"-"`
	ExifH       int       `json:"-"`
}

func Info(filename string, info os.FileInfo) (img Img, mod time.Time, err error) {
	img.N = norm.NFC.String(info.Name())
	img.D = info.ModTime().Add(time.Duration(-info.ModTime().Nanosecond()))
	mod = info.ModTime()
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
	b, err := ioutil.ReadAll(out)
	defer cmd.Wait()
	if err != nil {
		err = fmt.Errorf("reading from pipe of \"exiftool\" with \"%v\" failed: %w", filename, err)
		return
	}
	data := strings.Split(strings.Trim(string(b), " \r\n"), "\t")
	layout := "2006:01:02 15:04:05"
	date, err := time.ParseInLocation(layout, data[0], time.Local)
	if err == nil {
		img.D = date
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
