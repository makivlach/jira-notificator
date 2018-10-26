package main

import "github.com/vlachmilan/jira-notificator/client"

type notifier struct {
	data *client.Notifications
}

func (n notifier) isDiferrent(data *client.Notifications) bool {
	return n.data != data
}
