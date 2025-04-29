// port for uploading images from triple-s
package port

import (
	"io"
)

type ImageUploader interface {
	CreateBucket(bucketName string) error
	UploadImage(postID, filename string, r io.Reader) (string, error)
}
