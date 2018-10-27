package jira

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNotificationWorker_Start(t *testing.T) {
	notificationChan := make(chan *Notifications)
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

func mockStateFunc(f *notificationWorker) stateFunc {
	return nil
}
