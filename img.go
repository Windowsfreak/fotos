package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/unicode/norm"
)

type Img struct {
	N   string    `json:"n"`
	W   int       `json:"w"`
	H   int       `json:"h"`
	D   time.Time `json:"d,omitempty"`
	C   string    `json:"c,omitempty"`
	Lat float64   `json:"lat,omitempty"`
	Lon float64   `json:"lon,omitempty"`
}

func Info(filename string, info os.FileInfo) (img Img, mod time.Time, err error) {
	img.N = norm.NFC.String(info.Name())
	img.D = info.ModTime().Add(time.Duration(-info.ModTime().Nanosecond()))
	mod = info.ModTime()
	cmd := exec.Command("exiftool", "-T", "-datetimeoriginal", "-gps:GPSLatitude", "-gps:GPSLongitude", filename)
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
	if data[1] == "-" {
		data[1] = ""
	}
	if data[2] == "-" {
		data[2] = ""
	}
	img.Lat = ddStr(data[1])
	img.Lon = ddStr(data[2])
	return
}

func ddStr(in string) float64 {
	in = strings.ReplaceAll(in, " deg ", " ")
	in = strings.ReplaceAll(in, "'", "")
	in = strings.ReplaceAll(in, "\"", "")
	in = strings.ReplaceAll(in, "º", "")
	s := strings.Split(in, " ")
	if len(s) < 3 {
		return 0
	}
	deg, _ := strconv.ParseFloat(s[0], 64)
	min, _ := strconv.ParseFloat(s[1], 64)
	sec, _ := strconv.ParseFloat(s[2], 64)
	return math.Round(dd(deg, min, sec)*1000000) / 1000000
}

func dd(deg float64, min float64, sec float64) float64 {
	return deg + (min / 60) + (sec / 3600)
}
