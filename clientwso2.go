package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	minio "github.com/minio/minio-go"
	"github.com/minio/minio/pkg/auth"
)

func main() {
	//Returns Access Token
	accessToken, err := getAccessToken()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("token=", accessToken)

	//Exchange Access Token from IDP for Minio Credentials
	cred, err := getMinioCred(accessToken) //minioCred.go
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("AccessKey: %s SecretKey: %s ExpirationTime: %f\n", cred.AccessKey, cred.SecretKey, cred.ExpTime)

	endpoint := "127.0.0.1:9000"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, cred.AccessKey, cred.SecretKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called mymusic.
	bucketName := "mymusic"
	location := "demo"

	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := minioClient.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	}
	log.Printf("Successfully created %s\n", bucketName)

	// Upload the zip file
	objectName := "code2.java"
	filePath := "/Users/sanatmouli/Downloads/code2.java"
	contentType := "application.txt"

	// Upload the zip file with FPutObject
	n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
}

func getMinioCred(accessToken string) (*auth.Credentials, error) {
	minioTokenUrl := "http://localhost:4000"
	resource := "/getminiotoken"
	u, err := url.ParseRequestURI(minioTokenUrl)
	if err != nil {
		//return nil, err
		log.Fatalln(err)
	}
	u.Path = resource
	urlStr := u.String()
	data := url.Values{}
	data.Add("AccessToken", accessToken)

	client := &http.Client{}
	r, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		//return nil, err
		log.Fatalln(err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(r)
	if err != nil {
		//return nil, err
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//return nil, err
		log.Fatalln(err)
	}

	cred := &auth.Credentials{}
	json.Unmarshal(body, cred)
	defer resp.Body.Close()
	return cred, nil
}
