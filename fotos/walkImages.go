package fotos

import (
	"fmt"
	"golang.org/x/text/unicode/norm"
	"io/fs"
	"os"
	"strings"
	"time"
)

var imgThreads = make(chan bool, 4)

type ImgOrMisc struct {
	Img  Img
	Misc string
}
type WalkFilesResult struct {
	Imgs          []Img
	MiscFiles     []string
	MostRecentImg Img
}

func WalkFiles(files []fs.DirEntry, imgs []Img, inPrefix string, inFolder string, outPrefix string, outFolder string, nonce string) <-chan WalkFilesResult {
	var imgPromises []<-chan ImgOrMisc
	inPath := inPrefix + "/" + inFolder
	outPath := outPrefix + "/" + outFolder

	for _, dirEntry := range files {
		if dirEntry.IsDir() {
			continue
		}
		name := norm.NFC.String(dirEntry.Name())
		if strings.HasPrefix(name, ".") {
			continue
		}
		file, err := dirEntry.Info()
		if err != nil {
			println(fmt.Errorf("could not read file info for \"%v\": %w", inPath+"/"+name, err).Error())
			continue
		}
		imgThreads <- true

		imgPromise := make(chan ImgOrMisc)
		imgPromises = append(imgPromises, imgPromise)

		go func(file fs.FileInfo) {
			var (
				img Img
				err error
			)
			p := hashFile(inFolder+"/"+name, nonce)

			if existingImg, ok := CheckImageAge(file, imgs, name, outPath, p); ok {
				img, err = existingImg, nil
			} else {
				img, err = ConvertWith(inPath+"/"+file.Name(), outPath+"/"+p, file)
			}
			<-imgThreads
			if err != nil {
				println(fmt.Errorf("could not convert image \"%v\": %w", inPath+"/"+name, err).Error())
				imgPromise <- ImgOrMisc{Img{}, name}
			} else if (img == Img{}) {
				imgPromise <- ImgOrMisc{Img{}, name}
			}
			img.Path = p
			imgPromise <- ImgOrMisc{img, ""}
		}(file)
	}

	resultPromise := make(chan WalkFilesResult)
	go func() {
		var images []Img
		var misc []string
		for _, promise := range imgPromises {
			resolved := <-promise
			if resolved.Misc == "" {
				images = append(images, resolved.Img)
			} else {
				misc = append(misc, resolved.Misc)
			}
		}
		SortImgs(images)
		SortStrings(misc)
		resultPromise <- WalkFilesResult{images, misc, Img{}}
	}()
	return resultPromise
}

func CheckImageAge(file fs.FileInfo, imgs []Img, name string, outPath string, p string) (Img, bool) {
	existingImg, ok := FindImg(imgs, name)
	if ok {
		if fileInfo, err := os.Stat(outPath + "/" + p + ".o.jxl"); err == nil {
			return existingImg, fileInfo.ModTime().After(file.ModTime().Add(5 * time.Second))
		}
	}
	return existingImg, false
}
func ImgMainTest() error {
	inPrefix := "/Users/bjoern/projects/bjoern/fotos/images"
	outPrefix := "/Users/bjoern/projects/bjoern/fotos/images/Neu"
	inFolder := "Original"
	outFolder := splitHash(hashFile("Original", "0"))
	inPath := inPrefix + "/" + inFolder
	outPath := outPrefix + "/" + outFolder
	files, err := os.ReadDir(inPath)
	if err != nil {
		return fmt.Errorf("could not list contents of folder \"%v\": %w", inPath, err)
	}
	if err := os.MkdirAll(outPath, os.ModePerm); err != nil {
		return fmt.Errorf("could not create folder structure for \"%v\": %w", outPath, err)
	}
	resultPromise := WalkFiles(files, []Img{}, inPrefix, inFolder, outPrefix, outFolder, "0")
	value := <-resultPromise
	println(len(value.Imgs))
	return nil
}
