package main

import (
	"bytes"
	"github.com/200sc/klangsynthese/mp3"
	"github.com/gen2brain/beeep"
	"github.com/vlachmilan/jira-notificator/pkg/jira"
	"gopkg.in/AlecAivazis/survey.v1"
	"log"
	"os"
)

const (
	notificationInterval = 10
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

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}

	c := jira.NewClient(answers.Host)

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
	channel := make(chan []jira.Notification)
	finished := make(chan bool)

	worker, err := jira.NewWorker(c, channel, finished)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	notificator := jira.NewNotificator(beeep.Alert)
	go worker.Start(notificationInterval)

	for {
		select {
		case notifications := <-channel:
			log.Println("Provedení aktualizace")

			err := notificator.Notify(notifications)
			if err != nil {
				log.Fatalln(err)
				os.Exit(1)
			}

			if err := playNotificationMusic(); err != nil {
				log.Fatalln(err)
				os.Exit(1)
			}
		case <-finished:
			return
		}
	}
}

func playNotificationMusic() error {
	file, err := Asset(notificationSound)
	if err != nil {
		return err
	}

	reader := &reader{bytes.NewReader(file)}
	a, err := mp3.Load(reader)
	return <-a.Play()
}

type reader struct {
	*bytes.Reader
}

func (*reader) Close() error {
	return nil
}
