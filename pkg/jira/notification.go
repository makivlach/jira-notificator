package jira

import (
	"github.com/gen2brain/beeep"
	"github.com/hajimehoshi/oto"
	"github.com/pkg/errors"
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

	player, err := oto.NewPlayer(44100, 2, 2, 40000)
	if err != nil {
		return errors.Wrap(err, "Chyba při inicializaci zvukového zařízení")
	}

	worker, err := NewWorker(c, channel, finished)
	if err != nil {
		return errors.Wrap(err, "Chyba navázání připojení na data notifikací")
	}

	notificator := &notificator{beeep.Alert}
	go worker.Start(data.Interval)

	for {
		select {
		case notifications := <-channel:
			log.Println(data.Text)

			_, err := player.Write(data.Sound)
			if err != nil {
				return errors.Wrap(err, "Chyba při přehrávání tónu notifikace")
			}

			err = notificator.notify(notifications)
			if err != nil {
				return errors.Wrap(err, "Chyba při vytváření notifikace")
			}
		case <-finished:
			return errors.Wrap(worker.e, "Chyba při provádění aktualizace dat")
		}
	}
}
