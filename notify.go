package main

import "github.com/vlachmilan/jira-notificator/jira"

type alerter func(title, message, appIcon string) error

type notificator struct {
	alert alerter
}

func (n notificator) notify(notifications *jira.Notifications) error {
	for _, notification := range notifications.Notifications {
		err := n.alert("Jira", notification.Title, "assets/information.png")

		if err != nil {
			return err
		}
	}
	return nil
}
