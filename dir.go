package main

import (
	"sort"
	"time"
)

type Dir struct {
	MinDir
	Path string   `json:"-"`
	Subs []MinDir `json:"subs,omitempty"`
	Imgs []Img    `json:"imgs,omitempty"`
	Misc []string `json:"misc,omitempty"`
}

func (d *Dir) Sort() {
	sort.SliceStable(d.Subs, func(i, j int) bool {
		return cl.CompareString(d.Subs[i].N, d.Subs[j].N) < 0
	})
	sort.SliceStable(d.Imgs, func(i, j int) bool {
		return cl.CompareString(d.Imgs[i].N, d.Imgs[j].N) < 0
	})
}
func (d *Dir) AddFolder(subDir MinDir) {
	d.Subs = append(d.Subs, subDir)
	d.Files += subDir.Files
	d.Folders += subDir.Folders + 1
	d.TotalImages += subDir.TotalImages
	if d.Oldest.After(subDir.Oldest) {
		d.Oldest = subDir.Oldest
	}
	if d.Newest.Before(subDir.Newest) {
		d.Newest = subDir.Newest
	}
}
func (d *Dir) AddImage(img Img) {
	d.Imgs = append(d.Imgs, img)
	d.Images++
	d.TotalImages++
	if (img.D != time.Time{} && d.Oldest.After(img.D)) {
		d.Oldest = img.D
	}
	if d.Newest.Before(img.D) {
		d.Newest = img.D
	}
}
func (d *Dir) AddMisc(name string) {
	d.Misc = append(d.Misc, name)
}

func (d Dir) MakeOld() OldDir {
	out := OldDir{
		Subs: map[string]MinDir{},
		Imgs: map[string]Img{},
	}
	for _, d := range d.Subs {
		out.Subs[d.N] = d
	}
	for _, i := range d.Imgs {
		out.Imgs[i.N] = i
	}
	return out
}

type MinDir struct {
	N           string    `json:"n"`
	Files       int       `json:"files,omitempty"`
	Folders     int       `json:"folders,omitempty"`
	Images      int       `json:"images,omitempty"`
	TotalImages int       `json:"totalImages,omitempty"`
	Oldest      time.Time `json:"oldest,omitempty"`
	Newest      time.Time `json:"newest,omitempty"`
	C           string    `json:"c"`
}

type OldDir struct {
	Subs map[string]MinDir
	Imgs map[string]Img
}
