// port for uploading images from triple-s
package port

import (
	"io"
)

type ImageUploader interface {
	CreateBucket(bucketName string) error
	UploadPostImage(postID, filename string, r io.Reader) (string, error)
	UploadCommentImage(postID, commentID, filename string, r io.Reader) (string, error)
}
