package fotos

import (
	"fmt"
	"strings"
	"time"
)

type Stats struct {
	NoFeed                bool
	Path                  string
	BytesRead             int64
	BytesSkipped          int64
	FoldersProcessed      int64
	FoldersIgnored        int64
	FoldersSkipped        int64
	FoldersDeleted        int64
	ImagesProcessed       int64
	ImagesIgnored         int64
	ImagesSkipped         int64
	ImagesDeleted         int64
	CurrentFolder         string
	CurrentFolderProgress int
	CurrentFolderFiles    int
	LastImageName         string
	LastImageSize         int64
	LastErrorName         string
	StartTime             time.Time
}

var stats = Stats{StartTime: time.Now()}

func feedLines() {
	println("\n\n\n\n\n\n\n\n\n\n\n")
}

func printStats(err error) {
	if !stats.NoFeed {
		fmt.Print(strings.Repeat("\r\033[K\033[A", 10))
	}
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Started %s, Statistics from %s\n", stats.StartTime.Format("02.01.2006 15:04:05"), time.Now().Format("02.01.2006 15:04:05"))
	fmt.Printf(" .------------.------.--------.------------.\n")
	fmt.Printf(" |            | Dirs | Images |        kiB |  Settings:\n")
	fmt.Printf(" | Processed: | %*d | %*d | %*d |  inFolder: %s\n", 4, stats.FoldersProcessed, 6, stats.ImagesProcessed, 10, stats.BytesRead/1024, inFolder)
	fmt.Printf(" |   Ignored: | %*d | %*d |            | outFolder: %s\n", 4, stats.FoldersIgnored, 6, stats.ImagesIgnored, outFolder)
	fmt.Printf(" |   Skipped: | %*d | %*d | %*d |      path: %s\n", 4, stats.FoldersSkipped, 6, stats.ImagesSkipped, 10, stats.BytesSkipped/1024, stats.Path)
	fmt.Printf(" |   Deleted: | %*d | %*d |            |  minFDate: %s\n", 4, stats.FoldersDeleted, 6, stats.ImagesDeleted, minFolderDate.Format("02.01.2006 15:04:05"))
	fmt.Printf(" '------------'------'--------'------------'\n")
	fmt.Printf("Current folder: %s (%d/%d)\n", stats.CurrentFolder, stats.CurrentFolderProgress, stats.CurrentFolderFiles)
	fmt.Printf("    Last image: %s (%d kiB)\n", stats.LastImageName, stats.LastImageSize/1024)
	fmt.Printf("    Last error: %s\n", stats.LastErrorName)
}
