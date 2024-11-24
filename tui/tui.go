package tui

import (
	bpclient "bitpin-tui/bitpin_client"
	"github.com/rivo/tview"
	"log"
)

var (
	user   *bpclient.UserInfo
	client *bpclient.Client
)

func tui() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	app := tview.NewApplication().EnableMouse(true)
	addEventHandlers(app)

	firstLogin()
	mainMenu(app)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func Run() {
	tui()
}
