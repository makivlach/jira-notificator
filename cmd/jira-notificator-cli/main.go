package main

import (
	"github.com/gen2brain/beeep"
	"github.com/vlachmilan/jira-notificator/pkg/jira"
	"gopkg.in/AlecAivazis/survey.v1"
	"log"
	"os"
)

const (
	notificationInterval = 10
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

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}

	c := jira.New(answers.Host)

	login(c, answers.Username, answers.Password)
	fetchNewNotifications(c)

	os.Exit(0)
}

func login(c jira.Client, login, password string) {
	err := c.Login(login, password)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	log.Println("Uživatel přihlášen")
}

func fetchNewNotifications(c jira.Client) {
	channel := make(chan *jira.Notifications)
	finished := make(chan bool)
	worker := jira.NewWorker(c, channel, finished)

	notificator := jira.NewNotificator(beeep.Alert)

	go worker.Start(notificationInterval)

	for !<-finished {
		notifications := <-channel
		log.Println("Provedení aktualizace")

		err := notificator.Notify(notifications)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
	}
}
