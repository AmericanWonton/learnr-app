package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//Key variables for Amazon access
var AWSAccessKeyId string
var AWSSecretKey string
var bucketname string

/* This loads in our Amazon credentials. Called initially */
func loadAmazonCreds() {
	//Check to see if ENV Creds are available first
	_, ok := os.LookupEnv("LEARNR_BUCKET")
	if !ok {
		message := "This ENV Variable is not present: " + "LEARNR_BUCKET"
		panic(message)
	}
	_, ok2 := os.LookupEnv("AWS_ACCESS_KEY")
	if !ok2 {
		message := "This ENV Variable is not present: " + "AWS_ACCESS_KEY"
		panic(message)
	}
	_, ok3 := os.LookupEnv("AWS_SECRET_KEY")
	if !ok3 {
		message := "This ENV Variable is not present: " + "AWS_SECRET_KEY"
		panic(message)
	}

	bucketname = os.Getenv("LEARNR_BUCKET")
	AWSAccessKeyId = os.Getenv("AWS_ACCESS_KEY")
	AWSSecretKey = os.Getenv("AWS_SECRET_KEY")
}

func placeAmazonFile(amazonFileLocation string, userid string, learnrID string, learnrinfoid string) (bool, string, string) {
	goodFileGet, returnMessage, filePlacement := true, "Working file created successfully", ""

	//Make the potential directory
	curDir, _ := os.Getwd()
	tempFileLocation := filepath.Join(curDir, userid, learnrID, learnrinfoid)
	os.MkdirAll(tempFileLocation, 0777)
	//Create file to copy to and place in temp directory
	item := "fileMove.xlsx"
	file, err := os.Create(item)
	if err != nil {
		logMsg := "Error creating file: " + err.Error()
		fmt.Println(logMsg)
		goodFileGet, returnMessage, filePlacement = false, logMsg, ""
		return goodFileGet, returnMessage, filePlacement
	}

	//Initialize secret keys
	os.Setenv("AWS_ACCESS_KEY", AWSAccessKeyId)
	os.Setenv("AWS_SECRET_KEY", AWSSecretKey)

	fmt.Printf("DEBUG: Here is access key: %v\nHere is secret key: %v\nHere is filelocation: %v\n", AWSAccessKeyId,
		AWSSecretKey, amazonFileLocation)
	sess, _ := session.NewSession(&aws.Config{
		Region:                         aws.String("us-east-2"),
		Credentials:                    credentials.NewEnvCredentials(),
		DisableRestProtocolURICleaning: aws.Bool(true),
	})

	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketname),
			Key:    aws.String(amazonFileLocation),
		})
	if err != nil {
		logMsg := "Error downloading file from Amazon in placeAmazonFile: \n" + err.Error() + "\n"
		fmt.Println(logMsg)
		logWriter(logMsg)
		goodFileGet, returnMessage, filePlacement = false, logMsg, ""
		return goodFileGet, returnMessage, filePlacement
	} else {
		sucMsg := "Downloaded: " + file.Name() + string(numBytes) + " bytes" + "\n"
		fmt.Println(sucMsg)
		logWriter(sucMsg)
	}
	file.Close() //Closes the file in order to move it

	//Open this file again to move
	readFile, err := os.Open(file.Name())
	if err != nil {
		theErr := "Error opening this file: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodFileGet, returnMessage, filePlacement = false, theErr, ""
		return goodFileGet, returnMessage, filePlacement
	}
	writeToFile, err := os.Create(tempFileLocation)
	if err != nil {
		theErr := "Error creating writeToFile: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodFileGet, returnMessage, filePlacement = false, theErr, ""
		return goodFileGet, returnMessage, filePlacement
	}
	//Move file Contents to folder
	n, err := io.Copy(writeToFile, readFile)
	if err != nil {
		theErr := "Error copying the contents of the one Excel sheet to another: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodFileGet, returnMessage, filePlacement = false, theErr, ""
		return goodFileGet, returnMessage, filePlacement
	}
	fmt.Printf("DEBUG: move the contents of n: %v\n", n)
	readFile.Close()    //Close File
	writeToFile.Close() //Close File
	//Delete created file
	removeErr := os.Remove(file.Name())
	if removeErr != nil {
		theErr := "Error removing the file: " + removeErr.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodFileGet, returnMessage, filePlacement = false, theErr, ""
		return goodFileGet, returnMessage, filePlacement
	}

	/* Set returned variables */

	return goodFileGet, returnMessage, filePlacement
}
