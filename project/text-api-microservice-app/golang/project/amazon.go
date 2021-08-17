package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

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

/* This deletes our Amazon file after it has been read into our session */
func deleteExcelFromAmazon(amazonExcelLocation string) {
	os.Setenv("AWS_ACCESS_KEY", AWSAccessKeyId)
	os.Setenv("AWS_SECRET_KEY", AWSSecretKey)
	sess, err := session.NewSession(&aws.Config{
		Region:                         aws.String("us-east-2"),
		Credentials:                    credentials.NewEnvCredentials(),
		DisableRestProtocolURICleaning: aws.Bool(true),
	},
	)
	if err != nil {
		theErr := "Error connection to Amazon session in deleteFileFromS3: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		return
	}
	// Create S3 service client
	svc := s3.New(sess)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucketname), Key: aws.String(amazonExcelLocation)})
	if err != nil {
		theErr := "Whoops, error deleteing object: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		return
	}
	//Delete Confirmation
	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(amazonExcelLocation),
	})
	if err != nil {
		theErr := "Error waiting until document deletes: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		return
	} else {
		giveString := "Object successfully deleted in Amazon: " + amazonExcelLocation
		fmt.Println(giveString)
	}
}

/* This deletes the Excel sheet locally in our Server */
func deleteLocalFile(workFileLocation string) {
	fileError := os.Remove(workFileLocation)
	if fileError != nil {
		theErr := "Error removing file at this location: " + fileError.Error()
		fmt.Println(theErr)
		logWriter(theErr)
	}
}

func placeAmazonFile(amazonFileLocation string, theFileName string,
	userid string, learnrID string, learnrinfoid string) (bool, string, string) {
	goodFileGet, returnMessage, filePlacement := true, "Working file created successfully", ""

	//Make the potential directory
	curDir, _ := os.Getwd()
	tempFileLocation := filepath.Join(curDir, "aws-workfiles", userid, learnrID, learnrinfoid)
	os.MkdirAll(tempFileLocation, 0777)
	//Create file to copy to and place in temp directory
	item := theFileName
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

	sourceFileStat, err := os.Stat(theFileName)
	if err != nil {
		theErr := "Error with sourceStat: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodFileGet, returnMessage, filePlacement = false, theErr, ""
		return goodFileGet, returnMessage, filePlacement
	}

	if !sourceFileStat.Mode().IsRegular() {
		theErr := "Error with regular sourceStat: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodFileGet, returnMessage, filePlacement = false, theErr, ""
		return goodFileGet, returnMessage, filePlacement
	}

	source, err := os.Open(theFileName)
	if err != nil {
		theErr := "Error with file opening: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodFileGet, returnMessage, filePlacement = false, theErr, ""
		return goodFileGet, returnMessage, filePlacement
	}

	dst := filepath.Join(tempFileLocation, theFileName)
	destination, err := os.Create(dst)
	if err != nil {
		theErr := "Error with destination creation: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodFileGet, returnMessage, filePlacement = false, theErr, ""
		return goodFileGet, returnMessage, filePlacement
	}
	nBytes, err := io.Copy(destination, source)
	if err != nil {
		theErr := "Error copying files: " + err.Error() + strconv.Itoa(int(nBytes))
		fmt.Println(theErr)
		logWriter(theErr)
		goodFileGet, returnMessage, filePlacement = false, theErr, ""
		return goodFileGet, returnMessage, filePlacement
	}

	/* Close our files */
	source.Close()
	destination.Close()

	//Delete created file
	removeErr := os.Remove(theFileName)
	if removeErr != nil {
		theErr := "Error deleting file: " + removeErr.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodFileGet, returnMessage, filePlacement = false, theErr, ""
		return goodFileGet, returnMessage, filePlacement
	}

	/* Set returned variables, (the destination should be where our file is) */
	return goodFileGet, returnMessage, dst
}
