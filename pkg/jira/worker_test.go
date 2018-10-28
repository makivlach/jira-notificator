package jira

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testWorkerNotificationCount = 2
)

var testWorkerNotifications = []Notification{
	{
		Title: "Testovací Notifikace",
	},
}

func mockStateFunc(f *notificationWorker) stateFunc {
	return nil
}

func mockFetchNotificationCountStateFunc(f *notificationWorker) stateFunc {
	_ = fetchNotificationCount(f)
	return nil
}

func mockFetchNotificationsStateFunc(f *notificationWorker) stateFunc {
	_ = fetchNotifications(f)
	return nil
}

type mockClient struct{}

func (mockClient) FetchNotificationCount() (int, error) {
	return testWorkerNotificationCount, nil
}

func (mockClient) FetchNotifications() ([]Notification, error) {
	return testWorkerNotifications, nil
}
func (mockClient) Login(username, password string) error {
	return nil
}

func TestNotificationWorker_Start(t *testing.T) {
	notificationChan := make(chan []Notification)
	finishedChan := make(chan bool)

	worker := &notificationWorker{
		c:        &client{},
		channel:  notificationChan,
		finished: finishedChan,
		state:    mockStateFunc,
	}

	go worker.Start(0)
	isFinished := <-finishedChan

	assert.Equal(t, true, isFinished)
}

func TestNotificationWorker_fetchNotificationCount(t *testing.T) {
	notificationChan := make(chan []Notification)
	finishedChan := make(chan bool)

	worker := &notificationWorker{
		c:        mockClient{},
		channel:  notificationChan,
		finished: finishedChan,
		state:    mockFetchNotificationCountStateFunc,
	}

	go worker.Start(0)
	<-finishedChan
	if assert.NoError(t, worker.e) {
		assert.Equal(t, testWorkerNotificationCount, worker.notificationCount)
	}
}

func TestNotificationWorker_fetchNotifications(t *testing.T) {
	notificationChan := make(chan []Notification)
	finishedChan := make(chan bool)

	worker := &notificationWorker{
		c:        mockClient{},
		channel:  notificationChan,
		finished: finishedChan,
		notificationData: []Notification{
			{
				Title: "Testovací Notifikace 2",
			},
		},
		state: mockFetchNotificationsStateFunc,
	}

	go worker.Start(0)

	var notifications []Notification
	select {
	case notifications = <-notificationChan:
	}

	if assert.NoError(t, worker.e) {
		assert.Equal(t, testWorkerNotifications, notifications)
	}
}
