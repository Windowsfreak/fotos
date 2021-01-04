package pictures

import (
	"context"
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"

	"github.com/go-kit/kit/endpoint"

	"fotos/domain"
	"fotos/fotos"
)

var re = regexp.MustCompile("^[^\\/?%*:|\"<>\\.][^\\/?%*:|\"<>]*$")

func makeAddPictureValidationMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			msg, ok := req.(domain.AddPictureRequest)
			if !ok {
				return nil, domain.ErrInvalidMessageType
			}
			if err := validateAddPictureRequest(msg); err != nil {
				return nil, fmt.Errorf("%w: %s", domain.ErrMissingArgument, err)
			}
			return next(ctx, req)
		}
	}
}

func makeDeletePictureValidationMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			msg, ok := req.(domain.DeletePictureRequest)
			if !ok {
				return nil, domain.ErrInvalidMessageType
			}
			if err := validateDeletePictureRequest(msg); err != nil {
				return nil, fmt.Errorf("%w: %s", domain.ErrMissingArgument, err)
			}
			return next(ctx, req)
		}
	}
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func validateAddPictureRequest(msg domain.AddPictureRequest) error {
	if len(domain.Config.PreSharedKey) < 1 {
		log.Fatal("pre-shared key missing in config")
	}
	if msg.PreSharedKey != domain.Config.PreSharedKey {
		return domain.ErrInvalidKey
	}
	if !re.MatchString(msg.Gallery) || len(msg.Gallery) < 1 {
		return domain.ErrInvalidGalleryName
	}
	if !strings.HasPrefix(msg.Url, "http") {
		return domain.ErrUrlProtocol
	}
	if !containsString(fotos.MediaFiles, strings.ToLower(path.Ext(msg.Url))) {
		return domain.ErrUrlUnsupportedExt
	}
	if !re.MatchString(msg.Gallery) || len(msg.Gallery) < 1 {
		return domain.ErrInvalidGalleryName
	}
	return nil
}

func validateDeletePictureRequest(msg domain.DeletePictureRequest) error {
	if len(domain.Config.PreSharedKey) < 1 {
		log.Fatal("pre-shared key missing in config")
	}
	if msg.PreSharedKey != domain.Config.PreSharedKey {
		return domain.ErrInvalidKey
	}
	if !re.MatchString(msg.Gallery) || len(msg.Gallery) < 1 {
		return domain.ErrInvalidGalleryName
	}
	if !re.MatchString(msg.Filename) || len(msg.Filename) < 1 {
		return domain.ErrNotFound
	}
	return nil
}

func ValidateFilename(filename string) error {
	if !re.MatchString(filename) || len(filename) < 1 {
		return domain.ErrInvalidFileName
	}
	return nil
}
