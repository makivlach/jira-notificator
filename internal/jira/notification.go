package jira

import (
	"github.com/gen2brain/beeep"
	"github.com/hajimehoshi/oto"
	"github.com/pkg/errors"
	"log"
	"time"
)

const (
	errorWrapAudioInit        = "An error occurred while initializing audio device"
	errorWrapConnection       = "An error occurred while trying to establish connection for notification data"
	errorWrapPlayTone         = "An error occurred while trying to play the notification tone"
	errorWrapFireNotification = "An error occurred while tying to fire the notifications"
	errorWrapTryUpdateData    = "An error occurred while trying to update data"
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

	player, err := oto.NewPlayer(44100, 2, 2, 40000)
	if err != nil {
		return errors.Wrap(err, errorWrapAudioInit)
	}

	worker, err := NewWorker(c, channel, finished)
	if err != nil {
		return errors.Wrap(err, errorWrapConnection)
	}

	notificator := &notificator{beeep.Alert}
	go worker.Start(data.Interval)

	for {
		select {
		case notifications := <-channel:
			log.Println(data.Text)

			_, err := player.Write(data.Sound)
			if err != nil {
				return errors.Wrap(err, errorWrapPlayTone)
			}

			err = notificator.notify(notifications)
			if err != nil {
				return errors.Wrap(err, errorWrapFireNotification)
			}
		case <-finished:
			return errors.Wrap(worker.e, errorWrapTryUpdateData)
		}
	}
}
