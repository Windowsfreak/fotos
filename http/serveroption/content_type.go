package serveroption

import (
	"context"
	"net/http"

	"fotos/http/contextkeys"
	"fotos/http/header"
)

// ExtractContentTypeIntoContext extracts content type from an http request and injects it into the provided context.
func ExtractContentTypeIntoContext(ctx context.Context, r *http.Request) context.Context {
	ct := r.Header.Get(header.ContentType)
	ctx = context.WithValue(ctx, contextkeys.ContentType, ct)
	return ctx
}
