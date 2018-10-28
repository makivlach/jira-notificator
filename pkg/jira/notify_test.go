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

func TestNewNotificator(t *testing.T) {
	expected := &Notificator{nil}
	actual := NewNotificator(nil)
	assert.Equal(t, expected, actual)
}

func TestNotificator_Notify(t *testing.T) {
	notificator := &Notificator{
		alerterMock,
	}

	err := notificator.Notify(testNotifyNotifications)
	assert.Equal(t, nil, err)
}
