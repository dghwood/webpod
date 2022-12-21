package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"log"
)

const bucketName = "bsnek-316609.appspot.com"

func Store(fileBytes []byte, fileName string) string {
	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(fileName)
	writer := obj.NewWriter(ctx)
	_, writerErr := writer.Write(fileBytes)
	if writerErr != nil {
		log.Fatal(writerErr)
	}
	writer.Close()
	// Looks like this needs to be done after writing the file
	acl := obj.ACL()
	// TODO: check errors
	acl.Set(ctx, storage.AllUsers, storage.RoleReader)

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, fileName)
}
