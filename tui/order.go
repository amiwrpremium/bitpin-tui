package tui

import (
	bpclient "bitpin-tui/bitpin_client"
	"bitpin-tui/db"
	"bitpin-tui/utils"
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
	"sync"
	"time"
)

func createOrderStatusTable(app *tview.Application, orders []*bpclient.OrderStatus) *tview.Table {
	table := tview.NewTable().
		SetBorders(true)

	table.SetCell(0, 0, tview.NewTableCell("Id").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("Symbol").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("Type").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 3, tview.NewTableCell("Side").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 4, tview.NewTableCell("Base Amount").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 5, tview.NewTableCell("Quote Amount").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 6, tview.NewTableCell("Price").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 7, tview.NewTableCell("Identifier").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 8, tview.NewTableCell("State").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 9, tview.NewTableCell("Created At").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 10, tview.NewTableCell("Relative Time").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 11, tview.NewTableCell("Dealed Base Amount").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 12, tview.NewTableCell("Dealed Quote Amount").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 13, tview.NewTableCell("Commission").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 13, tview.NewTableCell("Req To Cancel").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))

	for i, order := range orders {
		color := tcell.ColorWhite
		if order.Side == "buy" {
			order.Side = "BUY"
			color = tcell.ColorGreen
		} else {
			order.Side = "SELL"
			color = tcell.ColorRed
		}

		if order.Identifier == "" {
			order.Identifier = "-"
		}

		table.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("%d", order.Id)).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 1, tview.NewTableCell(order.Symbol).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 2, tview.NewTableCell(order.Type).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 3, tview.NewTableCell(order.Side).SetTextColor(color).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 4, tview.NewTableCell(utils.FormatWithCommas(order.BaseAmount)).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 5, tview.NewTableCell(utils.FormatWithCommas(order.QuoteAmount)).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 6, tview.NewTableCell(utils.FormatWithCommas(order.Price)).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 7, tview.NewTableCell(order.Identifier).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 8, tview.NewTableCell(order.State).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 9, tview.NewTableCell(order.CreatedAt.String()).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 10, tview.NewTableCell(utils.FormatWithCommas(strconv.FormatInt(int64(time.Since(order.CreatedAt).Seconds()), 10))+"ms ago").SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 11, tview.NewTableCell(order.DealedBaseAmount).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 12, tview.NewTableCell(order.DealedQuoteAmount).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 13, tview.NewTableCell(order.Commission).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 13, tview.NewTableCell(strconv.FormatBool(order.ReqToCancel)).SetAlign(tview.AlignCenter))
	}

	table.Select(0, 0).SetFixed(1, 1)

	return table

}

func openOrders(app *tview.Application) {
	form := tview.NewForm().
		AddInputField("Symbol (optional)", "", 0, nil, nil).
		AddDropDown("Side (optional)", []string{"", "buy", "sell"}, 0, nil)

	form.AddButton("Submit", func() {
		symbol := form.GetFormItemByLabel("Symbol (optional)").(*tview.InputField).GetText()
		_, side := form.GetFormItemByLabel("Side (optional)").(*tview.DropDown).GetCurrentOption()

		orders, err := client.GetOpenOrders(bpclient.GetOrdersParams{Symbol: symbol, Side: side})

		if err != nil {
			errorModal(app, fmt.Sprintf("Failed to get open orders: %v", err), openOrders)
			return
		}

		if len(orders) == 0 {
			errorModal(app, "No open orders found", openOrders)
			return
		}

		table := createOrderStatusTable(app, orders)
		app.SetRoot(table, true)
	})

	form.
		AddButton("Back", func() {
			mainMenu(app)
		})

	form.
		SetBorder(true).
		SetTitle("Open Orders").
		SetTitleAlign(tview.AlignLeft).
		SetRect(0, 0, 30, 10)

	app.SetRoot(form, true)
}

func cancelOrder(app *tview.Application) {
	form := tview.NewForm().
		AddInputField("Order Id", "", 0, nil, nil)

	form.AddButton("Submit", func() {
		confirmModal(app, "Are you sure you want to cancel this order?", func(application *tview.Application) {
			orderId, _ := strconv.Atoi(form.GetFormItemByLabel("Order Id").(*tview.InputField).GetText())
			err := client.CancelOrder(orderId)
			if err != nil {
				errorModal(app, fmt.Sprintf("Failed to cancel order: %v", err), cancelOrder)
				return
			}

			modal := tview.NewModal().
				SetText("Order canceled").
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					mainMenu(app)
				})

			app.SetRoot(modal, true).SetFocus(modal)
		}, mainMenu)
	})

	form.
		AddButton("Back", func() {
			mainMenu(app)
		})
}

