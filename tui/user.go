package tui

import (
	bpclient "bitpin-tui/bitpin_client"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

func balances(app *tview.Application) {
	form := tview.NewForm()
	form.AddCheckbox("Non Zero Balances", false, nil).
		AddCheckbox("Live", false, nil)

	form.AddButton("Submit", func() {
		live := form.GetFormItemByLabel("Live").(*tview.Checkbox).IsChecked()
		//nonZero := form.GetFormItemByLabel("Non Zero Balances").(*tview.Checkbox).IsChecked()

		if live {
			table := tview.NewTable().
				SetBorders(true)

			table.SetCell(0, 0, tview.NewTableCell("Asset").SetTextColor(tcell.ColorYellow))
			table.SetCell(0, 1, tview.NewTableCell("Balance").SetTextColor(tcell.ColorYellow))
			table.SetCell(0, 2, tview.NewTableCell("Frozen").SetTextColor(tcell.ColorYellow))

			for {
				blnc, err := client.GetBalances(bpclient.GetBalancesParams{Limit: 200})
				if err != nil {
					errorModal(app, fmt.Sprintf("Failed to get balances: %v", err), balances)
					return
				}

				for i, balance := range blnc {
					table.SetCell(i+1, 0, tview.NewTableCell(balance.Asset).SetTextColor(tcell.ColorGreen))
					table.SetCell(i+1, 1, tview.NewTableCell(balance.Balance))
					table.SetCell(i+1, 2, tview.NewTableCell(balance.Frozen))
				}

				app.SetRoot(table, true)
				time.Sleep(5 * time.Second)
			}

		} else {
			blnc, err := client.GetBalances(bpclient.GetBalancesParams{Limit: 200})
			if err != nil {
				errorModal(app, fmt.Sprintf("Faileds to get balances: %v", err), balances)
				return
			}

			table := tview.NewTable().
				SetBorders(true)

			table.SetCell(0, 0, tview.NewTableCell("Asset").SetTextColor(tcell.ColorYellow))
			table.SetCell(0, 1, tview.NewTableCell("Balance").SetTextColor(tcell.ColorYellow))
			table.SetCell(0, 2, tview.NewTableCell("Frozen").SetTextColor(tcell.ColorYellow))

			for i, balance := range blnc {
				table.SetCell(i+1, 0, tview.NewTableCell(balance.Asset).SetTextColor(tcell.ColorGreen))
				table.SetCell(i+1, 1, tview.NewTableCell(balance.Balance))
				table.SetCell(i+1, 2, tview.NewTableCell(balance.Frozen))
			}

			app.SetRoot(table, true)
		}
	})

	form.
		AddButton("Back", func() {
			mainMenu(app)
		})

	form.
		SetBorder(true).
		SetTitle("Balances").
		SetTitleAlign(tview.AlignLeft).
		SetRect(0, 0, 30, 10)

	app.SetRoot(form, true)
}
