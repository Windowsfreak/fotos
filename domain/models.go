package domain

import "time"

type Error struct {
	ErrorMessage string `json:"errorMessage"`
}

type AddPictureRequest struct {
	UserId        string `json:"userId"`
	UserName      string `json:"userName"`
	Discriminator string `json:"discriminator"`
	Gallery       string `json:"gallery"`
	Url           string `json:"url"`
	PreSharedKey  string `json:"preSharedKey"`
}

type DeletePictureRequest struct {
	UserId        string `json:"userId"`
	UserName      string `json:"userName"`
	Discriminator string `json:"discriminator"`
	Gallery       string `json:"gallery"`
	Filename      string `json:"filename"`
	PreSharedKey  string `json:"preSharedKey"`
}

type PictureResponse struct {
	UserId        string `json:"userId"`
	UserName      string `json:"userName"`
	Discriminator string `json:"discriminator"`
	Gallery       string `json:"gallery"`
	Filename      string `json:"filename"`
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
	InvalidatePaths            []string      `yaml:"InvalidatePaths"`
}
