package jira

import (
	"errors"
	"github.com/adam-hanna/arrayOperations"
	"log"
	"time"
)

type stateFunc func(f *notificationWorker) stateFunc

type notificationWorker struct {
	state             stateFunc
	c                 Client
	e                 error
	notificationCount int
	notificationData  *Notifications
	finished          chan bool
	channel           chan *Notifications
}

func NewWorker(client Client, channel chan *Notifications, finished chan bool) *notificationWorker {
	return &notificationWorker{
		state:    fetchNotificationCount,
		c:        client,
		finished: finished,
		channel:  channel,
	}
}

func (w *notificationWorker) Start(refreshInterval time.Duration) {
	state := w.state
	for state != nil {
		state = state(w)
		time.Sleep(time.Second * refreshInterval)
	}
	close(w.channel)

	if w.e != nil {
		log.Fatalln(w.e)
	}

	w.finished <- true
}

func fetchNotifications(f *notificationWorker) stateFunc {
	notifications, err := f.c.FetchNotifications()
	if err != nil {
		f.e = err
		return nil
	}

	z, ok := arrayOperations.Difference(f.notificationData, notifications)
	if !ok {
		f.e = errors.New("cannot find difference")
		return nil
	}

	notifications, ok = z.Interface().(*Notifications)
	if !ok {
		f.e = errors.New("cannot convert new notifications to an object")
		return nil
	}

	f.channel <- notifications

	return fetchNotificationCount
}

func fetchNotificationCount(f *notificationWorker) stateFunc {
	count, err := f.c.FetchNotificationCount()
	if err != nil {
		f.e = err
		return nil
	}

	if count != f.notificationCount {
		f.notificationCount = count
		return fetchNotifications
	}

	return fetchNotificationCount
}
