package mail

import (
	"net/smtp"
	"io/ioutil"
	"net/http"
	"strings"
	"fmt"
	"encoding/json"
	"../database"
)

type MailingCoinflip struct {
	database.Coinflip
}

type Message struct {
	Sender string
	To []string
	Subject string
	Body string
}

type settings struct {
	SmtpUser string
	SmtpPassword string
}

var (
	smtpUser = readSmtpUser()
	smtpPassword = readSmtpPassword()
)

const (
	smtpServer = "email-smtp.us-east-1.amazonaws.com"
	senderEmail = "harm@mindshards.com"
)

const resultMessage = `
Brilliant! Everybody checked in.
The result of your coin flip is:

%s

Remember this results is based on absolute randomness.

Thanks for using Coinflips.net!

`

const confirmMessage = `
Someone created a coin toss with you.
Please confirm your email address by clicking on the link below:

%s
`

func NewMailingCoinflip(coinflipKey string) (*MailingCoinflip) {
	coin, _ := database.FindCoinflip(coinflipKey)
	return &MailingCoinflip{Coinflip: *coin}
}

func (coinflip *MailingCoinflip) MailResultToParticipants() string {
	result, _ := coinflip.getResult()
	participants := coinflip.FindParticipants()

	for _, participant := range participants {
		message := Message{
			Sender:  senderEmail,
			To:      []string{participant.Email},
			Subject: "The results are in!",
			Body:    fmt.Sprintf(resultMessage, result),
		}
		if err := message.send(); err != nil {
			fmt.Printf("Couldn't send email: %v\n", err)
		}
	}
	return result
}

func (coinflip *MailingCoinflip) MailCreationToParticipants() {
	participants := coinflip.FindParticipants()

	for _, participant := range participants {
		message := Message{
			Sender:  senderEmail,
			To:      []string{participant.Email},
			Subject: "What will it be? " + coinflip.Head + " or " + coinflip.Tail + "?",
			Body:    fmt.Sprintf(confirmMessage, "http://www.coinflips.net/register/" + coinflip.EncodedKey() + "?email=" + participant.Email),
		}
		if err := message.send(); err != nil {
			fmt.Printf("Couldn't send email: %v\n", err)
		}
	}
}

func (coinflip *MailingCoinflip) getResult() (string, error) {
	response, err := http.Get("http://www.random.org/integers/?num=1&min=0&max=1&col=1&base=10&format=plain&rnd=new")
	if err != nil {
		return "", err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		strContents := strings.Trim(string(contents), " \n")
		if strContents == "0" {
			return coinflip.Head, nil
		} else if strContents == "1" {
			return coinflip.Tail, nil
		} else {
			return "weirdness", nil
		}
	}
	return "never happens", nil
}

func (message *Message) send() error {
	auth := smtp.PlainAuth(
		"",
		smtpUser,
		smtpPassword,
		smtpServer,
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		smtpServer + ":587",
		auth,
		senderEmail,
		message.To,
		[]byte("Subject: " + message.Subject + "\r\n" +
		"From: Coinflips.net <" + senderEmail + ">\r\n" +
		"Reply-To: " + senderEmail + "\r\n" +
		"To: " + message.To[0] + "\r\n" +
		"\r\n" + message.Body),
	)
	if err != nil {
		return err
	}
	return nil
}

func readSmtpUser() string {
	return readSettings().SmtpUser
}

func readSmtpPassword() string {
	return readSettings().SmtpPassword
}

func readSettings() settings {
	file, err := ioutil.ReadFile("./settings.json")
	if err != nil {
		fmt.Errorf("File error: %v\n", err)
	}
	var s settings
	err = json.Unmarshal(file, &s)
	if err != nil {
	  fmt.Errorf("File failed to decode: %v\n", err)
	}
	return s
}
