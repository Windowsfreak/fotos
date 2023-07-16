package fotos

import (
	"encoding/json"
	"errors"
	"fmt"
	"fotos/domain"
	"golang.org/x/text/unicode/norm"
	"gopkg.in/yaml.v2"
	"io/fs"
	"math"
	"os"
	"path"
	"strings"
	"time"
)

func Walk(inPrefix string, inFolder string, outPrefix string, nonce string) <-chan MinDir {
	result := make(chan MinDir)

	inPath := inPrefix + "/" + inFolder
	outFolder := splitHash(hashFile(inFolder, nonce))
	outJson := outFolder + "/" + hashFile("index.json", nonce)

	files, err := os.ReadDir(inPath)
	if err != nil {
		println(fmt.Errorf("could not list contents of folder \"%v\": %w", inPath, err).Error())
		return nil
	}

	oldDir := Dir{}
	data, err := os.ReadFile(outPrefix + "/" + outJson + ".json")
	if errors.Is(err, fs.ErrNotExist) {
		if err := os.MkdirAll(outPrefix+"/"+outFolder, os.ModePerm); err != nil {
			println(fmt.Errorf("could not create folder structure for \"%v\": %w", outPrefix+"/"+outFolder, err).Error())
			return nil
		}
		if file, err := os.Create(outPrefix + "/" + outFolder + "/index.htm"); err != nil {
			println(fmt.Errorf("could not create stub index.htm for \"%v\": %w", outPrefix+"/"+outFolder, err).Error())
		} else {
			err := file.Close()
			if err != nil {
				println(fmt.Errorf("could not close index.htm for \"%v\": %w", outPrefix+"/"+outFolder, err).Error())
			}
		}
	} else if err != nil {
		println(fmt.Errorf("file \"%v\" could not be read: %w", outPrefix+"/"+outJson, err).Error())
		return nil
	} else {
		err = json.Unmarshal(data, &oldDir)
		if err != nil {
			println(fmt.Errorf("contents of file \"%v\" could not be unmarshalled: %w", outPrefix+"/"+outJson, err).Error())
			return nil
		}
	}

	if oldDir.ModTime > domain.Config.LastCompletion {
		println("skipping " + inFolder)
		go func() {
			result <- oldDir.MinDir
		}()
		return result
	}

	removedFolders := GetNonexistentFolders(oldDir.Subs, files)
	for _, folder := range removedFolders {
		WalkClean(inFolder+"/"+folder.Name, outPrefix, folder.Nonce)
	}
	removedImages := getNonexistentFiles(oldDir.Imgs, files)
	for _, file := range removedImages {
		RemoveThumbnails(outPrefix + "/" + outFolder + "/" + file.Path)
	}

	PrepareSubsNonces(&oldDir, files)

	// Save backups of subdirectory Nonces
	txt, err := json.Marshal(oldDir)
	if err != nil {
		println(fmt.Errorf("could not save file \"%v\": %w", outPrefix+"/"+outJson, err).Error())
	}
	err = os.WriteFile(outPrefix+"/"+outJson+".json", txt, os.ModePerm)
	if err != nil {
		println(fmt.Errorf("could not save file \"%v\": %w", outPrefix+"/"+outJson, err).Error())
	}

	minDirsPromise := WalkSubdirectories(files, inPrefix, inFolder, outPrefix, oldDir.Subs)
	walkFilesPromise := WalkFiles(files, oldDir.Imgs, inPrefix, inFolder, outPrefix, outFolder, nonce)

	go func() {
		minDirs := <-minDirsPromise
		walkFilesResult := <-walkFilesPromise

		dir := Dir{
			MinDir: MinDir{
				Path:    outJson,
				Name:    inFolder,
				Nonce:   nonce,
				ModTime: time.Now().Unix(),
			},
			Subs: minDirs,
			Imgs: walkFilesResult.Imgs,
			Misc: walkFilesResult.MiscFiles,
		}
		dir.NumDirs, dir.NumImgs, dir.NumMisc = AggregateTotals(dir)
		dir.Image, dir.Color, dir.Recent, dir.Oldest = AggregateImages(dir, outFolder)

		txt, err := json.Marshal(dir)
		if err != nil {
			println(fmt.Errorf("could not save file \"%v\": %w", outPrefix+"/"+outJson, err).Error())
		}
		err = os.WriteFile(outPrefix+"/"+outJson+".json", txt, os.ModePerm)
		if err != nil {
			println(fmt.Errorf("could not save file \"%v\": %w", outPrefix+"/"+outJson, err).Error())
		}
		result <- dir.MinDir
	}()

	return result
}

func WalkClean(inFolder string, outPrefix string, nonce string) {
	outFolder := splitHash(hashFile(inFolder, nonce))
	outJson := outFolder + "/" + hashFile("index.json", nonce) + ".json"

	dir := Dir{}
	data, err := os.ReadFile(outPrefix + "/" + outJson)
	if err != nil {
		println(fmt.Errorf("file \"%v\" could not be read: %w", outPrefix+"/"+outJson, err).Error())
		println("i.e. not all files and folders could be removed")
		return
	} else {
		err = json.Unmarshal(data, &dir)
		if err != nil {
			println(fmt.Errorf("contents of file \"%v\" could not be unmarshalled: %w", outPrefix+"/"+outJson, err).Error())
			println("i.e. not all files and folders could be removed")
			return
		}
	}
	for _, minDir := range dir.Subs {
		WalkClean(inFolder+"/"+minDir.Name, outPrefix, minDir.Nonce)
	}
	// println("Here I would delete the folder " + outFolder + " with all contents")
	if err := os.RemoveAll(outPrefix + "/" + outFolder); err != nil {
		println(fmt.Errorf("deleting folder \"%v\" unsuccessful: %w", outPrefix+"/"+outFolder, err).Error())
		println("i.e. not all files and folders could be removed")
	}
}

