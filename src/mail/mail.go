package mail

import (
	"net/smtp"
	"io/ioutil"
	"net/http"
	"strings"
	"strings"
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

func (coinflip *DecoratedCoinflip) mailResultToParticipants() string {
	result := coinflip.getResult()
	participants := coinflip.FindParticipants()

	for participant := range participants {
		msg := &mail.Message{
			Sender:  senderEmail,
			To:      []string{participant.Email},
			Subject: "The results are in!",
			Body:    fmt.Sprintf(resultMessage, result),
		}
		if err := mail.Send(msg); err != nil {
			context.Errorf("Couldn't send email: %v", err)
		}
	}
	return result
}

func (coinflip *DecoratedCoinflip) mailCreationToParticipants() {
	participants := coinflip.FindParticipants()
	query := datastore.NewQuery("Participant").Ancestor(coinflipKey)

	for participant := range participants {
		msg := &mail.Message{
			Sender:  senderEmail,
			To:      []string{participant.Email},
			Subject: "What will it be? " + coinflip.Head + " or " + coinflip.Tail + "?",
			Body:    fmt.Sprintf(confirmMessage, "http://www.coinflips.net/register/"+coinflipKey.Encode()+"?email="+participant.Email),
		}
		if err := mail.Send(context, msg); err != nil {
			context.Errorf("Couldn't send email: %v", err)
		}
	}
}

func (coinflip *DecoratedCoinflip) getResult() result, error {
	client := urlfetch.Client(context)
	response, err := client.Get("http://www.random.org/integers/?num=1&min=0&max=1&col=1&base=10&format=plain&rnd=new")
	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
		}
		strContents := strings.Trim(string(contents), " \n")
		if strContents == "0" {
			return coinflip.Head
		} else if strContents == "1" {
			return coinflip.Tail
		} else {
			return "weirdness"
		}
	}
	return "never happens"
}
