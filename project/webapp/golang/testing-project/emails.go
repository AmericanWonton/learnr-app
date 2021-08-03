package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
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

//This is used for email verification
type EmailVerify struct {
	Username string    `json:"Username"`
	Email    string    `json:"Email"`
	ID       int       `json:"ID"`
	TimeMade time.Time `json:"TimeMade"`
	Active   bool      `json:"Active"`
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
	_, ok5 := os.LookupEnv("MY_EMAIL")
	if !ok5 {
		message := "This ENV Variable is not present: " + "MY_EMAIL"
		panic(message)
	}
	_, ok6 := os.LookupEnv("MY_PWORD")
	if !ok6 {
		message := "This ENV Variable is not present: " + "MY_PWORD"
		panic(message)
	}

	theClientID = os.Getenv("EMAIL_CLIENTID")
	theClientSecret = os.Getenv("EMAIL_CLIENTSECRET")
	theAccessToken = os.Getenv("EMAIL_ACCESSTOKEN")
	theRefreshToken = os.Getenv("EMAIL_REFRESHTOKEN")
	senderAddress = os.Getenv("MY_EMAIL")
	senderPWord = os.Getenv("MY_PWORD")
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
		//Try to send Email to User
		sendUserMessage := "Hello from LearnR! Thanks for reaching out, I am currently reading your message and will respond back soon."
		/* Encode our image to include for email signature */
		templateData := struct {
			Username   string
			TheMessage []string
		}{
			Username:   dataEmail.YourUser.UserName,
			TheMessage: []string{sendUserMessage},
		}
		r := NewRequest([]string{dataEmail.YourUser.Email[0]}, "Thanks, just recieved your message", "Hello, World!")
		err1 := r.ParseTemplate("./static/emailTemplates/messagerecieved.html", templateData)
		if err1 != nil {
			fmt.Printf("Could not parse the template: %v\n", err1.Error())
			log.Fatal("Could not parse the template" + err1.Error())
		}
		_, goodSend := sendEmailTemplate(dataEmail.YourUser.Email[0], sendUserMessage, "Thanks for messaging us!", r)
		if goodSend {
			//Alter messsage
			logWriter("Successfuly email sending")
		} else {
			errMsg := "Failure to send Email"
			fmt.Println(errMsg)
			logWriter(errMsg)
		}
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

/* This function sends a verification email to the User. They
will need to check their email and wait for verification. */
func sendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	//Get the byte slice from the JSON
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}

	//Declare DataType from JSON
	type MessageInfo struct {
		YourNameInput  string `json:"YourNameInput"`
		YourEmailInput string `json:"YourEmailInput"`
		YourUserID     int    `json:"YourUserID"`
		YourUser       User   `json:"YourUser"`
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

	//Make a random number up for the User
	randIDVerif := simpleEmailID()

	//Make email verification and add it to our DB
	verifyEmail := EmailVerify{
		Username: dataEmail.YourNameInput,
		Email:    dataEmail.YourEmailInput,
		ID:       randIDVerif,
		TimeMade: time.Now(),
		Active:   true,
	}

	goodAdd, message := addEmailVerif(verifyEmail)
	if !goodAdd {
		//Email verification failed
		theErr := "Failure to send email verification: " + message
		fmt.Println(theErr)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.TheErr = theErr
		theReturnMessage.ResultMsg = theErr
	} else {
		/* Email verification added to DB. Need to send User email to confirm this */
		type TemplateData struct {
			Username      string `json:"Username"`
			UserID        int    `json:"UserID"`
			Email         string `json:"Email"`
			ItemID        int    `json:"ItemID"`
			DownloadImage string `json:"DownloadImage"`
			PurchaseLink  string `json:"PurchaseLink"`
		}
		templateData := TemplateData{
			Username: dataEmail.YourNameInput,
			UserID:   dataEmail.YourUserID,
			Email:    dataEmail.YourEmailInput,
			ItemID:   verifyEmail.ID,
		}
		//Good to create and send Email template
		emailRequest := NewRequest([]string{dataEmail.YourEmailInput}, "Account Creation", "Your account Creation")
		err1 := emailRequest.ParseTemplate("./static/emailTemplates/emailverif.html", templateData)
		if err1 != nil {
			fmt.Printf("Could not parse the template: %v\n", err1.Error())
			log.Fatal("Could not parse the template" + err1.Error())
		}
		//Send Email
		_, goodSend := sendEmailTemplate(dataEmail.YourEmailInput, "Account Creation",
			"Account Creation", emailRequest)
		if goodSend {
			//Confirmation sent
			theMessage := "Sent confirmation email. Please check to see if you have it in junk or spam. If not sent, please" +
				" check to see if the email entered is valid."
			theReturnMessage.SuccOrFail = 0
			theReturnMessage.ResultMsg = theMessage
		} else {
			errMsg := "Failure to send Email to User about creating their acccount"
			fmt.Println(errMsg)
			theReturnMessage.SuccOrFail = 1
			theReturnMessage.ResultMsg = errMsg
			theReturnMessage.TheErr = errMsg
		}
	}

	//Return JSON
	theJSONMessage, err := json.Marshal(theReturnMessage)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* This creates a short, random ID for User to verify */
func simpleEmailID() int {

	//Create RandomID
	randInt := 0
	randIntString := ""
	min, max := 0, 9 //The min and Max value for our randInt
	//Create the random number, convert it to string
	for i := 0; i < 6; i++ {
		randInt = rand.Intn(max-min) + min
		randIntString = randIntString + strconv.Itoa(randInt)
	}
	//Once we have a string of numbers, we can convert it back to an integer
	theID, err := strconv.Atoi(randIntString)
	if err != nil {
		fmt.Printf("We got an error converting a string back to a number, %v\n", err)
		fmt.Printf("Here is randInt: %v\n and randIntString: %v\n", randInt, randIntString)
		fmt.Println(err)
		log.Fatal(err)
	}

	return theID
}
