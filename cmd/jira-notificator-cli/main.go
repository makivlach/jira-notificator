package main

import (
	"github.com/vlachmilan/jira-notificator/pkg/jira"
	"gopkg.in/AlecAivazis/survey.v1"
	"log"
	"os"
)

const (
	notificationInterval = 15
	notificationSound    = "assets/notify.mp3"
)

// the questions to ask
var qs = []*survey.Question{
	{
		Name:      "host",
		Prompt:    &survey.Input{Message: "Zadejte adresu Jiry (například: https://something.atlassian.net): "},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name:      "username",
		Prompt:    &survey.Input{Message: "Zadejte přihlašovací jméno: "},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name:      "password",
		Prompt:    &survey.Password{Message: "Zadejte přihlašovací heslo: "},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
}

func main() {
	answers := struct {
		Host     string
		Username string
		Password string
	}{}

	//perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}

	c, err := jira.NewClient(answers.Host, answers.Username, answers.Password)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	sound, err := Asset(notificationSound)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	err = c.Login()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	println("Uživatel přihlášen")

	err = jira.FetchNewNotifications(c, jira.NotificationData{Sound: sound, Interval: notificationInterval, Text: "Provedení aktualizace"})
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	os.Exit(0)
}
