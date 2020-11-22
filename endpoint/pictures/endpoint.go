package pictures

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"fotos/domain"
)

func makeAddPictureEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(domain.AddPictureRequest)
		req.PreSharedKey = ""
		err := s.StoreAddPictureUserInformation(req)
		if err != nil {
			return req, err
		}
		filename, err := s.AddPicture(req.Gallery, req.Url)
		return domain.PictureResponse{
			UserId:        req.UserId,
			UserName:      req.UserName,
			Discriminator: req.Discriminator,
			Gallery:       req.Gallery,
			Filename:      filename,
		}, err
	}
}

func makeDeletePictureEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(domain.DeletePictureRequest)
		req.PreSharedKey = ""
		err := s.StoreDeletePictureUserInformation(req)
		if err != nil {
			return req, err
		}
		err = s.DeletePicture(req.Gallery, req.Filename)
		return req, err
	}
}

func makeGetRandomPictureEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		randomPicture, err := s.GetRandomPicture()
		return randomPicture, err
	}
}
