package main

import (
    "fmt"
    "log"
    "time"
    "strconv"
    "github.com/mattn/go-runewidth"
    "github.com/rivo/tview"
    "github.com/gdamore/tcell/v2"
    "github.com/midry3125/cliRadio/mail"
    "github.com/midry3125/cliRadio/radiko"
)

type Widgets struct {
    App            *tview.Application
    StationList    *tview.List
    ProgramTable   *tview.List
    MailFrom       *tview.InputField
    MailTo         *tview.InputField
    MailMessage    *tview.TextArea
    CurrentPlayer   radiko.Stream
    Programs        *radiko.Programs
    nowPlaying       bool
    beforeidx        int
}

var (
    password mail.PasswordManager = mail.NewPasswordManager()
)

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
                case tcell.KeyCtrlT:
                    w.App.SetFocus(w.MailFrom)
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
        end := strconv.Itoa(p.End)
        w.ProgramTable.AddItem(p.Title, fmt.Sprintf("↑ : %s:%s～%s:%s", start[:2], start[2:], end[:2], end[2:]), 0, func(){})
    }
}

func (w *Widgets) UpdateNowProgram(id string, app *tview.Application) {
    index := -1
    now := time.Now()
    hours, _ := strconv.Atoi(now.Format("1504"))
    if w.Programs.Day != now.Day() {
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
    mail_from := tview.NewInputField().
                       SetLabel("From:")
    mail_to := tview.NewInputField().
                       SetLabel("To:")
    mail_from.SetDoneFunc(func(key tcell.Key) {
        if key == tcell.KeyEnter {
            app.SetFocus(mail_to)
        }
    })
    mail_title := tview.NewInputField().
                        SetLabel("Title:")
    mail_to.SetDoneFunc(func(key tcell.Key) {
        if key == tcell.KeyEnter {
            app.SetFocus(mail_title)
        }
    })
    mail_message := tview.NewTextArea()
    mail_message.SetBorder(true).
                 SetTitle("Message")
    mail_title.SetDoneFunc(func(key tcell.Key) {
        if key == tcell.KeyEnter {
            app.SetFocus(mail_message)
        }
    })
    data_flex := tview.NewFlex()
    data_flex.AddItem(stationList, 0, 1, true).
              AddItem(progTable, 0, 2, false)
    mail_addr_flex := tview.NewFlex()
    mail_addr_flex.AddItem(mail_from, 0, 1, false).
                   AddItem(mail_to, 0, 1, false)
    mail_message_flex := tview.NewFlex()
    mail_message_flex.AddItem(mail_message, 0, 1, false)
    mail_flex := tview.NewFlex().
                       SetDirection(tview.FlexRow)
    mail_flex.AddItem(mail_addr_flex, 0, 1, false).
              AddItem(mail_title, 0, 1, false).
              AddItem(mail_message_flex, 0, 3, false)
    mail_flex.SetBorder(true).
              SetTitle("Mail")
    layout := tview.NewFlex().
                    SetDirection(tview.FlexRow)
    layout.SetBorder(true).
           SetTitle(" Ctrl+S: Move to stations list  Ctrl+T: Move to mail form ")
    layout.AddItem(data_flex, 0, 1, true).
           AddItem(mail_flex, 0, 2, false)
    pages.AddPage("main", layout, true, true)
    app.SetRoot(pages, true)
    mail_message.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Key() == tcell.KeyCtrlW {
            modal := tview.NewModal().
                           AddButtons([]string{"Send", "Cancel"}).
                           SetText("Do you want to send this?")
            modal.SetDoneFunc(func(index int, label string) {
                if label == "Send" {
                    from := mail_from.GetText()
                    p, ok := password.Get(from)
                    if !ok {
                        modal.ClearButtons().
                              AddButtons([]string{"Back"}).
                              SetText("Error\nUnknown the password of this mail address").
                              SetDoneFunc(func(_ int, _ string) {
                            app.SetRoot(pages, true)
                            app.SetFocus(mail_message)
                        })
                  app.SetFocus(modal)
                  return
                    }
                    mailInfo := mail.Mail{from, mail_to.GetText(), p}
                    if err := mailInfo.Send(mail_title.GetText(), mail_message.GetText()); err != nil {
                        modal.ClearButtons().
                              AddButtons([]string{"Back"}).
                              SetText("Error\n"+err.Error()).
                              SetDoneFunc(func(_ int, _ string) {
                                  app.SetRoot(pages, true)
                                  app.SetFocus(mail_message)
                              })
                        app.SetFocus(modal)
                        return
                    } else {
                        mail_message.SetText("", true)
                        app.SetRoot(pages, true)
                        app.SetFocus(mail_message)
                    }
                }
            })
            app.SetRoot(modal, true)
            return nil
        } else {
            return event
        }
    })
    app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        switch event.Key() {
        case tcell.KeyCtrlT:
            app.SetFocus(mail_from)
        case tcell.KeyCtrlS:
            app.SetFocus(stationList)
        default:
            return event
        }
        return nil
    })
    return Widgets{
        App:            app,
        StationList:    stationList,
        ProgramTable:   progTable,
        MailFrom:       mail_from,
        MailTo:         mail_to,
        MailMessage:    mail_message,
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