package main

import (
	"log"
	"github.com/mattn/go-runewidth"
    "github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/midry3125/cliRadio/radiko"
)

type Widgets struct {
	App            *tview.Application
	StationList    *tview.List
	InfoView       *tview.TextView
	nowPlaying     bool
	currentPlayer  radiko.Stream
}

func (w *Widgets) SetStationData(stations []radiko.StationInfo) {
	for _, s := range stations {
	    w.StationList.AddItem(s.Name, "", 0, func() {
			station := stations[w.StationList.GetCurrentItem()]
			if w.nowPlaying {
				w.currentPlayer.Stop()
			}
			w.currentPlayer = radiko.Stream{
				ID:    station.ID,
				Token: radiko.Auth(),
			}
			go w.currentPlayer.StartStream(w.App, w.InfoView)
			w.nowPlaying = true
			w.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyCtrlC {
					w.currentPlayer.Stop()
				}
				return event
			})
		})
	}
}

func createApp() Widgets {
	app := tview.NewApplication()
    pages := tview.NewPages()
    stationList := tview.NewList()
	stationList.SetMainTextColor(tcell.GetColor("#00FFFF")).
	            SetSelectedBackgroundColor(tcell.GetColor("#32cd32")).
	            SetSelectedFocusOnly(true)
	infoView := tview.NewTextView()
	infoView.SetBorder(true).
	         SetTitle("Log")
	infoView.SetMaxLines(30).
	         SetWrap(true)
	grid := tview.NewGrid()
    grid.AddItem(stationList, 0, 0, 1, 1, 0, 0, true).
	     AddItem(infoView, 1, 0, 1, 1, 0, 0, false)
    pages.AddPage("main", grid, true, true)
    app.SetRoot(pages, true)
    return Widgets{
		App:            app,
		StationList:    stationList,
		InfoView:       infoView,
	}
}

func main() {
	widget := createApp()
	stations, _ := radiko.ParseStationIDPage()
	widget.SetStationData(stations)
	err := widget.App.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
    runewidth.DefaultCondition = &runewidth.Condition{EastAsianWidth: false}
}