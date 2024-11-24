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

func info(app *tview.Application) {
	table := tview.NewTable().
		SetBorders(true)

	table.SetCell(0, 0, tview.NewTableCell("User Identifier").SetTextColor(tcell.ColorYellow))
	table.SetCell(0, 1, tview.NewTableCell(user.UserIdentifier).SetTextColor(tcell.ColorGreen))
	table.SetCell(1, 0, tview.NewTableCell("Phone").SetTextColor(tcell.ColorYellow))
	table.SetCell(1, 1, tview.NewTableCell(user.Phone).SetTextColor(tcell.ColorGreen))
	table.SetCell(2, 0, tview.NewTableCell("State").SetTextColor(tcell.ColorYellow))
	table.SetCell(2, 1, tview.NewTableCell(user.State).SetTextColor(tcell.ColorGreen))
	table.SetCell(3, 0, tview.NewTableCell("Is Phone Confirmed").SetTextColor(tcell.ColorYellow))
	table.SetCell(3, 1, tview.NewTableCell(strconv.FormatBool(user.IsPhoneConfirmed)).SetTextColor(tcell.ColorGreen))
	table.SetCell(4, 0, tview.NewTableCell("Is Email Confirmed").SetTextColor(tcell.ColorYellow))
	table.SetCell(4, 1, tview.NewTableCell(strconv.FormatBool(user.IsEmailConfirmed)).SetTextColor(tcell.ColorGreen))
	table.SetCell(5, 0, tview.NewTableCell("First Name").SetTextColor(tcell.ColorYellow))
	table.SetCell(5, 1, tview.NewTableCell(user.FirstName).SetTextColor(tcell.ColorGreen))
	table.SetCell(6, 0, tview.NewTableCell("Last Name").SetTextColor(tcell.ColorYellow))
	table.SetCell(6, 1, tview.NewTableCell(user.LastName).SetTextColor(tcell.ColorGreen))
	table.SetCell(7, 0, tview.NewTableCell("Full Name").SetTextColor(tcell.ColorYellow))
	table.SetCell(7, 1, tview.NewTableCell(user.Fullname).SetTextColor(tcell.ColorGreen))
	table.SetCell(8, 0, tview.NewTableCell("Birth Date Text").SetTextColor(tcell.ColorYellow))
	table.SetCell(8, 1, tview.NewTableCell(user.BirthDateText).SetTextColor(tcell.ColorGreen))
	table.SetCell(9, 0, tview.NewTableCell("Email").SetTextColor(tcell.ColorYellow))
	table.SetCell(9, 1, tview.NewTableCell(user.Email).SetTextColor(tcell.ColorGreen))
	table.SetCell(10, 0, tview.NewTableCell("Type").SetTextColor(tcell.ColorYellow))
	table.SetCell(10, 1, tview.NewTableCell(user.Type).SetTextColor(tcell.ColorGreen))
	table.SetCell(11, 0, tview.NewTableCell("Remaining Daily Withdraw").SetTextColor(tcell.ColorYellow))
	table.SetCell(11, 1, tview.NewTableCell(utils.FormatWithCommas(strconv.FormatInt(user.RemainingDailyWithdraw, 10))).SetTextColor(tcell.ColorGreen))
	table.SetCell(12, 0, tview.NewTableCell("Remaining Monthly Withdraw").SetTextColor(tcell.ColorYellow))
	table.SetCell(12, 1, tview.NewTableCell(utils.FormatWithCommas(strconv.FormatInt(user.RemainingMonthlyWithdraw, 10))).SetTextColor(tcell.ColorGreen))
	table.SetCell(13, 0, tview.NewTableCell("Tetherban").SetTextColor(tcell.ColorYellow))
	table.SetCell(13, 1, tview.NewTableCell(strconv.FormatBool(user.Tetherban)).SetTextColor(tcell.ColorGreen))

	app.SetRoot(table, true)
}
