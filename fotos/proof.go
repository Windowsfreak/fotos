package fotos

import (
	"encoding/json"
	"errors"
	"fmt"
	"fotos/domain"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

func Validate2(proofLink string) (string, string, error) {
	const originalPath = "/fotos/images/"
	const proofSeparator = "?proof="
	sections := strings.Split(proofLink, proofSeparator)
	if len(sections) != 2 {
		return "", "", fmt.Errorf("no proof provided")
	}
	needle := strings.TrimPrefix(sections[0], originalPath)
	if needle == sections[0] {
		return "", "", fmt.Errorf("hostname mismatch")
	}
	slash := strings.LastIndex(needle, "/")
	outFolder := splitHash(hashFile(needle[:slash], sections[1][:8]))
	outJson := outFolder + "/" + hashFile("index.json", sections[1][:8])
	data, err := os.ReadFile(domain.Config.OutPrefix + "/" + outJson + ".json")
	var dir Dir
	if errors.Is(err, fs.ErrNotExist) {
		return "", "", fmt.Errorf("proof does not exist")
	} else if err != nil {
		return "", "", fmt.Errorf("proof inaccessible")
	} else {
		err = json.Unmarshal(data, &dir)
		if err != nil {
			return "", "", fmt.Errorf("proof invalid")
		}
	}
	for _, img := range dir.Imgs {
		if dir.Name+"/"+img.Name == needle && sections[1][8:] == img.Path {
			return dir.Name + "/" + img.Name, img.Name, nil
		}
	}
	return "", "", fmt.Errorf("proof mismatch")
}

func NewServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	addr := ":" + port
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	if path, img, err := Validate2(r.URL.Path + "?" + r.URL.RawQuery); err == nil {
		data, err := os.ReadFile(domain.Config.InPrefix + "/" + path)
		if err == nil {
			w.Header().Set("Content-Type", "application/x-download")
			w.Header().Set("Content-Disposition", "attachment; filename="+img)
			_, err := w.Write(data)
			if err != nil {
				println("Error writing response: " + err.Error())
			}
		} else {
			if _, err := fmt.Fprintf(w, "%s", "file not found"); err != nil {
				println("Error writing response: " + err.Error())
			}
			println(err.Error())
		}
	} else {
		if _, err := fmt.Fprintf(w, "%s", err.Error()); err != nil {
			println("Error writing response: " + err.Error())
		}
	}
}
