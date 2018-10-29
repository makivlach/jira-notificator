package jira

import (
	"github.com/gen2brain/beeep"
	"log"
	"time"
)

type alerter func(title, message, appIcon string) error

type NotificationData struct {
	Sound    []byte
	Interval time.Duration
	Text     string
}

type notificator struct {
	alert alerter
}

func (n notificator) notify(notifications []Notification) error {
	for _, notification := range notifications {
		err := n.alert("Jira", notification.Title, "assets/information.png")

		if err != nil {
			return err
		}
	}
	return nil
}

func FetchNewNotifications(c Client, data NotificationData) error {
	channel := make(chan []Notification)
	finished := make(chan bool)

	player, err := NewPlayer(data.Sound)
	if err != nil {
		return err
	}

	worker, err := NewWorker(c, channel, finished)
	if err != nil {
		return err
	}

	notificator := &notificator{beeep.Alert}
	go worker.Start(data.Interval)

	for {
		select {
		case notifications := <-channel:
			log.Println(data.Text)

			err := notificator.notify(notifications)
			if err != nil {
				return err
			}

			if err := player.Play(); err != nil {
				return err
			}
		case <-finished:
			return nil
		}
	}
}
