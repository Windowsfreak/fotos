package fotos

import (
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"sort"
)

var cl = collate.New(language.German)

func SortMinDirs(subs []MinDir) {
	sort.Slice(subs, func(i, j int) bool {
		return cl.CompareString(subs[i].Name, subs[j].Name) < 0
	})
}

func SortImgs(imgs []Img) {
	sort.Slice(imgs, func(i, j int) bool {
		return imgs[i].Date < imgs[j].Date
	})
}

func SortStrings(misc []string) {
	sort.Slice(misc, func(i, j int) bool {
		return cl.CompareString(misc[i], misc[j]) < 0
	})
}
