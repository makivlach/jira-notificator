package jira

import (
	"github.com/google/go-cmp/cmp"
	"log"
	"time"
)

type stateFunc func(f *notificationWorker) stateFunc

type notificationWorker struct {
	state             stateFunc
	c                 Client
	e                 error
	notificationCount int
	notificationData  []Notification
	finished          chan bool
	channel           chan []Notification
}

func NewWorker(client Client, channel chan []Notification, finished chan bool) *notificationWorker {
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

	if !cmp.Equal(f.notificationData, notifications) {
		var unique []Notification
		for _, v := range notifications {
			skip := false
			for _, u := range f.notificationData {
				if cmp.Equal(v, u) {
					skip = true
					break
				}
			}
			if !skip {
				unique = append(unique, v)
			}
		}
		f.notificationData = notifications
		f.channel <- unique
	}

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
