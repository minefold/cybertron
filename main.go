package main

import (
	"fmt"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
	"net/http"
	"os"
)

var (
	bucket *s3.Bucket
)

func main() {
	auth, err := aws.EnvAuth()
	if err != nil || os.Getenv("S3_BUCKET") == "" {
		log.Fatal("set AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and S3_BUCKET")
	}

	s := s3.New(auth, aws.USEast)
	bucket = s.Bucket(os.Getenv("S3_BUCKET"))

	http.Handle("/", new(Router))

	fmt.Println("Listening")
	err = http.ListenAndServe(":7070", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