func PrepareSubsNonces(oldDir *Dir, files []os.DirEntry) {
	for _, dirEntry := range files {
		name := norm.NFC.String(dirEntry.Name())
		if strings.HasPrefix(name, ".") || !dirEntry.IsDir() {
			continue
		}
		_, ok := FindDir(oldDir.Subs, name)
		if !ok {
			sub := MinDir{
				Name:    name,
				Nonce:   MakeNonce(),
				ModTime: 0,
			}
			oldDir.Subs = append(oldDir.Subs, sub)
		}
	}
}

func WalkSubdirectories(files []os.DirEntry, inPrefix string, inFolder string, outPrefix string, subs []MinDir) <-chan []MinDir {
	var subdirPromises []<-chan MinDir

	for _, dirEntry := range files {
		name := norm.NFC.String(dirEntry.Name())
		if strings.HasPrefix(name, ".") || !dirEntry.IsDir() {
			continue
		}

		oldDir, ok := FindDir(subs, name)
		nonce := oldDir.Nonce
		if !ok {
			println(fmt.Errorf("subdirectory \"%v\" was not prepared", inFolder+"/"+name).Error())
			println("i.e. not all subdirectories could be created")
			continue
			// nonce = MakeNonce()
		}

		subdirPromise := Walk(inPrefix, inFolder+"/"+name, outPrefix, nonce)
		subdirPromises = append(subdirPromises, subdirPromise)
	}

	resultPromise := make(chan []MinDir)
	go func() {
		var subs []MinDir
		for _, promise := range subdirPromises {
			sub := <-promise
			sub.Name = path.Base(sub.Name)
			subs = append(subs, sub)
		}
		SortMinDirs(subs)
		resultPromise <- subs
	}()
	return resultPromise
}

func GetNonexistentFolders(minDirs []MinDir, dirEntries []fs.DirEntry) []MinDir {
	existingFolders := make(map[string]bool)

	// Build a map of existing folder names
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			name := norm.NFC.String(dirEntry.Name())
			existingFolders[name] = true
		}
	}

	var nonexistentFolders []MinDir
	for _, minDir := range minDirs {
		if !existingFolders[minDir.Name] {
			nonexistentFolders = append(nonexistentFolders, minDir)
		}
	}

	return nonexistentFolders
}

func getNonexistentFiles(imgs []Img, dirEntries []fs.DirEntry) []Img {
	existingFiles := make(map[string]bool)

	// Build a map of existing filenames
	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			name := norm.NFC.String(dirEntry.Name())
			existingFiles[name] = true
		}
	}

	var nonexistentFiles []Img
	for _, img := range imgs {
		if !existingFiles[img.Name] {
			nonexistentFiles = append(nonexistentFiles, img)
		}
	}

	return nonexistentFiles
}

func AggregateTotals(dir Dir) (folders int, images int, files int) {
	folders = len(dir.Subs)
	images = len(dir.Imgs)
	files = len(dir.Misc)
	for _, subDir := range dir.Subs {
		folders += subDir.NumDirs
		images += subDir.NumImgs
		files += subDir.NumMisc
	}
	return
}

func AggregateImages(dir Dir, localPath string) (string, string, int64, int64) {
	var image, color string
	oldest := int64(math.MaxInt64)
	recent := int64(math.MinInt64)

	for _, img := range dir.Imgs {
		if img.Date > recent {
			recent = img.Date
			image = localPath + "/" + img.Path
			color = img.Color
		}
		if img.Date < oldest {
			oldest = img.Date
		}
	}

	for _, subDir := range dir.Subs {
		if subDir.Recent > recent {
			recent = subDir.Recent
			image = subDir.Image
			color = subDir.Color
		}
		if subDir.Oldest < oldest {
			oldest = subDir.Oldest
		}
	}

	return image, color, recent, oldest
}

func DirMainTest() error {
	if err := ReadConfig(); err != nil {
		return err
	}

	/*server := NewServer()
	println(server.ListenAndServe().Error())

	proof, img, err := Validate2(Proof2)
	println(proof)
	println(img)
	println(err.Error())
	os.Exit(0)*/

	inPrefix := domain.Config.InPrefix
	outPrefix := domain.Config.OutPrefix
	inFolder := domain.Config.InFolder
	resultPromise := Walk(inPrefix, inFolder, outPrefix, domain.Config.Nonce)
	value := <-resultPromise

	domain.Config.LastCompletion = time.Now().Unix()
	domain.Config.OutFolder = value.Path
	return WriteConfig()
}

func ReadConfig() error {
	f, err := os.Open("config.yml")
	if err != nil {
		println(fmt.Errorf("could not load config file: %w", err).Error())
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&domain.Config)
	if err != nil {
		println(fmt.Errorf("could not unmarshal config file: %w", err).Error())
		return err
	}

	if len(domain.Config.SecretKey) < 1 {
		err = fmt.Errorf("secret key missing in config")
		println(err.Error())
		return err
	}
	return nil
}

func WriteConfig() error {
	data, err := yaml.Marshal(domain.Config)
	if err != nil {
		println(fmt.Errorf("could not marshal config file: %w", err).Error())
		return err
	}

	err = os.WriteFile("config.yml", data, 0644)
	if err != nil {
		println(fmt.Errorf("could not save config file: %w", err).Error())
		return err
	}
	return nil
}
