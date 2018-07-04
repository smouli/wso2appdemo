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
)

func main() {
	//Returns Access Token
	accessToken, err := getAccessToken()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("token=", accessToken)

	//introspectEndpoint := "https://localhost:9443/oauth2/introspect"

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
	//minioClient, err := minio.New(endpoint, "EM5Q098B5FWARF4NV931", "/BJZ2WT0Mobv5wSphP/aGmzGiUrcT8yBXCIDtXNh", useSSL)
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
	// filePath := "/Users/sanatmouli/Downloads/code2.java"
	// contentType := "application.txt"

	// Upload the zip file with FPutObject
	// n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
	object, err := minioClient.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
    	fmt.Println(err)
    	return
	}
	_ , err = ioutil.ReadAll(object)
	

	// if err := minioClient.GetObject(bucketName, objectName, "code2.java", minio.GetObjectOptions{}); err != nil {
	// 	fmt.Println("HERE")
	// 	log.Fatalln(err)
	// }
	log.Printf("Successfully got  %s of size %d\n", objectName)
}

type credentials struct {
	AccessKey string `json:"accessKey,omitempty"`
	SecretKey string `json:"secretKey,omitempty"`
	ExpTime   float64
}

func getMinioCred(accessToken string) ( /*auth.Credentials*/ credentials, error) {
	minioTokenUrl := "http://localhost:9000/"
	resource := "/minio/admin/v1/sts"
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
		log.Fatalln(err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(r)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	cred := &credentials{}
	json.Unmarshal(body, cred)
	defer resp.Body.Close()
	return *cred, nil

}
