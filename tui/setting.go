package tui

import (
	"bitpin-tui/db"
	"github.com/rivo/tview"
)

func setting(app *tview.Application) {
	form := tview.NewForm()

	allSettings := db.GetALlSettings()

	for key, value := range allSettings {
		form.AddInputField(key, value, 0, nil, nil)
	}

	form.AddButton("Submit", func() {
		for key, _ := range allSettings {
			db.UpsertSetting(key, form.GetFormItemByLabel(key).(*tview.InputField).GetText())
		}

		mainMenu(app)
	})

	form.AddButton("Back", func() {
		mainMenu(app)
	})

	form.
		SetBorder(true).
		SetTitle("Settings").
		SetTitleAlign(tview.AlignLeft).
		SetRect(0, 0, 30, 10)

	app.SetRoot(form, true)
}
