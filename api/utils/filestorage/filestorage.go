package filestorage

import (

"time"
"log"
"fmt"
"context"
"github.com/minio/minio-go/v7"
"github.com/minio/minio-go/v7/pkg/credentials"
)


func gelinkFile (file string) string {

var endpoint = "34.87.48.124:9000"
var accessKeyID = "AKIAIOSFODNN7EXAMPLE"
var secretAccessKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
var useSSL = false

    // Initialize minio client object.
    minioClient, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
        Secure: useSSL,
    })
    if err != nil {
        log.Fatalln(err)
    }

	log.Printf(" chay roi do %#v\n", minioClient) // minioClient is now setup
	
	found, err := minioClient.BucketExists(context.Background(), "tamlinh")
	if err != nil {
		fmt.Println("khong tim thay  %#v\n",err)
		return ""
	}
	if found {
		fmt.Println(" Tim thay roi Bucket found")
	}

	expiry := time.Second * 24 * 60 * 60 // 1 day.
	presignedURL, err := minioClient.PresignedPutObject(context.Background(), "tamlinh", "uploads/ll06W02y9i1aOzTXv086ZfloExHlnxGjaDqol4oD.png", expiry)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println("Successfully generated presigned URL", presignedURL)

	return "presignedURL"

}