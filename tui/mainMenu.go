package tui

import (
	"github.com/rivo/tview"
)

func mainMenu(app *tview.Application) {
	form := tview.NewForm()

	form.AddButton("Get Ticker", func() {
		tickers(app)
	}).AddButton("Get Recent Trades", func() {
		recentTrades(app)
	})

	if user == nil {
		form.AddButton("Get Order Book", func() {
			orderBook(app)
		})
	}

	if user != nil {
		form.AddButton("Get Order Book", func() {
			orderBookV2(app)
		}).AddButton("Get Balances", func() {
			balances(app)
		}).AddButton("Get Open Orders", func() {
			openOrders(app)
		}).AddButton("Create Order", func() {
			createOrder(app)
		}).AddButton("Cancel Order", func() {
			cancelOrder(app)
		}).AddButton("Pussy Out", func() {
			pussyOut(app)
		}).AddButton("Info", func() {
			info(app)
		}).AddButton("Refresh", func() {
			firstLogin()
		}).AddButton("Logout", func() {
			logout(app)
		})
	} else {
		form.AddButton("Login", func() {
			login(app)
		})
	}

	form.AddButton("Reset", func() {
		reset(app)
	}).AddButton("Setting", func() {
		setting(app)
	}).AddButton("Quit", func() {
		app.Stop()
	})

	title := "Bitpin TUI"

	if user != nil {
		title += " - " + user.Fullname
	}

	form.
		SetBorder(true).
		SetTitle(title).
		SetTitleAlign(tview.AlignLeft).
		SetRect(0, 0, 30, 10)

	app.SetRoot(form, true)
}
