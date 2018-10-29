package jira

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var testNotifyNotifications = []Notification{
	{
		Title: "Testovací Notifikace",
	},
}

func alerterMock(title, message, appIcon string) error {
	log.Println(fmt.Sprintf("%s: Application %s just alerted message \"%s\"", appIcon, title, message))
	return nil
}

func TestNotificator_Notify(t *testing.T) {
	notificator := &notificator{
		alerterMock,
	}

	err := notificator.notify(testNotifyNotifications)
	assert.Equal(t, nil, err)
}
