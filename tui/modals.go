package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func errorModal(app *tview.Application, msg string, goTo func(app *tview.Application)) {
	modal := tview.NewModal().
		SetBackgroundColor(tcell.ColorRed).
		SetText(msg).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if goTo != nil {
				goTo(app)
			} else {
				return
			}
		})

	app.SetRoot(modal, true).SetFocus(modal)
}

func confirmModal(app *tview.Application, text string, yesCallback func(app *tview.Application), noCallback func(app *tview.Application)) {
	modal := tview.NewModal().
		SetText(text).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				yesCallback(app)
			} else {
				noCallback(app)
			}
		})

	app.SetRoot(modal, true).SetFocus(modal)
}

func warningModal(app *tview.Application, text string, goTo func(app *tview.Application)) {
	modal := tview.NewModal().
		SetBackgroundColor(tcell.ColorYellow).
		SetText(text).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if goTo != nil {
				goTo(app)
			} else {
				return
			}
		})

	app.SetRoot(modal, true).SetFocus(modal)
}

func loggedInModal(app *tview.Application) {
	modal := tview.NewModal().
		SetBackgroundColor(tcell.ColorGreen).
		SetText(fmt.Sprintf("Logged in as %s", user.Fullname)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			mainMenu(app)
		})

	app.SetRoot(modal, true).SetFocus(modal)
}
