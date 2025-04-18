package imageuploader

import (
	"fmt"
	"os"
	"path/filepath"
)

// if my post will be deleted (archived), should I change anything with the folders?
// in theory, even if I delete a post, and its bucket is still alive, I should not be able to create another folder with the same name (uuid)

// if a user uploads images with the same name to the same buckets, how should I handle that?

// Need to store URL paths in DB!! (via service)

// Creates a folder for each post (bucketName == PostID)
func BucketCreate(bucketName, rootdir string) error {

	testFile := filepath.Join(rootdir, ".testfile")
	file, err := os.Create(testFile)
	if err != nil {
		// log this
		return fmt.Errorf("problem with writing in the directory: %w", err)
	}

	file.Close()
	os.Remove(testFile)

	bucketPath := fmt.Sprintf("%s/%s", rootdir, bucketName)
	if _, err := os.Stat(bucketPath); !os.IsNotExist(err) {
		// log this
		return fmt.Errorf("bucket already exists")
	}

	err = os.Mkdir(bucketPath, 0o777)
	if err != nil {
		// log this
		return fmt.Errorf("failed to create bucket: %s", bucketName)
	}

	return nil
}
