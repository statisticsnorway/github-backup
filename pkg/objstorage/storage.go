package objstorage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"os"
	"path/filepath"
	"time"
)

var client *storage.Client
var objBasePath = filepath.Join("/", time.Now().Format("2006/01/02"))

func init() {
	ctx := context.Background()
	c, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		panic(err)
	}
	client = c
}

func CopyToBucket(localSrcFile *os.File, bucketName string) error {
	fmt.Printf("copying '%s' to bucket '%s'\n", localSrcFile.Name(), bucketName)
	srcFilename, err := FilenameWithoutPath(localSrcFile)
	if err != nil {
		return err
	}
	ctx := context.Background()
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(filepath.Join(objBasePath, srcFilename))
	bucketWriter := obj.NewWriter(ctx)
	defer bucketWriter.Close()
	written, err := io.Copy(bucketWriter, localSrcFile)
	if err != nil {
		return err
	}
	fmt.Printf("wrote %d bytes to '%s'\n", written, bucketName)
	return nil
}

func FilenameWithoutPath(f *os.File) (string, error) {
	info, err := f.Stat()
	if err != nil {
		return "", err
	}
	return info.Name(), nil
}
