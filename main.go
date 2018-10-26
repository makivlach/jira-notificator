package main

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"github.com/vlachmilan/jira-notificator/client"
	"log"
	"os"
	"time"
)

const (
	notificationHost     = "https://predplatit.atlassian.net"
	notificationInterval = 30
)

func main() {
	c := client.New(notificationHost)

	err := c.Login("", "")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Print("Uživatel přihlášen")

	notificationChannel := make(chan *client.Notifications)

	go func() {
		for {
			notificationChannel <- c.FetchNewNotifications()
			fmt.Print("Provedena aktualizace")
			time.Sleep(time.Second * notificationInterval)
		}
	}()

	var notifications *client.Notifications
	for {
		notifications = <-notificationChannel

		for _, v := range notifications.Notifications {
			err := beeep.Alert("Jira", v.Title, "assets/information.png")
			if err != nil {
				log.Fatal(err)
			}
			break
		}

	}

	os.Exit(0)
}
