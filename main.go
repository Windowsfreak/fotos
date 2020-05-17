package main

import (
	"flag"
	"runtime"
	"strings"
	"time"
)

var (
	concurrency                chan struct{}
	inFolder                   = ""
	outFolder                  = ""
	minFolderDate              time.Time
	exclusions                 = map[string]struct{}{}
	alwaysProcessRotatedImages = false
	plausibility               = false
)

func main() {
	threads, path := shellArguments()
	stats.Path = path
	concurrency = make(chan struct{}, threads)
	p := strings.Split(path, "/")
	name := p[len(p)-1]
	if path == "" {
		name = "Fotos"
	}
	feedLines()
	printStats(nil)
	running := true
	go func() {
		for running {
			printStats(nil)
			time.Sleep(5 * time.Second)
		}
	}()
	_, _, err := Walk(name, path)
	if err != nil {
		printStats(err)
	}
	running = false
}

func shellArguments() (int, string) {
	threadsPtr := flag.Int("threads", runtime.NumCPU(), "the number of goroutines that are allowed to run concurrently")
	inFolderPtr := flag.String("in", "./volume/Fotos", "the root folder from which files are read")
	outFolderPtr := flag.String("out", "./volume/Thumbs", "the root folder where thumbnails are stored")
	pathPtr := flag.String("path", "", "the relative path from which photos are updated, no leading or trailing slashes")
	maxAgePtr := flag.Duration("maxage", 0, "the maximum age after which folders are scanned again, e.g. 48h")
	excludePtr := flag.String("exclude", "", "excluded folders, comma separated, e.g. snapshot")
	alwaysProcessRotatedImagesPtr := flag.Bool("always-process-rotated-images", false, "Do not skip unmodified images when exif rotation is set if rotated")
	plausibilityPtr := flag.Bool("plausibility", false, "Assume correct rotation when image is upright and requires 90 degree rotation")
	flag.Parse()
	threads := *threadsPtr
	inFolder = *inFolderPtr
	outFolder = *outFolderPtr
	path := *pathPtr
	maxAge := *maxAgePtr
	alwaysProcessRotatedImages = *alwaysProcessRotatedImagesPtr
	plausibility = *plausibilityPtr
	if maxAge > 0 {
		minFolderDate = time.Now().Add(-maxAge)
	}
	for _, v := range strings.Split(*excludePtr, ",") {
		exclusions[v] = struct{}{}
	}
	return threads, path
}
