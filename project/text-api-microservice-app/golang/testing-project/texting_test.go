package main

/* Declarative Test structs/values */

type OurJSON struct {
	TheUser        User       `json:"TheUser"`
	TheLearnR      Learnr     `json:"TheLearnR"`
	TheLearnRInfo  LearnrInfo `json:"TheLearnRInfo"`
	PersonName     string     `json:"PersonName"`
	PersonPhoneNum string     `json:"PersonPhoneNum"`
	Introduction   string     `json:"Introduction"`
}
type LearnRTestSends struct {
	JSONSend            OurJSON
	ExpectedNum         int
	ExpectedTruth       bool
	ExpectedStringArray []string
}

var learnrTestSendResults []LearnRTestSends

func createLearnRTextSession() {
	//Get User

	//Create first LearnR, success
	learnrTestSendResults = append(learnrTestSendResults, LearnRTestSends{
		JSONSend: OurJSON{},
	})
}
