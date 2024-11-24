package tui

import (
	bpclient "bitpin-tui/bitpin_client"
	"bitpin-tui/db"
	"bitpin-tui/utils"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"math"
	"strconv"
	"time"
)

func orderBook(app *tview.Application) {
	form := tview.NewForm()
	favSymbols := db.GetFavorites("order_book")

	if len(favSymbols) > 0 {
		favSymbols = append(favSymbols, "Other")
		form.AddDropDown("Favorite", favSymbols, 0, func(option string, optionIndex int) {
			if option == "Other" {
				//form.RemoveFormItem(0)
				form.Clear(false)
				form.AddInputField("Symbol", "", 0, nil, nil)
				form.AddInputField("Depth", "20", 0, nil, nil)
				form.AddInputField("Interval", "1", 0, nil, nil)
				favSymbols = []string{}
			}
		})
	} else {
		form.AddInputField("Symbol", "", 0, nil, nil)
	}

	form.AddInputField("Depth", "20", 0, nil, nil)
	form.AddInputField("Interval", "1", 0, nil, nil)

	form.AddButton("Submit", func() {
		symbol := ""
		if len(favSymbols) == 0 {
			symbol = form.GetFormItemByLabel("Symbol").(*tview.InputField).GetText()
		} else {
			_, symbol = form.GetFormItemByLabel("Favorite").(*tview.DropDown).GetCurrentOption()
		}

		if symbol == "" {
			errorModal(app, "Symbol is required", orderBook)
			return
		}
		if !utils.StringEndsWith(symbol, "_USDT") && !utils.StringEndsWith(symbol, "_IRT") {
			errorModal(app, "Invalid symbol. Symbol must end with _USDT or _IRT", orderBook)
			return
		}

		db.UpsertFavorite("order_book", symbol)

		depth, _ := strconv.Atoi(form.GetFormItemByLabel("Depth").(*tview.InputField).GetText())
		interval, _ := strconv.Atoi(form.GetFormItemByLabel("Interval").(*tview.InputField).GetText())

		// two text views for the order book (bids and asks)
		bidsTextView := newPrimitive(app, "", tcell.ColorGreen)
		asksTextView := newPrimitive(app, "", tcell.ColorRed)
		lastUpdatedTextView := newPrimitive(app, "Last updated at: ", tcell.ColorOrange)

		// two loggers for the order book (bids and asks)
		bidsLogger := log.New(bidsTextView, "", 0)
		asksLogger := log.New(asksTextView, "", 0)

		// create a grid layout
		grid := tview.NewGrid().
			SetRows(1, 1, 1, 0, 1).
			SetColumns(0, 0).
			SetBorders(true).
			AddItem(newPrimitive(app, symbol+" Order Book", tcell.ColorPurple), 0, 0, 1, 2, 0, 0, false).
			AddItem(lastUpdatedTextView, 4, 0, 1, 2, 0, 0, false)

		// add "ask" and "bid" as a first row in the grid
		grid.AddItem(newPrimitive(app, "Asks", tcell.ColorRed), 1, 0, 1, 1, 0, 0, false)
		grid.AddItem(newPrimitive(app, "Bids", tcell.ColorGreen), 1, 1, 1, 1, 0, 0, false)

		grid.AddItem(newPrimitive(app, "Price \t\t\t Amount", tcell.ColorWhite), 2, 0, 1, 1, 0, 0, false)
		grid.AddItem(newPrimitive(app, "Price \t\t\t Amount", tcell.ColorWhite), 2, 1, 1, 1, 0, 0, false)

		grid.AddItem(asksTextView, 3, 0, 1, 1, 0, 0, false)
		grid.AddItem(bidsTextView, 3, 1, 1, 1, 0, 0, false)

		// create a channel for the order book
		orderBookChannel := make(chan *bpclient.OrderBook)

		lastUpdatedAt := time.Now()

		// create a go routine that listens to the order book channel and updates the text views
		go func() {
			for {
				orderBook := <-orderBookChannel

				bidsTextView.Clear()
				asksTextView.Clear()
				lastUpdatedTextView.Clear()

				minDepth := int(math.Min(math.Min(float64(len(orderBook.Bids)), float64(len(orderBook.Asks))), float64(depth)))

				for i := 0; i < minDepth; i++ {
					bid := orderBook.Bids[i]
					ask := orderBook.Asks[i]

					_, _ = bidsTextView.Write([]byte(fmt.Sprintf("%s \t\t\t %s\n", utils.FormatWithCommas(bid[0]), utils.FormatWithCommas(bid[1]))))
					_, _ = asksTextView.Write([]byte(fmt.Sprintf("%s \t\t\t %s\n", utils.FormatWithCommas(ask[0]), utils.FormatWithCommas(ask[1]))))
				}

				_, _ = lastUpdatedTextView.Write([]byte(fmt.Sprintf("Last updated: %dms ago\n", time.Since(lastUpdatedAt).Milliseconds())))

				lastUpdatedAt = time.Now()
				app.Draw()
			}
		}()

		// create a go routine that retrieves the order book in a loop and sends it to the order book channel
		go func() {
			for {
				orderBook, err := client.GetOrderBook(symbol)
				if err != nil {
					bidsLogger.Printf("error fetching order book: %v\n", err)
					asksLogger.Printf("error fetching order book: %v\n", err)
					continue
				}
				orderBookChannel <- orderBook

				time.Sleep(time.Duration(interval) * time.Second)
			}
		}()

		if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
			panic(err)
		}
	})

	form.
		AddButton("Back", func() {
			mainMenu(app)
		})

	form.
		SetBorder(true).
		SetTitle("Order Book").
		SetTitleAlign(tview.AlignLeft).
		SetRect(0, 0, 30, 10)

	app.SetRoot(form, true)
}

