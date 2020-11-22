package pictures

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"fotos/domain"
	"fotos/fotos"
	"fotos/repository"
)

type Service interface {
	StoreAddPictureUserInformation(r domain.AddPictureRequest) error
	StoreDeletePictureUserInformation(r domain.DeletePictureRequest) error
	AddPicture(gallery string, source string) (string, error)
	DeletePicture(gallery string, filename string) error
	GetRandomPicture() (domain.PictureResponse, error)
}

type service struct {
	logger *logrus.Entry
	r      repository.Repository
}

func NewService(
	logger *logrus.Entry,
	repository repository.Repository,
) Service {
	return &service{
		logger: logger,
		r:      repository,
	}
}

func (s *service) StoreAddPictureUserInformation(r domain.AddPictureRequest) error {
	num, err := repository.ConvertStringToInt64(r.UserId)
	if err != nil {
		return err
	}
	return s.r.Store(num, r.UserName, r.Discriminator)
}

func (s *service) StoreDeletePictureUserInformation(r domain.DeletePictureRequest) error {
	num, err := repository.ConvertStringToInt64(r.UserId)
	if err != nil {
		return err
	}
	return s.r.Store(num, r.UserName, r.Discriminator)
}

func (s *service) AddPicture(gallery string, source string) (filename string, err error) {
	r, _ := http.NewRequest("GET", source, nil)
	filename = path.Base(r.URL.Path)
	err = ValidateFilename(filename)
	if err != nil {
		return
	}
	ext := path.Ext(filename)
	var fn string
	var fp string
	for ok := true; ok; ok = fileExists(fp) {
		fn = filenameWithoutExtension(filename)
		fn += "-" + randomString(5)
		fp = domain.Config.InFolder + gallery + "/" + fn + ext
	}
	filename = fn + ext
	err = downloadFile(fp, r)
	if err != nil {
		return
	}
	go func() {
		c := domain.Config
		c.InvalidatePaths = []string{gallery}
		fotos.TryRun(c, s.r)
	}()
	return
}

func (s *service) DeletePicture(gallery string, fn string) error {
	fp := domain.Config.InFolder + gallery + "/" + fn
	if !fileExists(fp) {
		return domain.ErrNotFound
	}
	err := os.Remove(fp)
	if err != nil {
		return err
	}
	fp = domain.Config.InFolder + gallery + "/"
	empty, err := isEmpty(fp)
	if err != nil {
		return err
	}

	if empty {
		err = os.Remove(fp)
	}
	if err != nil {
		return err
	}
	go func() {
		fotos.TryRun(domain.Config, s.r)
	}()
	return nil
}

func (s *service) GetRandomPicture() (domain.PictureResponse, error) {
	gallery, image, err := randomWalk(domain.Config.OutFolder)
	if err != nil {
		return domain.PictureResponse{
			Gallery:  gallery,
			Filename: image,
		}, err
	}
	segments := strings.Split(gallery, "/")
	num, err := repository.ConvertStringToInt64(segments[0])
	if err != nil {
		return domain.PictureResponse{
			Gallery:  gallery,
			Filename: image,
		}, nil
	}
	username, discriminator, err := s.r.Fetch(num)
	return domain.PictureResponse{
		UserId:        gallery,
		UserName:      username,
		Discriminator: discriminator,
		Gallery:       gallery,
		Filename:      image,
	}, err
}

func randomWalk(fp string) (string, string, error) {
	data, err := ioutil.ReadFile(fp + "/index.json")
	if err != nil {
		return "", "", fmt.Errorf("file \"index.json\" in \"%v\" could not be read: %w", fp, err)
	}
	p := fotos.Dir{}
	err = json.Unmarshal(data, &p)
	if p.TotalImages <= 0 {
		return "", "", fmt.Errorf("no image found")
	}
	pos := rand.Intn(p.TotalImages)
	if pos < p.Images {
		return "", p.Imgs[pos].N, nil
	}
	pos -= p.Images
	for _, sub := range p.Subs {
		if pos < sub.TotalImages {
			gallery, image, err := randomWalk(fp + "/" + sub.N)
			return sub.N + "/" + gallery, image, err
		}
		pos -= sub.TotalImages
	}
	return "", "", fmt.Errorf("no image found")
}

func downloadFile(fp string, r *http.Request) error {
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = os.MkdirAll(path.Dir(fp), os.ModePerm); err != nil {
		return fmt.Errorf("could not create folder structure for \"%v\": %w", fp, err)
	}

	out, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func fileExists(fn string) bool {
	info, err := os.Stat(fn)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func filenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

const charset = "abcdefghijklmnopqrstuvwxyz"

var random = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[random.Intn(len(charset))]
	}
	return string(b)
}
