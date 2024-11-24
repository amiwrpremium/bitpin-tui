package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func addEventHandlers(app *tview.Application) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			app.Stop()
		}

		if event.Key() == tcell.KeyESC {
			mainMenu(app)
		}

		return event
	})
}

func newPrimitive(app *tview.Application, text string, color tcell.Color) *tview.TextView {
	return tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetText(text).SetTextAlign(tview.AlignCenter).
		SetTextColor(color).
		SetChangedFunc(func() {
			app.Draw()
		})
}
