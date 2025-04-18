// port for uploading images from triple-s
package port

import (
	"context"
	"io"
)

type ImageUploader interface {
	Upload(ctx context.Context, file io.Reader, filename string) (string, error)
}
