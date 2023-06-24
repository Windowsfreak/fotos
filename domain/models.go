package domain

import (
	"io"
	"time"
)

type Error struct {
	ErrorMessage string `json:"errorMessage"`
}

type AddImageRequest struct {
	Folder       string `json:"folder"`
	Url          string `json:"url"`
	PreSharedKey string `json:"preSharedKey"`
}

type UploadImageRequest struct {
	File         io.Reader `json:"filename"`
	Folder       string    `json:"folder"`
	Filename     string    `json:"filename"`
	Url          string    `json:"url"`
	PreSharedKey string    `json:"preSharedKey"`
}

type DeleteImageRequest struct {
	Folder       string `json:"folder"`
	Filename     string `json:"filename"`
	PreSharedKey string `json:"preSharedKey"`
}

type GetImageRequest struct {
	Folder       string `json:"folder"`
	Filename     string `json:"filename"`
	PreSharedKey string `json:"preSharedKey"`
}

type GetThumbnailRequest struct {
	Folder       string `json:"folder"`
	Filename     string `json:"filename"`
	Size         int    `json:"size"`
	PreSharedKey string `json:"preSharedKey"`
}

type ImageResponse struct {
	Folder   string `json:"folder"`
	Filename string `json:"filename"`
}

type ConfigStruct struct {
	Threads                    int           `yaml:"Threads"`
	InFolder                   string        `yaml:"InFolder"`
	OutFolder                  string        `yaml:"OutFolder"`
	Path                       string        `yaml:"Path"`
	MaxAge                     time.Duration `yaml:"MaxAge"`
	Exclude                    string        `yaml:"Exclude"`
	AlwaysProcessRotatedImages bool          `yaml:"AlwaysProcessRotatedImages"`
	Plausibility               bool          `yaml:"Plausibility"`
	ServerAddr                 string        `yaml:"ServerAddr"`
	PreSharedKey               string        `yaml:"PreSharedKey"`
	RootFolderCaption          string        `yaml:"RootFolderCaption"`
	InvalidatePaths            []string      `yaml:"InvalidatePaths"`
}
