package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

/* INFORMATION FOR OUR EMAIL VARIABLES */
var senderAddress string
var senderPWord string

var GmailService *gmail.Service //This gets initialized in init

var theClientID string
var theClientSecret string
var theAccessToken string
var theRefreshToken string
var redirectURL string

//Might be needed to cancel stuff
var TheCancelFunc context.CancelFunc

type UserJSON struct {
	TheName    string `json:"TheName"`
	TheEmail   string `json:"TheEmail"`
	TheMessage string `json:"TheMessage"`
}

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func loadInEmailCreds() {
	//Check to see if ENV Creds are available first
	_, ok := os.LookupEnv("EMAIL_CLIENTID")
	if !ok {
		message := "This ENV Variable is not present: " + "EMAIL_CLIENTID"
		panic(message)
	}
	_, ok2 := os.LookupEnv("EMAIL_CLIENTSECRET")
	if !ok2 {
		message := "This ENV Variable is not present: " + "EMAIL_CLIENTSECRET"
		panic(message)
	}
	_, ok3 := os.LookupEnv("EMAIL_ACCESSTOKEN")
	if !ok3 {
		message := "This ENV Variable is not present: " + "EMAIL_ACCESSTOKEN"
		panic(message)
	}
	_, ok4 := os.LookupEnv("EMAIL_REFRESHTOKEN")
	if !ok4 {
		message := "This ENV Variable is not present: " + "EMAIL_REFRESHTOKEN"
		panic(message)
	}

	theClientID = os.Getenv("EMAIL_CLIENTID")
	theClientSecret = os.Getenv("EMAIL_CLIENTSECRET")
	theAccessToken = os.Getenv("EMAIL_ACCESSTOKEN")
	theRefreshToken = os.Getenv("EMAIL_REFRESHTOKEN")
}

//Initialized at begininning of program
func OAuthGmailService() {
	config := oauth2.Config{
		ClientID:     theClientID,
		ClientSecret: theClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURL,
	}

	token := oauth2.Token{
		AccessToken:  theAccessToken,
		RefreshToken: theRefreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	//Create a context to use for our gmail services

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		errMsg := "Unable to retrieve Gmail client: " + err.Error()
		fmt.Println(errMsg)
		logWriter(errMsg)
		panic(errMsg)
	}

	GmailService = srv
	if GmailService != nil {
		succMsg := "Email service is initialized"
		logWriter(succMsg)
	} else {
		panic("GmailService is nil")
	}
}

/* This funciton is taken from Ajax on messageme page to process successful emails */
func emailMe(w http.ResponseWriter, r *http.Request) {
	//Get the byte slice from the JSON
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}

	//Declare DataType from JSON
	type MessageInfo struct {
		YourNameInput    string `json:"YourNameInput"`
		YourEmailInput   string `json:"YourEmailInput"`
		YourMessageInput string `json:"YourMessageInput"`
		YourUserID       int    `json:"YourUserID"`
		YourUser         User   `json:"YourUser"`
	}

	//Marshal the user data into our type
	var dataEmail MessageInfo
	json.Unmarshal(bs, &dataEmail)

	//Declare return information for JSON
	type ReturnMessage struct {
		TheErr     string `json:"TheErr"`
		ResultMsg  string `json:"ResultMsg"`
		SuccOrFail int    `json:"SuccOrFail"`
	}
	theReturnMessage := ReturnMessage{}

	/* Set data for email function, then send an email */
	userMarsh, _ := json.Marshal(dataEmail.YourUser)
	theMessage := "Message from " + dataEmail.YourNameInput + ":\n" +
		dataEmail.YourMessageInput + "\n" + "Here is the User info: \n" + string(userMarsh) + "\n"
	theSubject := "Posted Message From " + dataEmail.YourNameInput
	theRequest := NewRequest([]string{senderAddress}, theSubject, "Hello, World!")
	sendResult, goodSend := sendEmail(senderAddress, theMessage, theSubject, theRequest)
	//Set return JSON
	if goodSend != true {
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = "Message to me failed"
		theReturnMessage.TheErr = "Message to me Failed: \n" + sendResult
	} else {
		theReturnMessage.SuccOrFail = 0
		theReturnMessage.ResultMsg = "Message sent to me"
		theReturnMessage.TheErr = ""
	}
	//Return JSON
	theJSONMessage, err := json.Marshal(theReturnMessage)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* This sends the actual email out to the intended party */
func sendEmail(sendTo string, themessage string, thesubject string, theRequest *Request) (string, bool) {
	sendResult, goodSend := "", true // Used to determine if our email was sent successfully

	var message gmail.Message
	emailTo := "To: " + sendTo + "\r\n"
	subject := "Subject: " + thesubject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + themessage)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	//Create a context to use for our gmail services
	/*
		parentContext := context.Background()
		aContext, cancel := context.WithTimeout(parentContext, (3 * time.Second))

		TheCancelFunc = cancel
	*/

	// Send the message
	_, err := GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		errMsg := "Error sending this message to the User: " + err.Error()
		sendResult, goodSend = errMsg, false
		fmt.Println(errMsg)
		logWriter(errMsg)
	} else {
		errMsg := "Succussfully sent email to: " + sendTo
		sendResult, goodSend = errMsg, true
	}

	return sendResult, goodSend
}

/* This is a test function to send email with a TEMPLATE, not just text */
func sendEmailTemplate(sendTo string, themessage string, thesubject string, theRequest *Request) (string, bool) {
	sendResult, goodSend := "", true // Used to determine if our email was sent successfully

	var message gmail.Message
	emailTo := "To: " + sendTo + "\r\n"
	subject := "Subject: " + thesubject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + theRequest.body)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	//Create a context to use for our gmail services
	/*
		parentContext := context.Background()
		aContext, cancel := context.WithTimeout(parentContext, (3 * time.Second))

		TheCancelFunc = cancel
	*/

	// Send the message
	_, err := GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		errMsg := "Error sending this message to the User: " + err.Error()
		sendResult, goodSend = errMsg, false
		fmt.Println(errMsg)
		logWriter(errMsg)
	} else {
		errMsg := "Succussfully sent email to: " + sendTo
		sendResult, goodSend = errMsg, true
	}

	return sendResult, goodSend
}

/* This makes a new request...not sure it's needed */
func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

//Needed to parse email templates
func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
