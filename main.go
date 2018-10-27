package main

import (
	"github.com/gen2brain/beeep"
	"github.com/vlachmilan/jira-notificator/jira"
	"log"
	"os"
)

const (
	notificationHost     = "https://something.atlassian.net"
	notificationInterval = 10
)

func main() {
	c := jira.New(notificationHost)

	login(c)
	fetchNewNotifications(c)

	os.Exit(0)
}

func login(c jira.Client) {
	err := c.Login("", "")
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	log.Println("Uživatel přihlášen")
}

func fetchNewNotifications(c jira.Client) {
	channel := make(chan *jira.Notifications)
	finished := make(chan bool)
	worker := NewWorker(c, channel, finished)

	notificator := &notificator{beeep.Alert}

	go worker.start(notificationInterval)

	for !<-finished {
		notifications := <-channel
		log.Println("Provedení aktualizace")

		err := notificator.notify(notifications)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
	}
}
