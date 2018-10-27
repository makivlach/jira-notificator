package jira

type alerter func(title, message, appIcon string) error

func NewNotificator(alertFunc alerter) Notificator {
	return Notificator{alertFunc}
}

type Notificator struct {
	alert alerter
}

func (n Notificator) Notify(notifications *Notifications) error {
	for _, notification := range notifications.Notifications {
		err := n.alert("Jira", notification.Title, "assets/information.png")

		if err != nil {
			return err
		}
	}
	return nil
}
