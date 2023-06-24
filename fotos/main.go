package fotos

import (
	"flag"
	"runtime"
	"strings"
	"time"

	"fotos/domain"
)

var (
	concurrency                chan struct{}
	inFolder                   = ""
	outFolder                  = ""
	minFolderDate              time.Time
	exclusions                 = map[string]struct{}{}
	alwaysProcessRotatedImages = false
	plausibility               = false
	invalidatePaths            []string
	Running                    = false
	Repeat                     = false
	RepeatConfig               domain.ConfigStruct
	rootFolderCaption          = ""
)

func TryRun(config domain.ConfigStruct) {
	if Running {
		Repeat = true
		i := RepeatConfig.InvalidatePaths
		RepeatConfig = config
		for _, s := range config.InvalidatePaths {
			i = append(i, s)
		}
		RepeatConfig.InvalidatePaths = i
	} else {
		runCustom(config)
	}
}

func runCustom(config domain.ConfigStruct) {
	if config.MaxAge > 0 {
		minFolderDate = time.Now().Add(-config.MaxAge)
	}
	for _, v := range strings.Split(config.Exclude, ",") {
		exclusions[v] = struct{}{}
	}
	inFolder = config.InFolder
	outFolder = config.OutFolder
	path := config.Path
	alwaysProcessRotatedImages = config.AlwaysProcessRotatedImages
	plausibility = config.Plausibility
	invalidatePaths = config.InvalidatePaths
	stats.Path = config.Path
	concurrency = make(chan struct{}, config.Threads)
	rootFolderCaption = config.RootFolderCaption
	p := strings.Split(path, "/")
	name := p[len(p)-1]
	if path == "" {
		name = rootFolderCaption
	}
	stats.NoFeed = true
	stats.StartTime = time.Now()
	printStats(nil)
	Running = true
	go func() {
		for Running {
			printStats(nil)
			time.Sleep(5 * time.Second)
		}
	}()
	_, _, err := Walk(name, path)
	printStats(err)
	Running = false
	if Repeat {
		Repeat = false
		c := RepeatConfig
		RepeatConfig.InvalidatePaths = []string{}
		runCustom(c)
	}
}

func Main() {
	threads, path := shellArguments()
	stats.Path = path
	concurrency = make(chan struct{}, threads)
	p := strings.Split(path, "/")
	name := p[len(p)-1]
	if path == "" {
		name = rootFolderCaption
	}
	feedLines()
	running := true
	go func() {
		time.Sleep(5 * time.Second)
		for running {
			printStats(nil)
			time.Sleep(5 * time.Second)
		}
	}()
	_, _, err := Walk(name, path)
	printStats(err)
	running = false
}

func shellArguments() (int, string) {
	threadsPtr := flag.Int("threads", runtime.NumCPU(), "the number of goroutines that are allowed to run concurrently")
	inFolderPtr := flag.String("in", "./volume/Fotos", "the root folder from which files are read")
	outFolderPtr := flag.String("out", "./volume/Thumbs", "the root folder where thumbnails are stored")
	pathPtr := flag.String("path", "", "the relative path from which photos are updated, no leading or trailing slashes")
	maxAgePtr := flag.Duration("maxage", 0, "the maximum age after which folders are scanned again, e.g. 48h")
	excludePtr := flag.String("exclude", "", "excluded folders, comma separated, e.g. snapshot")
	alwaysProcessRotatedImagesPtr := flag.Bool("always-process-rotated-images", false, "Do not skip unmodified rotated CR2 files")
	plausibilityPtr := flag.Bool("plausibility", false, "Assume correct rotation when image is upright and requires 90 degree rotation")
	rootFolderCaptionPtr := flag.String("root", "Fotos", "name root folder, e.g. Fotos")
	flag.Parse()
	threads := *threadsPtr
	inFolder = *inFolderPtr
	outFolder = *outFolderPtr
	path := *pathPtr
	maxAge := *maxAgePtr
	alwaysProcessRotatedImages = *alwaysProcessRotatedImagesPtr
	plausibility = *plausibilityPtr
	rootFolderCaption = *rootFolderCaptionPtr
	if maxAge > 0 {
		minFolderDate = time.Now().Add(-maxAge)
	}
	for _, v := range strings.Split(*excludePtr, ",") {
		exclusions[v] = struct{}{}
	}
	return threads, path
}
