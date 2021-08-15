package main

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

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

/* This analyzes an Excel sheet to make sure it is formatted correctly
for sending LearnRs. If failed, it will return a message and bool */
func examineExcelSheet(excelPath string, fileName string) (bool, string) {
	goodExcel, message := true, ""
	excelErrors := []string{} //This is collected and put into our message variable at the end
	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		errMsg := "Issue opening Excel sheet: " + err.Error()
		fmt.Println(errMsg)
		goodExcel, message = false, errMsg
		return goodExcel, message
	}
	/* Check to see if  first few cells are formatted correctly */
	//Person Name
	cell, err := f.GetCellValue("Sheet1", "A1")
	if err != nil {
		errMsg := "Error working with this Excel Sheet: " + err.Error()
		fmt.Println(errMsg)
		goodExcel, message = false, errMsg
		return goodExcel, message
	} else if !(strings.ToLower(cell) == "person name") {
		errMsg := "Error working with this Excel Sheet: " + "A1 must be 'person name'"
		fmt.Println(errMsg)
		goodExcel, message = false, errMsg
		return goodExcel, message
	}
	//Phone Number
	cell, err = f.GetCellValue("Sheet1", "B1")
	if err != nil {
		errMsg := "Error working with this Excel Sheet: " + err.Error()
		fmt.Println(errMsg)
		goodExcel, message = false, errMsg
		return goodExcel, message
	} else if !(strings.ToLower(cell) == "phone number") {
		errMsg := "Error working with this Excel Sheet: " + "B1 must be 'phone number'"
		fmt.Println(errMsg)
		goodExcel, message = false, errMsg
		return goodExcel, message
	}
	//What to Say
	cell, err = f.GetCellValue("Sheet1", "C1")
	if err != nil {
		errMsg := "Error working with this Excel Sheet: " + err.Error()
		fmt.Println(errMsg)
		goodExcel, message = false, errMsg
		return goodExcel, message
	} else if !(strings.ToLower(cell) == "what to say") {
		errMsg := "Error working with this Excel Sheet: " + "C1 must be 'what to say'"
		fmt.Println(errMsg)
		goodExcel, message = false, errMsg
		return goodExcel, message
	}
	/* Check through person name to make sure values are okay*/
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
	}
	theRows := 0
	for _, row := range rows {
		//Loop through columns to check each field
		if theRows == 0 {

		} else {
			theColumns := 0
			for _, colCell := range row {
				//For each column case, check each column
				switch theColumns {
				case 0:
					//Check Person Name
					personName := colCell
					if len(personName) > 20 {
						theErr := "Error; person name is too long, needs to be under 20 characters: " + personName
						excelErrors = append(excelErrors, theErr)
						goodExcel = false
					} else if len(personName) <= 0 {
						theErr := "Error; person name is too short, needs to be at least 1 character: " + personName
						excelErrors = append(excelErrors, theErr)
						goodExcel = false
					} else {

					}
					break
				case 1:
					//Check Phone Number
					personPhone := colCell
					if len(personPhone) > 11 {
						theErr := "Error; person phone number is too long, needs to be under 11 characters: " + personPhone
						excelErrors = append(excelErrors, theErr)
						goodExcel = false
					} else if len(personPhone) <= 0 {
						theErr := "Error; person phone number is too short, needs to be at least 1 character: " + personPhone
						excelErrors = append(excelErrors, theErr)
						goodExcel = false
					} else if personPhone == "911" {
						theErr := "Error; cannot use emergency numbers for phone number: " + personPhone
						excelErrors = append(excelErrors, theErr)
						goodExcel = false
					} else {

					}
					break
				case 2:
					//Check what to Say
					personSay := colCell
					if len(personSay) > 120 {
						theErr := "Error; message to user cannot be larger than 120 characters: " + personSay
						excelErrors = append(excelErrors, theErr)
						goodExcel = false
					} else if len(personSay) <= 0 {
						theErr := "Error; person message is too short, needs to be at least 1 character: " + personSay
						excelErrors = append(excelErrors, theErr)
						goodExcel = false
					} else {

					}
				default:
					//Wrong column, there's an issue
					theErr := "Error; column distribution is incorrect. Please contain all data in the first 3 columns"
					excelErrors = append(excelErrors, theErr)
					goodExcel = false
				}
				theColumns = theColumns + 1 //Increment column counter for logic above
			}
		}
		theRows = theRows + 1
	}

	//Format message to display the errors
	if !goodExcel {
		message = "There were errors with the Excel sheet, please review and submit again: \n"
		for n := 0; n < len(excelErrors); n++ {
			message = message + excelErrors[n] + "\n"
		}
	} else {
		message = "Excel sheet was successful; here are any errors returned: "
		for n := 0; n < len(excelErrors); n++ {
			message = message + excelErrors[n] + "\n"
		}
	}

	return goodExcel, message
}

/* This function sends our Excel sheet with multiple people
to send LearnRs to in Amazon buckets. It will be worked by our 'texting project'
Microservice, then deleted afterwards */
func sendExcelToBucket(aHex string, s *session.Session,
	file multipart.File, fileHeader *multipart.FileHeader, aUser User) (bool, string, string) {
	goodSend, message := true, ""

	// the file content into a buffer
	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	// create a unique file name for the file
	stringUserID := strconv.Itoa(aUser.UserID)
	tempFileName := "excelSheets/" + stringUserID + "/" + aHex + filepath.Ext(fileHeader.Filename)

	/* Upload function for certain content type */
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucketname),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	if err != nil {
		errMsg := "Error submitting file for Amazon bucket in UploadFileToS3: " + err.Error()
		logWriter(errMsg)
		fmt.Printf("Error submitting file for Amazon bucket in UploadFileToS3: \n%v\n", err.Error())
		message = errMsg
		goodSend = false
	}

	return goodSend, message, tempFileName
}
