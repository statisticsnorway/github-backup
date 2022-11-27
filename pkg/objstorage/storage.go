package objstorage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var client *storage.Client

func init() {
	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	check(err)
	client = c
}

func CopyToBucket(localSrcFile *os.File, bucketName string) {
	fmt.Printf("copying %s to bucket\n", localSrcFile.Name())
	ctx := context.Background()
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(localSrcFile.Name())
	bucketWriter := obj.NewWriter(ctx)
	defer bucketWriter.Close()
	buf := make([]byte, 51200)
	for {
		bytesRead, err := localSrcFile.Read(buf)
		if err == io.EOF {
			break
		}
		check(err)
		if bytesRead > 0 {
			_, err := bucketWriter.Write(buf[:bytesRead])
			check(err)
		}
	}
	fmt.Printf("copied %s to bucket\n", localSrcFile.Name())
}

func FilenameFor(reponame string) string {
	now := time.Now().Format("2006-01-02T15:04:05-0700")
	withSlashes := fmt.Sprintf("ghbackup_%s_%s.tar.gz", reponame, now)
	return strings.ReplaceAll(withSlashes, "/", "_")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
