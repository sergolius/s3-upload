package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	bucket = os.Getenv("S3_BUCKET")
	region = os.Getenv("S3_REGION")
	ID     = os.Getenv("S3_ID")
	SECRET = os.Getenv("S3_SECRET")
	TOKEN  = os.Getenv("S3_TOKEN")
	LOG    = os.Getenv("S3_LOG")
)

var LogLevel = map[string]int{
	"warn":  0,
	"debug": 1,
	"trace": 2,
	"info":  2,
}

// validateEnv
func validateEnv() {
	for _, r := range []struct {
		name, value string
	}{
		{"S3_BUCKET", bucket},
		{"S3_REGION", region},
		{"S3_ID", ID},
		{"S3_SECRET", SECRET},
	} {
		if r.value == "" {
			log.Fatalf("%q is required \n", r.name)
		}
	}
}

func main() {
	start := time.Now()

	validateEnv()

	flag.Parse()

	logLevel := LogLevel[LOG]
	folders := flag.Args()
	if 0 == len(folders) {
		log.Fatal("folders to upload is not defined")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(ID, SECRET, TOKEN),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Uploading.")

	filePaths := make(chan string)

	go func() {
		for _, folder := range folders {
			if err := fileList(folder, filePaths); err != nil {
				log.Printf("%s: %v \n\n", folder, err)
			}
		}
		close(filePaths)
	}()

	for filePath := range filePaths {
		if strings.HasPrefix(filepath.Base(filePath), ".") {
			continue
		}

		if resp, err := putFile(sess, bucket, filePath); err != nil {
			log.Printf("upload failed: %q, err: %v \n", filePath, err)
		} else {
			if 1 <= logLevel {
				log.Printf("%q - Uploaded \n", filePath)
			}
			if 2 <= logLevel {
				log.Printf("Result: %v \n", resp)
			}
		}
	}
	end := time.Now()

	log.Printf("Done. %s \n", end.Sub(start))
}

// fileList walks the file tree, seeks for only files
func fileList(folder string, fileChannel chan string) error {
	return filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		fileChannel <- path
		return nil
	})
}

// putFile sends file to S3
func putFile(sess *session.Session, bucket, filePath string) (*s3.PutObjectOutput, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
		Body:   file,
	}
	return s3.New(sess).PutObject(params)
}
