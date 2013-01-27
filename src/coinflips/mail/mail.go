package mail

import (
	"net/smtp"
	"io/ioutil"
	"net/http"
	"strings"
	"fmt"
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

const (
	smtpUser = "AKIAJNQ3R2RRD6DC7YBA"
	smtpPassword = "ApI4BKG6pkRiN+LeB2S8o/mz7cucsAO+QFDffbo3LbpH"
	smtpServer = "email-smtp.us-east-1.amazonaws.com"
	senderEmail = "Coinflips.net <harm@mindshards.com>"
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
			fmt.Errorf("Couldn't send email: %v", err)
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
			fmt.Errorf("Couldn't send email: %v", err)
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
			fmt.Printf("%s", err)
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
		smtpServer + ":25",
		auth,
		"harmaarts@gmail.com",
		[]string{"harmaarts@gmail.com"},
		[]byte(message.Subject + "\r\n\r\n" + message.Body),
	)
	if err != nil {
		return err
	}
	return nil
}
