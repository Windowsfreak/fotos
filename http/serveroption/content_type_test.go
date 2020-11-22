package serveroption

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"fotos/http/contextkeys"
)

func TestPutContentTypeIntoContext(t *testing.T) {
	inputRequest := httptest.NewRequest("GET", "https://localhost/", nil)
	inputRequest.Header.Add("content-type", "application/important")
	ctx := ExtractContentTypeIntoContext(context.Background(), inputRequest)
	actualContentType := ctx.Value(contextkeys.ContentType).(string)
	expectedContentType := "application/important"

	assert.Equal(t, expectedContentType, actualContentType)
}
