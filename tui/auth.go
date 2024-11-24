package tui

import (
	bpclient "bitpin-tui/bitpin_client"
	"bitpin-tui/db"
	"fmt"
	"github.com/rivo/tview"
	"os"
)

func firstLogin() {
	session := db.GetSession()

	client_, err := bpclient.NewClient(bpclient.ClientOptions{
		ApiKey:       session.ApiKey,
		SecretKey:    session.SecretKey,
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
	})

	if err != nil {
		if client == nil {
			panic(fmt.Sprintf("client is nil %v", err))
		}
	}

	client = client_
	info, err := client.GetUserInfo()

	if err != nil {
		return
	}
	user = info

	db.UpsertSession(client.ApiKey, client.SecretKey, client.AccessToken, client.RefreshToken)
}

func doLogin(app *tview.Application) {
	form := tview.NewForm().
		AddPasswordField("API Key", "", 0, '*', nil).
		AddPasswordField("API Secret", "", 0, '*', nil)

	form.
		AddButton("Submit", func() {
			client, err := bpclient.NewClient(bpclient.ClientOptions{
				ApiKey:    form.GetFormItemByLabel("API Key").(*tview.InputField).GetText(),
				SecretKey: form.GetFormItemByLabel("API Secret").(*tview.InputField).GetText(),
			})

			if err != nil {
				errorModal(app, fmt.Sprintf("Failed to log in: %v", err), login)
				return
			}

			info, err := client.GetUserInfo()
			if err != nil {
				errorModal(app, fmt.Sprintf("Failed to get user info: %v", err), login)
				return
			}
			user = info

			db.UpsertSession(client.ApiKey, client.SecretKey, client.AccessToken, client.RefreshToken)

			loggedInModal(app)
		}).
		AddButton("Quit", func() {
			app.Stop()
		}).
		AddButton("Back", func() {
			mainMenu(app)
		})

	form.
		SetBorder(true).
		SetTitle("Login").
		SetTitleAlign(tview.AlignLeft).
		SetRect(0, 0, 30, 10)

	app.SetRoot(form, true)

}

func doLogout(app *tview.Application) {
	user = nil
	err := db.DeleteSession()

	if err != nil {
		modal := tview.NewModal().
			SetText(fmt.Sprintf("Failed to log out: %v", err)).
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				mainMenu(app)
			})

		app.SetRoot(modal, true).SetFocus(modal)
	}

	modal := tview.NewModal().
		SetText("Logged out").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			mainMenu(app)
		})

	app.SetRoot(modal, true).SetFocus(modal)
}

func doReset(app *tview.Application) {
	err := os.Remove("bitpin-tui.db")
	if err != nil {
		errorModal(app, fmt.Sprintf("Failed to reset: %v", err), mainMenu)
		return
	}
	errorModal(app, "Database reset successfully. Please restart the app.", func(app *tview.Application) {
		app.Stop()
	})

}

func logout(app *tview.Application) {
	confirmModal(app, "Are you sure you want to log out?", doLogout, mainMenu)
}

func login(app *tview.Application) {
	warningModal(app, "Currently credentials are stored in plain text in the database. Please be aware of this.\n\nYou can log out at any time to remove them.", doLogin)
}

func reset(app *tview.Application) {
	confirmModal(app, "Are you sure you want to reset the database?", doReset, mainMenu)
}
