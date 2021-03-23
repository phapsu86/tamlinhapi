package controllers

import (
	"fmt"
	//"github.com/gorilla/mux"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/phapsu86/tamlinh/api/auth"
	"github.com/phapsu86/tamlinh/api/responses"
	"github.com/phapsu86/tamlinh/api/utils/formaterror"

	"context"

	"github.com/phapsu86/tamlinh/api/models"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
// 	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")

// }

func (server *Server) UploadFile(w http.ResponseWriter, r *http.Request) {

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, err)
		return
	}
	str := fmt.Sprint(uid)
	data := []byte("hello" + str)
	fmt.Printf("ooooolalaa %x", md5.Sum(data))
	hashAvatar := md5.Sum(data)
	file, handler, err := r.FormFile("file")
	fileName := hex.EncodeToString(hashAvatar[:]) + "_" + r.FormValue("file_name")
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
	
		panic(err)
		//errs := formaterror.ReturnErr(err)

	}
	defer file.Close()

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		
		panic(err)

	}

	defer f.Close()
	//  _, _ = io.WriteString(w, "File "+fileName+" Uploaded successfully")
	_, _ = io.Copy(f, file)
	// Minio

	endpoint := "tamlinh.aibanhmi.com:9000"
	accessKeyID := "AKIAIOSFODNN7EXAMPLE"
	secretAccessKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now set up

	// Make a new bucket called mymusic.
	bucketName := "avatar"
	location := "us-east-1"
	ctx := context.Background()
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	contentType := "image/*"

	// Upload the zip file with FPutObject
	//n, err := minioClient.PutObject("my-bucketname", "my-objectname", reader, "application/octet-stream")
	n, err := minioClient.FPutObject(ctx, bucketName, fileName, handler.Filename, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", fileName, n)
	//update avatar to user
	mdUser := models.User{}
	userInfo, err := mdUser.FindUserByID(server.DB, uid)
	if err != nil {
		responses.JSON(w, http.StatusOK, formaterror.ReturnErr(err))
		return
	}
	userInfo.Avartar = fileName
	u, err := userInfo.UpdateAvartarUser(server.DB, uid)
	if err != nil {
		responses.JSON(w, http.StatusOK, formaterror.ReturnErr(err))
		return
	}
	//Get link file
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+fileName+"\"")
	// Generates a presigned url which expires in a day.
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), bucketName, fileName, time.Second*24*60*60, reqParams)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully generated presigned URL", presignedURL)

	///=======================================
	if u != nil {
		fmt.Println("Successfully %v", u)
	}
	result := formaterror.ReturnUploadSuccess()
	result.Url = presignedURL.String()
	result.Msg = "File " + fileName + " Uploaded successfully"
	responses.JSON(w, http.StatusOK, result)
}
