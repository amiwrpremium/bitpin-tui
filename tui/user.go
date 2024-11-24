package tui

import (
	bpclient "bitpin-tui/bitpin_client"
	"bitpin-tui/utils"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
	"strings"
	"time"
)

func createBalanceTable(balances []bpclient.Balance, excludeZero bool) *tview.Table {
	table := tview.NewTable().
		SetBorders(true)

	table.SetCell(0, 0, tview.NewTableCell("Asset").SetTextColor(tcell.ColorYellow))
	table.SetCell(0, 1, tview.NewTableCell("Balance").SetTextColor(tcell.ColorYellow))
	table.SetCell(0, 2, tview.NewTableCell("Frozen").SetTextColor(tcell.ColorYellow))

	iter := 0

	for _, balance := range balances {
		floatBalance, err := strconv.ParseFloat(balance.Balance, 64)
		if err != nil {
			errorModal(nil, fmt.Sprintf("Failed to parse balance: %v", err), nil)
			return nil
		}
		if excludeZero && floatBalance == 0 {
			continue
		}

		table.SetCell(iter+1, 0, tview.NewTableCell(balance.Asset).SetTextColor(tcell.ColorGreen))
		table.SetCell(iter+1, 1, tview.NewTableCell(utils.FormatWithCommas(balance.Balance)))
		table.SetCell(iter+1, 2, tview.NewTableCell(utils.FormatWithCommas(balance.Frozen)))

		iter++
	}

	return table
}

func getBalanceTable(app *tview.Application, excludeZero bool, assets []string) *tview.Table {
	blnc, err := client.GetBalances(bpclient.GetBalancesParams{Limit: 1000, Assets: assets})
	if err != nil {
		errorModal(app, fmt.Sprintf("Failed to get balances: %v", err), balances)
		return nil
	}

	return createBalanceTable(blnc, excludeZero)
}

func balances(app *tview.Application) {
	form := tview.NewForm()
	form.AddCheckbox("Exclude Zero Balances", true, nil).
		AddInputField("Interval", "-1", 0, nil, nil).
		AddInputField("Assets", "", 0, nil, nil).
		AddButton("Submit", func() {
			excludeNonZero := form.GetFormItemByLabel("Exclude Zero Balances").(*tview.Checkbox).IsChecked()
			interval, _ := strconv.Atoi(form.GetFormItemByLabel("Interval").(*tview.InputField).GetText())
			assets := form.GetFormItemByLabel("Assets").(*tview.InputField).GetText()

			if interval > 0 {
				for {
					app.SetRoot(tview.NewFlex(), true)
					table := getBalanceTable(app, excludeNonZero, strings.Split(assets, ","))
					app.SetRoot(table, true)
					time.Sleep(time.Duration(interval) * time.Second)

					app.Draw()
				}

			} else {
				table := getBalanceTable(app, excludeNonZero, strings.Split(assets, ","))
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
