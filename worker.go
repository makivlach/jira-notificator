package main

import (
	"errors"
	"github.com/adam-hanna/arrayOperations"
	"github.com/vlachmilan/jira-notificator/jira"
	"log"
	"time"
)

type stateFunc func(f *notificationWorker) stateFunc

type notificationWorker struct {
	state             stateFunc
	c                 jira.Client
	notificationCount int
	notificationData  *jira.Notifications
	e                 error
	channel           chan *jira.Notifications
}

func (w *notificationWorker) start(refreshInterval time.Duration) {
	state := w.state
	for state != nil {
		state = state(w)
		time.Sleep(time.Second * refreshInterval)
	}
	close(w.channel)

	if w.e != nil {
		log.Fatalln(w.e)
	}
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

	notifications, ok = z.Interface().(*jira.Notifications)
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
