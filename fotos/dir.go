package fotos

type Dir struct {
	MinDir
	Subs []MinDir `json:"subs,omitempty"`
	Imgs []Img    `json:"imgs,omitempty"`
	Misc []string `json:"misc,omitempty"`
}

type MinDir struct {
	Path    string `json:"p"`
	Name    string `json:"n"`
	Image   string `json:"i"`
	Nonce   string `json:"q"`
	NumDirs int    `json:"j,omitempty"`
	NumImgs int    `json:"k,omitempty"`
	NumMisc int    `json:"l,omitempty"`
	Oldest  int64  `json:"o,omitempty"`
	Recent  int64  `json:"r,omitempty"`
	Color   string `json:"c"`
	ModTime int64  `json:"m,omitempty"`
}

func FindDir(items []MinDir, name string) (MinDir, bool) {
	for i := range items {
		if items[i].Name == name {
			return items[i], true
		}
	}
	return MinDir{}, false
}
