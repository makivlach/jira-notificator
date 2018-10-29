// 12 august 2018

package main

import (
	"C"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/vlachmilan/jira-notificator/pkg/jira"
	"log"
	_ "log"
	"os"
	_ "os"
)

const (
	notificationSound = "assets/notify.mp3"
	icon              = "assets/icon.ico"
	interval          = 15
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
		Text:     "Provedení katualizace",
	}
}

func makeBasicControlsPage() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	group := ui.NewGroup("Přihlašovací údaje")
	group.SetMargined(true)
	vbox.Append(group, true)

	group.SetChild(ui.NewNonWrappingMultilineEntry())

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)

	host := ui.NewEntry()
	username := ui.NewEntry()
	password := ui.NewPasswordEntry()

	button := ui.NewButton("Přihlásit se")
	button.OnClicked(func(*ui.Button) {
		host.Disable()
		username.Disable()
		password.Disable()
		button.Disable()

		go func(window *ui.Window) {
			client, err := jira.NewClient(host.Text(), username.Text(), password.Text())
			if err != nil {
				ui.QueueMain(func() {
					ui.MsgBoxError(window, "Chyba", err.Error())
				})
				return
			}

			err = client.Login()
			if err != nil {
				ui.QueueMain(func() {
					ui.MsgBoxError(window, "Chyba", err.Error())
				})
				return
			}
			ui.QueueMain(func() {
				ui.MsgBox(window, "Úspěch", "Úspěšně jste se přihlásili!")
			})

			err = jira.FetchNewNotifications(client, notificationData)
			if err != nil {
				ui.QueueMain(func() {
					ui.MsgBoxError(window, "Chyba", err.Error())
				})
				return
			}
		}(mainwin)
	})

	entryForm.Append("Host", host, false)
	entryForm.Append("Uživatelské jméno", username, false)
	entryForm.Append("Heslo", password, false)
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
