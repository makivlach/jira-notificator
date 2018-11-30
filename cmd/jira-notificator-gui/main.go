package main

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/vlachmilan/jira-notificator/internal/jira"
	"log"
	"os"
)

const (
	notificationSound = "assets/notify.wav"
	icon              = "assets/icon.ico"
	interval          = 25
)

var mainwin *ui.Window
var client jira.Client
var notificationData jira.NotificationData

func init() {
	sound, err := Asset(notificationSound)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	notificationData = jira.NotificationData{
		Sound:    sound,
		Interval: interval,
		Text:     "Perform update",
	}
}

func makeBasicControlsPage() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	group := ui.NewGroup("User credentials")
	group.SetMargined(true)
	vbox.Append(group, true)

	group.SetChild(ui.NewNonWrappingMultilineEntry())

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)

	host := ui.NewEntry()
	username := ui.NewEntry()
	password := ui.NewPasswordEntry()

	button := ui.NewButton("Log in")
	button.OnClicked(func(*ui.Button) {
		if host.Text() == "" || username.Text() == "" || password.Text() == "" {
			ui.MsgBoxError(mainwin, "Error", "fill all required fields")
			return
		}

		host.Disable()
		username.Disable()
		password.Disable()
		button.Disable()

		go func(window *ui.Window) {
			client, err := jira.NewClient(host.Text(), username.Text(), password.Text())
			if err != nil {
				ui.QueueMain(func() {
					ui.MsgBoxError(window, "Error", err.Error())
					host.Enable()
					username.Enable()
					password.Enable()
					button.Enable()
				})
				return
			}

			err = client.Login()
			if err != nil {
				ui.QueueMain(func() {
					ui.MsgBoxError(window, "Error", err.Error())
					host.Enable()
					username.Enable()
					password.Enable()
					button.Enable()
				})
				return
			}
			ui.QueueMain(func() {
				ui.MsgBox(window, "Success", "User has been logged in!")
			})

			err = jira.FetchNewNotifications(client, notificationData)
			if err != nil {
				ui.QueueMain(func() {
					ui.MsgBoxError(window, "Error", err.Error())
					host.Enable()
					username.Enable()
					password.Enable()
					button.Enable()
				})
				return
			}
		}(mainwin)
	})

	entryForm.Append("Host", host, false)
	entryForm.Append("Login", username, false)
	entryForm.Append("Password", password, false)
	entryForm.Append("", button, false)

	return vbox
}

func setupUI() {
	mainwin = ui.NewWindow("Jira notifications", 400, 200, false)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return false
	})

	mainwin.SetChild(makeBasicControlsPage())
	mainwin.SetMargined(true)
	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}
