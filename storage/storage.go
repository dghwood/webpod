package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
)

const bucketName = "webpods.appspot.com"

func Store(fileBytes []byte, fileName string) (url string, err error) {
	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return url, err
	}
	defer client.Close()

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(fileName)
	writer := obj.NewWriter(ctx)
	_, err = writer.Write(fileBytes)
	if err != nil {
		return url, err
	}
	err = writer.Close()
	if err != nil {
		return url, err
	}

	err = obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		return url, err
	}
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, fileName), nil
}
