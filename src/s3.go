package hugo

// Taken from https://blog.tocconsulting.fr/upload-files-and-directories-to-aws-s3-using-golang/

import (
	"os"
	"fmt"
	"path/filepath"
 
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func isDirectory(path string) bool {
	fd, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	switch mode := fd.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	}
	return false
}
 
func uploadDirToS3(sess *session.Session, bucketName string, bucketPrefix string, dirPath string) (string, error) {
	fileList := []string{}
	filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
		fmt.Println("PATH ==> " + path)
		if isDirectory(path) {
			// Do nothing
			return nil
		} else {
			fileList = append(fileList, path)
			return nil
		}
	})
 
	for _, file := range fileList {
		err := uploadFileToS3(sess, bucketName, bucketPrefix, file)
		if err != nil {
			return "", err
		}
	}
	return bucketName, nil
}

 
func uploadFileToS3(sess *session.Session, bucketName string, bucketPrefix string, filePath string) error {
	fmt.Println("upload " + filePath + " to S3")
	// An s3 service
	s3Svc := s3.New(sess)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open file", file, err)
		return err
	}
	defer file.Close()
	var key string
	fileDirectory, _ := filepath.Abs(filePath)
	key = bucketPrefix + fileDirectory
	// Upload the file to the s3 given bucket
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName), // Required
		Key:    aws.String(key),        // Required
		Body:   file,
	}
	_, err = s3Svc.PutObject(params)
	if err != nil {
		fmt.Printf("Failed to upload data to %s/%s, %s\n",
			bucketName, key, err.Error())
		return err
	}
	return nil
}
 
func makeSession(profile string) *session.Session {
	// Enable loading shared config file
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	// Specify profile to load for the session's config
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profile,
	})
	if err != nil {
		fmt.Println("failed to create session,", err)
		fmt.Println(err)
		os.Exit(1)
	}
 
	return sess
}