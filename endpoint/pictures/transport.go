package pictures

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kitlogrus "github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/sirupsen/logrus"

	"fotos/domain"
	mhttp "fotos/http"
	"fotos/http/middleware"
	"fotos/http/serveroption"
)

// MakeAddPictureHandler returns a handler for the Pictures service.
func MakeAddPictureHandler(
	s Service,
	logger *logrus.Entry,
) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(serveroption.ExtractContentTypeIntoContext),
		kithttp.ServerBefore(serveroption.ExtractAcceptHeaderIntoContext),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(kitlogrus.NewLogrusLogger(logger))),
		kithttp.ServerErrorEncoder(middleware.MakeEncodeErrorFunc(logger)),
	}

	mw := endpoint.Chain(
		middleware.MakeAcceptHeaderValidationMiddleware(),
		makeAddPictureValidationMiddleware(),
	)

	endpointHandler := kithttp.NewServer(
		mw(makeAddPictureEndpoint(s)),
		decodeAddPictureRequest,
		mhttp.EncodeResponse,
		opts...,
	)

	return endpointHandler
}

// MakeDeletePictureHandler returns a handler for the Pictures service.
func MakeDeletePictureHandler(
	s Service,
	logger *logrus.Entry,
) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(serveroption.ExtractContentTypeIntoContext),
		kithttp.ServerBefore(serveroption.ExtractAcceptHeaderIntoContext),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(kitlogrus.NewLogrusLogger(logger))),
		kithttp.ServerErrorEncoder(middleware.MakeEncodeErrorFunc(logger)),
	}

	mw := endpoint.Chain(
		middleware.MakeAcceptHeaderValidationMiddleware(),
		makeDeletePictureValidationMiddleware(),
	)

	endpointHandler := kithttp.NewServer(
		mw(makeDeletePictureEndpoint(s)),
		decodeDeletePictureRequest,
		mhttp.EncodeResponse,
		opts...,
	)

	return endpointHandler
}

// MakeGetRandomPictureHandler returns a handler for the Pictures service.
func MakeGetRandomPictureHandler(
	s Service,
	logger *logrus.Entry,
) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(serveroption.ExtractAcceptHeaderIntoContext),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(kitlogrus.NewLogrusLogger(logger))),
		kithttp.ServerErrorEncoder(middleware.MakeEncodeErrorFunc(logger)),
	}

	mw := endpoint.Chain(
		middleware.MakeAcceptHeaderValidationMiddleware(),
	)

	endpointHandler := kithttp.NewServer(
		mw(makeGetRandomPictureEndpoint(s)),
		nopFunction,
		mhttp.EncodeResponse,
		opts...,
	)

	return endpointHandler
}

func decodeAddPictureRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request domain.AddPictureRequest
	err := mhttp.DecodeRequest(ctx, r, &request)
	if err != nil {
		query := r.URL.Query()
		request.Url = query.Get("url")
		request.Gallery = query.Get("id")
		request.UserId = query.Get("id")
		request.UserName = query.Get("username")
		request.Discriminator = query.Get("discriminator")
		request.PreSharedKey = query.Get("token")
		err = nil
		println(r.URL.RawQuery)
	}
	return request, err
}

func decodeDeletePictureRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request domain.DeletePictureRequest
	err := mhttp.DecodeRequest(ctx, r, &request)
	if err != nil {
		query := r.URL.Query()
		request.Filename = query.Get("filename")
		request.Gallery = query.Get("id")
		request.UserId = query.Get("id")
		request.UserName = query.Get("username")
		request.Discriminator = query.Get("discriminator")
		request.PreSharedKey = query.Get("token")
		err = nil
		println(r.URL.RawQuery)
	}
	return request, err
}

func nopFunction(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}
