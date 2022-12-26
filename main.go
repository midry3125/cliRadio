package main

import (
    "fmt"
    "log"
    "time"
    "strconv"
    "github.com/mattn/go-runewidth"
    "github.com/rivo/tview"
    "github.com/gdamore/tcell/v2"
    "github.com/midry3125/cliRadio/radiko"
)

type Widgets struct {
    App            *tview.Application
    StationList    *tview.List
    ProgramTable   *tview.List
    CurrentPlayer   radiko.Stream
    Programs        *radiko.Programs
    nowPlaying       bool
    beforeidx        int
}

func (w *Widgets) SetData(stations []radiko.StationInfo) {
    for _, s := range stations {
        w.StationList.AddItem(s.Name, "", 0, func() {
            index := w.StationList.GetCurrentItem()
            station := stations[index]
            w.StationList.SetItemText(w.beforeidx, stations[w.beforeidx].Name, "")
            w.StationList.SetItemText(index, station.Name, "↑ Now playing...")
            if w.nowPlaying {
                w.CurrentPlayer.Stop()
            }
            w.UpdatePrograms(station.ID, w.App)
            w.CurrentPlayer = radiko.Stream{
                ID:        station.ID,
                Token:     radiko.Auth(),
            }
            go w.CurrentPlayer.StartStream(w.App, func() {
                w.UpdateNowProgram(station.ID, w.App)
            })
            w.nowPlaying = true
            w.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
                switch event.Key() {
                case tcell.KeyCtrlS:
                    w.App.SetFocus(w.StationList)
                case tcell.KeyCtrlC:
                    w.CurrentPlayer.Stop()
                    return event
                default:
                    return event
                }
                return nil
            })
            w.beforeidx = index
        })
    }
}

func (w *Widgets) UpdatePrograms(id string, app *tview.Application) {
    w.Programs = radiko.GetPrograms(id, app)
    w.ProgramTable.Clear()
    for _, p := range w.Programs.Programs {
        start := strconv.Itoa(p.Start)
        switch len(start) {
        case 3:
            start = start[:1]+":"+start[1:]
        case 4:
            start = start[:2]+":"+start[2:]
        }
        end := strconv.Itoa(p.End)
        switch len(end) {
        case 3:
            end = end[:1]+":"+end[1:]
        case 4:
            end = end[:2]+":"+end[2:]
        }
        w.ProgramTable.AddItem(p.Title+fmt.Sprintf("  %s～%s", start, end), "", 0, func(){})
    }
}

func (w *Widgets) UpdateNowProgram(id string, app *tview.Application) {
    index := -1
    now := time.Now()
    now_h := now.Hour()
    hours, _ := strconv.Atoi(now.Format("1504"))
    if (w.Programs.Day != now.Day()) || (0 <= now_h && now_h <= 6) {
        hours += 2400
    }
    for n, p := range w.Programs.Programs {
        if hours < p.End {
            index = n
            break
        }
    }
    if index == -1 {
        w.UpdatePrograms(id, app)
        w.ProgramTable.SetCurrentItem(0)
    } else {
        w.ProgramTable.SetCurrentItem(index)
    }
}

func createApp() Widgets {
    app := tview.NewApplication()
    pages := tview.NewPages()
    stationList := tview.NewList()
    stationList.SetMainTextColor(tcell.GetColor("#00FFFF")).
                SetSelectedBackgroundColor(tcell.GetColor("#32cd32"))
    progTable := tview.NewList().
                       SetSelectedBackgroundColor(tcell.GetColor("#00FFFF"))
    progTable.SetBorder(true).
              SetTitle("Programs")
    layout := tview.NewFlex().
                    SetDirection(tview.FlexRow)
    layout.SetBorder(true).
           SetTitle(" Ctrl+S: Move to stations list  Ctrl+T: Move to mail form ")
    layout.AddItem(stationList, 0, 1, true).
           AddItem(progTable, 0, 2, false)
    pages.AddPage("main", layout, true, true)
    app.SetRoot(pages, true)
    app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Key() == tcell.KeyCtrlS {
            app.SetFocus(stationList)
        }
        return event
    })
    return Widgets{
        App:            app,
        StationList:    stationList,
        ProgramTable:   progTable,
    }
}

func main() {
    widget := createApp()
    stations := radiko.ParseStationIDPage()
    widget.SetData(stations)
    err := widget.App.Run()
    if err != nil {
        log.Fatal(err)
    }
}

func init() {
    runewidth.DefaultCondition = &runewidth.Condition{EastAsianWidth: false}
}