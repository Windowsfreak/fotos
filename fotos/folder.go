package fotos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/text/unicode/norm"

	"fotos/repository"
)

func getInfo(folder string) (OldDir, error) {
	myOutFolder := outFolder
	if folder != "" {
		if folder[0] == '/' {
			folder = folder[1:]
		}
		myOutFolder += "/" + folder
	}
	data, err := ioutil.ReadFile(myOutFolder + "/index.json")
	if err != nil {
		return Dir{}.MakeOld(), fmt.Errorf("file \"index.json\" in \"%v\" could not be read: %w", myOutFolder, err)
	}
	p := Dir{}
	err = json.Unmarshal(data, &p)
	if err != nil {
		return Dir{}.MakeOld(), fmt.Errorf("file \"index.json\" in \"%v\" could not be decoded: %w", myOutFolder, err)
	}
	return p.MakeOld(), nil
}

func StripLeadingSlash(text string) string {
	if text[0] == '/' {
		text = text[1:]
	}
	return text
}

// CheckFolderAge returns true if folder is still fresh enough
func CheckFolderAge(name string) (bool, error) {
	in, err := os.Stat(name)
	if err != nil {
		return false, fmt.Errorf("stat \"%v\" failed: %w", name, err)
	}
	return in.ModTime().After(minFolderDate), nil
}

func Walk(name string, folder string, r repository.Repository) (Dir, error) {
	myInFolder := inFolder
	myOutFolder := outFolder
	if folder != "" {
		folder = StripLeadingSlash(folder)
		myInFolder += "/" + folder
		myOutFolder += "/" + folder
	}
	d := name
	if r != nil {
		num, err := repository.ConvertStringToInt64(name)
		if err == nil {
			username, discriminator, err := r.Fetch(num)
			if err == nil {
				d = username + "#" + discriminator
			}
		}
	}
	myDir := Dir{MinDir: MinDir{D: d, N: name, Oldest: time.Now()}, Path: folder}
	if err := os.MkdirAll(myOutFolder, os.ModePerm); err != nil {
		return myDir, fmt.Errorf("could not create folder structure for \"%v\": %w", myOutFolder, err)
	}
	files, err := ioutil.ReadDir(myInFolder)
	if err != nil {
		return myDir, fmt.Errorf("could not list contents of folder \"%v\": %w", myInFolder, err)
	}
	oldDir, _ := getInfo(folder)
	err = WalkSubdirectories(files, folder, oldDir.Subs, myOutFolder, &myDir, r)
	if err != nil {
		return myDir, err
	}
	stats.CurrentFolder = myInFolder
	stats.CurrentFolderFiles = len(files)
	stats.CurrentFolderProgress = 0
	WalkFiles(files, oldDir.Imgs, myInFolder, myOutFolder, &myDir)
	newDir := myDir.MakeOld()
	for k := range oldDir.Subs {
		if strings.ReplaceAll(k, "/", "") == "" || strings.Contains(k, "/") {
			return myDir, fmt.Errorf("invalid Subs entry found in old index.json in folder \"%v\"", myOutFolder)
		}
		if _, ok := newDir.Subs[k]; !ok {
			folder := myOutFolder + "/" + k
			err := os.RemoveAll(folder)
			if err != nil {
				printStats(fmt.Errorf("deleting folder \"%v\" failed: %w", folder, err))
			}
		}
	}
	for k := range oldDir.Imgs {
		if strings.ReplaceAll(k, "/", "") == "" || strings.Contains(k, "/") {
			return myDir, fmt.Errorf("invalid Imgs entry found in old index.json in folder \"%v\"", myOutFolder)
		}
		if _, ok := newDir.Imgs[k]; !ok {
			file := myOutFolder + "/" + k
			err := os.RemoveAll(file + ".h.webp")
			if err != nil {
				printStats(fmt.Errorf("deleting file \"%v.h.webp\" failed: %w", file, err))
			}
			err = os.RemoveAll(file + ".s.webp")
			if err != nil {
				printStats(fmt.Errorf("deleting file \"%v.s.webp\" failed: %w", file, err))
			}
		}
	}
	myDir.SortByModifiedDesc()
	txt, _ := json.Marshal(myDir)
	err = ioutil.WriteFile(myOutFolder+"/index.json", txt, os.ModePerm)
	return myDir, err
}

func WalkFiles(files []os.FileInfo, oldImgs map[string]Img, myInFolder string, myOutFolder string, myDir *Dir) {
	var wg sync.WaitGroup
	for _, f := range files {
		name := norm.NFC.String(f.Name())
		if name[0] == '.' {
			stats.CurrentFolderProgress++
			continue
		}
		if !f.IsDir() {
			var img Img
			var ok bool
			if img, ok = oldImgs[name]; ok {
				ok, img.ModTime, _ = CheckFileAge(f, myOutFolder+"/"+name)
				if alwaysProcessRotatedImages {
					format := strings.ToLower(filepath.Ext(name))
					if format == ".cr2" {
						img, err := Info(myInFolder+"/"+name, f)
						if err == nil && img.Orientation > 1 {
							ok = false
						}
					}
				}
			}
			if ok {
				myDir.AddImage(img)
				stats.BytesSkipped += f.Size()
				stats.ImagesSkipped++
			} else {
				wg.Add(1)
				go func(f os.FileInfo, name string) {
					defer wg.Done()
					concurrency <- struct{}{}
					defer func() {
						<-concurrency
						stats.CurrentFolderProgress++
					}()
					stats.BytesRead += f.Size()
					img, err := Run(myInFolder+"/"+name, myOutFolder+"/"+name, f)
					if err != nil || (img == Img{}) {
						if err != nil {
							stats.LastErrorName = myInFolder + "/" + name
							printStats(fmt.Errorf("processing \"%v\" failed: %w", myInFolder+"/"+name, err))
						}
						myDir.AddMisc(name)
						stats.ImagesIgnored++
						myDir.Files++
					} else {
						myDir.AddImage(img)
						stats.ImagesProcessed++
					}
				}(f, name)
			}
		} else {
			stats.CurrentFolderProgress++
		}
	}
	wg.Wait()
}

func WalkSubdirectories(files []os.FileInfo, folder string, oldSubs map[string]MinDir, myOutFolder string, myDir *Dir, r repository.Repository) error {
	for _, f := range files {
		if f.IsDir() {
			if f.Name()[0] == '.' {
				continue
			}
			name := norm.NFC.String(f.Name())
			if _, ok := exclusions[StripLeadingSlash(folder+"/"+name)]; ok {
				stats.FoldersIgnored++
				continue
			}
			_, ok := oldSubs[name]
			var maxage bool
			if ok {
				maxage, _ = CheckFolderAge(myOutFolder + "/" + name + "/" + "index.json")
			}
			for _, path := range invalidatePaths {
				if strings.HasPrefix(path, StripLeadingSlash(folder+"/"+name)) {
					maxage = false
				}
			}
			if !ok || !maxage {
				subDir, err := Walk(name, folder+"/"+name, r)
				if err != nil {
					printStats(fmt.Errorf("walking through \"%v\" failed: %w", StripLeadingSlash(folder+"/"+name), err))
					if ok {
						return err
					}
				}
				myDir.AddFolder(subDir.MinDir)
				stats.FoldersProcessed++
			} else {
				dir := oldSubs[name]
				myDir.AddFolder(dir)
				stats.FoldersSkipped++
			}
		}
	}
	return nil
}
