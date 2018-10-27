package main

import "github.com/vlachmilan/jira-notificator/jira"

type alerter func(title, message, appIcon string) error

type notificator struct {
	alert alerter
}

func (n notificator) notify(notification jira.Notification) error {
	return n.alert("Jira", notification.Title, "assets/information.png")
}
