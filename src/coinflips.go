package main

import (
	"fmt"
	"github.com/hoisie/mustache"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type DecoratedCoinflip database.Coinflip

const (
	senderEmail = "Coinflips.net <harm@mindshards.com>"
)

func init() {
	http.HandleFunc("/", home)
	http.HandleFunc("/show/", show)
	http.HandleFunc("/create", create)
	http.HandleFunc("/register/", register)
	http.HandleFunc("/why", why)
	http.HandleFunc("/about", about)
}

func why(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, mustache.RenderFile("./flipco.in/views/layout.html", map[string]string{"body": mustache.RenderFile("./flipco.in/views/why.html", map[string]string{"title": "Why coin tosses? - Coinflips.net"})}))
}

func about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, mustache.RenderFile("./flipco.in/views/layout.html", map[string]string{"body": mustache.RenderFile("./flipco.in/views/about.html", map[string]string{"title": "About coin flips - Coinflips.net"})}))
}

func home(w http.ResponseWriter, r *http.Request) {
	/* static file serve */
	if len(r.URL.Path) != 1 {
		http.ServeFile(w, r, "./flipco.in/views"+r.URL.Path)
		return
	}

	/* the real root */
	count := database.TotalNumberOfParticipants()

	/*long, very long line */
	fmt.Fprint(w, mustache.RenderFile("./flipco.in/views/layout.html", map[string]string{"body": mustache.RenderFile("./flipco.in/views/home.html", map[string]string{"title": "Awesome coin tosses - Coinflips.net", "nr_of_flips": fmt.Sprint(count)})}))
}

func register(w http.ResponseWriter, r *http.Request) {
	/* not a GET request? redirect to home */
	if r.Method != "GET" {
		http.Redirect(w, r, "/", 302)
		return
	}

	coinflipKey, _ := strings.Split(r.URL.Path, "/")[2]
	coinflip, _ := newDecoratedCoinflip(database.FindCoinflip(coinflipKey))

	if coinflip.Result != "" {
		http.Redirect(w, r, "/show/" + coinflipKey.EncodeKey(), 302)
		return
	}

	participant, err := coinflip.FindParticipantByEmail(r.FormValue("email"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	participant.Seen = time.Now()
	participant.Update()

	count, err := coinflip.NumberOfUnregisteredParticipants()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		result := coinflip.mailResultToParticipants(context, coinflipKey)
		coinflip.Result = result
		coinflip.Update()
	}
	http.Redirect(w, r, "/show/" + coinflipKey.Encode(), 302)
}

func create(w http.ResponseWriter, r *http.Request) {
	/* not a POST request? redirect to root */
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 302)
		return
	}

	r.ParseForm()
	tail := r.Form["tail"][0]
	head := r.Form["head"][0]
	friends := r.Form["friends[]"]

	uniq_friends := uniq(friends)

	if tail == "" || head == "" || uniq_friends == nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	coin := database.Coinflip{
		Head: head,
		Tail: tail,
	}

	coinflipKey, err := database.CreateCoinflip(&coin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range uniq_friends {
		participant := Participant{Email: uniq_friends[i]}
		database.CreateParticipant(&participant, &coin)
	}

	coin.mailCreationToParticipants(coinflipKey)

	http.Redirect(w, r, "/show/" + coinflip.EncodeKey(), 302)
}

func show(w http.ResponseWriter, r *http.Request) {
	coinflipKey, _ := strings.Split(r.URL.Path, "/")[2]
	coinflip, _ := database.FindCoinFlip(coinflipKey)

	iterator := database.FindParticipants(coinflip)

	email_list := participantsMap(iterator, func(p Participant) map[string]string {
		var seen_at string
		if p.Seen.IsZero() {
			seen_at = "hasn't registered yet"
		} else {
			seen_at = p.Seen.Format("Monday 2 January 2006")
		}
		return map[string]string{"email": p.Email, "seen_at": seen_at}
	})

	var result string
	if coinflip.Result == "" {
		result = "Nothing yet! Not everybody checked in. Perhaps a little encouragement?"
	} else {
		result = coinflip.Result
	}

	str_to_str := map[string]string{"count": fmt.Sprint(len(email_list)), "head": coinflip.Head, "tail": coinflip.Tail, "result": result}
	str_to_slice := map[string][]map[string]string{"participants": email_list}
	fmt.Fprint(w, mustache.RenderFile("./flipco.in/views/layout.html", map[string]string{"body": mustache.RenderFile("./flipco.in/views/show.html", str_to_str, str_to_slice)}))
}

func uniq(friends []string) (uniq_friends []string) {
	for i := range friends {
		if friends[i] != "" {
			for j := i + 1; j < len(friends); j++ {
				if friends[j] == friends[i] {
					friends[j] = ""
				}
			}
			uniq_friends = append(uniq_friends, friends[i])
		}
	}
	if len(uniq_friends) > 10 {
		uniq_friends = uniq_friends[0:10]
	}
	return uniq_friends
}

func participantsMap(iterator *datastore.Iterator, f func(Participant) map[string]string) (mapped []map[string]string) {
	var participant Participant
	for _, err := iterator.Next(&participant); ; _, err = iterator.Next(&participant) {
		if err == datastore.Done {
			break
		}
		if err != nil {
			break
		}
		mapped = append(mapped, f(participant))
	}
	return mapped
}

func newDecoratedCoinflip(coinflip *database.Coinflip) DecoratedCoinflip {
	return DecoratedCoinflip(coinflip)
}

// TODO: move to mail package
func (coinflip *Coinflip) mailCreationToParticipants() {
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

// TODO: move to mail package
const confirmMessage = `
Someone created a coin toss with you.
Please confirm your email address by clicking on the link below:

%s
`

// TODO: move to mail package
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

// TODO: move to mail package
const resultMessage = `
Brilliant! Everybody checked in.
The result of your coin flip is:

%s

Remember this results is based on absolute randomness.

Thanks for using Coinflips.net!

`

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