func createOrder(app *tview.Application) {
	form := tview.NewForm().
		AddInputField("Symbol", "", 0, nil, nil).
		AddDropDown("Type", []string{"limit", "market"}, 0, nil).
		AddDropDown("Side", []string{"buy", "sell"}, 0, nil).
		AddInputField("Price", "", 0, nil, nil).
		AddInputField("Base Amount", "", 0, nil, nil)

	form.AddButton("Submit", func() {
		symbol := form.GetFormItemByLabel("Symbol").(*tview.InputField).GetText()
		_, orderType := form.GetFormItemByLabel("Type").(*tview.DropDown).GetCurrentOption()
		_, side := form.GetFormItemByLabel("Side").(*tview.DropDown).GetCurrentOption()
		price := form.GetFormItemByLabel("Price").(*tview.InputField).GetText()
		baseAmount := form.GetFormItemByLabel("Base Amount").(*tview.InputField).GetText()

		if symbol == "" {
			errorModal(app, "Symbol is required", createOrder)
			return
		}
		if orderType == "" {
			errorModal(app, "Type is required", createOrder)
			return
		}
		if side == "" {
			errorModal(app, "Side is required", createOrder)
			return
		}
		if orderType == "limit" && price == "" {
			errorModal(app, "Price is required", createOrder)
			return
		}
		if baseAmount == "" {
			errorModal(app, "Base Amount is required", createOrder)
			return
		}

		params := bpclient.CreateOrderParams{
			Symbol:     symbol,
			Type:       orderType,
			Side:       side,
			Price:      price,
			BaseAmount: baseAmount,
		}

		order, err := client.CreateOrder(params)
		if err != nil {
			errorModal(app, fmt.Sprintf("Failed to create order: %v", err), createOrder)
			return
		}

		modal := tview.NewModal().
			SetText("Order created").
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				table := createOrderStatusTable(app, []*bpclient.OrderStatus{order})
				app.SetRoot(table, true)
			})

		app.SetRoot(modal, true).SetFocus(modal)
	})

	form.
		AddButton("Back", func() {
			mainMenu(app)
		},
		)

	form.
		SetBorder(true).
		SetTitle("Create Order").
		SetTitleAlign(tview.AlignLeft).
		SetRect(0, 0, 30, 10)

	app.SetRoot(form, true)
}

func pussyOut(app *tview.Application) {
	confirmModal(app, "Are you sure?", doPussyOut, mainMenu)
}

func prependLogMessage(tv *tview.TextView, message string) {
	// Use a mutex if multiple goroutines access this to prevent race conditions
	currentText := tv.GetText(false)
	newText := message + "\n" + currentText
	tv.SetText(newText)
}

func doPussyOut(app *tview.Application) {
	cancelingChannelText := newPrimitive(app, "", tcell.ColorWhite)
	cancelResultChannelText := newPrimitive(app, "", tcell.ColorWhite)
	openOrdersText := newPrimitive(app, "", tcell.ColorWhite)

	grid := tview.NewGrid().
		SetRows(-1).         // Changed from (0, 0) to (-1) to use full height
		SetColumns(0, 0, 0). // Changed to include three columns
		SetBorders(true)

	grid.AddItem(openOrdersText, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(cancelingChannelText, 0, 1, 1, 1, 0, 0, false)
	grid.AddItem(cancelResultChannelText, 0, 2, 1, 1, 0, 0, false)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get the number of workers from the settings
	workerCount, err := strconv.Atoi(db.GetSetting("pussy_out_workers"))
	if err != nil {
		errorModal(app, fmt.Sprintf("Failed to get pussy out workers: %v", err), mainMenu)
		return
	}

	// Buffered channel to avoid blocking
	bufferSize, err := strconv.Atoi(db.GetSetting("pussy_out_buffer_size"))
	if err != nil {
		errorModal(app, fmt.Sprintf("Failed to get pussy out buffer size: %v", err), mainMenu)
		return
	}
	cancelChannel := make(chan int, bufferSize)

	var wg sync.WaitGroup

	// Cancel order workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case orderId, ok := <-cancelChannel:
					if !ok {
						return
					}
					message := fmt.Sprintf("Worker %d: Canceling order: %d", workerID, orderId)
					prependLogMessage(cancelingChannelText, message)

					err := client.CancelOrder(orderId)
					if err != nil {
						errorMessage := fmt.Sprintf("Worker %d: Error canceling order: %d: %v", workerID, orderId, err)
						prependLogMessage(cancelResultChannelText, errorMessage)
					} else {
						successMessage := fmt.Sprintf("Worker %d: Successfully canceled order: %d", workerID, orderId)
						prependLogMessage(cancelResultChannelText, successMessage)
					}
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}

	// Order fetching workers - running at max speed
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					prependLogMessage(openOrdersText, fmt.Sprintf("Worker %d: Fetching open orders", workerID))
					orders, err := client.GetOpenOrders(bpclient.GetOrdersParams{State: "active"})
					if err != nil {
						errorMessage := fmt.Sprintf("Worker %d: Error fetching open orders: %v", workerID, err)
						prependLogMessage(openOrdersText, errorMessage)
						continue
					}

					message := fmt.Sprintf("Worker %d: %d open orders fetched", workerID, len(orders))
					prependLogMessage(openOrdersText, message)

					for _, order := range orders {
						select {
						case cancelChannel <- order.Id:
							// Sent to channel
						case <-ctx.Done():
							return
						default:
							// Channel full, skip and continue
							prependLogMessage(openOrdersText, fmt.Sprintf("Worker %d: Skipped order %d (channel full)", workerID, order.Id))
						}
					}
				}
			}
		}(i)
	}

	// Handle application cleanup
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			cancel()
			close(cancelChannel)
			app.Stop()
			return nil
		}
		return event
	})

	if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		cancel()
		close(cancelChannel)
		panic(err)
	}

	wg.Wait()
}