func tickers(app *tview.Application) {
	tickers, err := client.GetTickers()
	if err != nil {
		errorModal(app, fmt.Sprintf("Failed to get tickers: %v", err), mainMenu)
		return
	}

	table := tview.NewTable().
		SetBorders(true)

	table.SetCell(0, 0, tview.NewTableCell("Symbol").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("Price").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("Daily Change Price").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 3, tview.NewTableCell("Low").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 4, tview.NewTableCell("High").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 5, tview.NewTableCell("Updated At").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))

	for i, ticker := range tickers {
		table.SetCell(i+1, 0, tview.NewTableCell(ticker.Symbol).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 1, tview.NewTableCell(utils.FormatWithCommas(ticker.Price)).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 2, tview.NewTableCell(utils.FormatWithCommas(fmt.Sprintf("%.2f", ticker.DailyChangePrice))).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 3, tview.NewTableCell(utils.FormatWithCommas(ticker.Low)).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 4, tview.NewTableCell(utils.FormatWithCommas(ticker.High)).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 5, tview.NewTableCell(utils.FormatWithCommas(strconv.FormatFloat(time.Since(time.Unix(int64(ticker.Timestamp), 0)).Seconds(), 'f', 0, 64))+"s ago").SetAlign(tview.AlignCenter))

	}

	table.Select(0, 0).SetFixed(1, 1)
	app.SetRoot(table, true)
}

func recentTrades(app *tview.Application) {
	form := tview.NewForm()
	favSymbols := db.GetFavorites("recent_trades")

	if len(favSymbols) > 0 {
		favSymbols = append(favSymbols, "Other")
		form.AddDropDown("Favorite", favSymbols, 0, func(option string, optionIndex int) {
			if option == "Other" {
				//form.RemoveFormItem(0)
				form.Clear(false)
				form.AddInputField("Symbol", "", 0, nil, nil)
				favSymbols = []string{}
			}
		})
	} else {
		form.AddInputField("Symbol", "", 0, nil, nil)
	}

	form.AddButton("Submit", func() {
		symbol := ""

		if len(favSymbols) == 0 {
			symbol = form.GetFormItemByLabel("Symbol").(*tview.InputField).GetText()
		} else {
			_, symbol = form.GetFormItemByLabel("Favorite").(*tview.DropDown).GetCurrentOption()
		}

		if symbol == "" {
			errorModal(app, "Symbol is required", recentTrades)
			return
		}
		if !utils.StringEndsWith(symbol, "_USDT") && !utils.StringEndsWith(symbol, "_IRT") {
			errorModal(app, "Invalid symbol. Symbol must end with _USDT or _IRT", recentTrades)
			return
		}

		db.UpsertFavorite("recent_trades", symbol)

		trades, err := client.GetRecentTrades(symbol)
		if err != nil {
			errorModal(app, fmt.Sprintf("Failed to get recent trades: %v", err), recentTrades)
			return
		}

		table := tview.NewTable().
			SetBorders(true)

		table.SetCell(0, 0, tview.NewTableCell("Id").SetTextColor(tcell.ColorYellow))
		table.SetCell(0, 1, tview.NewTableCell("Price").SetTextColor(tcell.ColorYellow))
		table.SetCell(0, 2, tview.NewTableCell("Base Amount").SetTextColor(tcell.ColorYellow))
		table.SetCell(0, 3, tview.NewTableCell("Quote Amount").SetTextColor(tcell.ColorYellow))
		table.SetCell(0, 4, tview.NewTableCell("Side").SetTextColor(tcell.ColorYellow))

		for i, trade := range trades {
			color := tcell.ColorWhite
			if trade.Side == "buy" {
				trade.Side = "BUY"
				color = tcell.ColorGreen
			} else {
				trade.Side = "SELL"
				color = tcell.ColorRed
			}

			table.SetCell(i+1, 0, tview.NewTableCell(trade.Id))
			table.SetCell(i+1, 1, tview.NewTableCell(trade.Price))
			table.SetCell(i+1, 2, tview.NewTableCell(trade.BaseAmount))
			table.SetCell(i+1, 3, tview.NewTableCell(trade.QuoteAmount))
			table.SetCell(i+1, 4, tview.NewTableCell(trade.Side).SetTextColor(color))
		}

		table.Select(0, 0).SetFixed(1, 1)

		app.SetRoot(table, true)
	})

	form.
		AddButton("Back", func() {
			mainMenu(app)
		})

	form.
		SetBorder(true).
		SetTitle("Recent Trades").
		SetTitleAlign(tview.AlignLeft).
		SetRect(0, 0, 30, 10)

	app.SetRoot(form, true)
}
