package main

import (
	"fmt"
	"github.com/whatupdave/s3/s3util"
	"log"
	"net/http"
	"os"
)

var (
	s3url string
)

func main() {
	s3url = os.Getenv("S3_URL")
	s3util.DefaultConfig.AccessKey = os.Getenv("AWS_ACCESS_KEY")
	s3util.DefaultConfig.SecretKey = os.Getenv("AWS_SECRET_KEY")

	http.Handle("/", new(Router))

	fmt.Println("Listening")
	err := http.ListenAndServe(":7070", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
