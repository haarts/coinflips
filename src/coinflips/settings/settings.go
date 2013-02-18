package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type settings struct {
	SmtpUser string
	SmtpPassword string
	Salt string
}

func ReadSmtpUser() string {
	return readSettings().SmtpUser
}

func ReadSmtpPassword() string {
	return readSettings().SmtpPassword
}

func ReadSalt() string {
	return readSettings().Salt
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
