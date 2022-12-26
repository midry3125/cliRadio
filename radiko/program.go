package radiko

import (
    "fmt"
    "log"
    "time"
    "encoding/xml"
    "github.com/rivo/tview"
)

const (
    programUrl string = "https://radiko.jp/v3/program/station/date/%s/%s.xml"
)

type Programs struct {
    XMLName  xml.Name      `xml:"radiko"`
    Programs []ProgramInfo `xml:"stations>station>progs>prog"`
    Day     int
}

type ProgramInfo struct {
    XMLName  xml.Name `xml:"prog"`
    Start    int      `xml:"ftl,attr"`
    End      int      `xml:"tol,attr"`
    Title    string   `xml:"title"`
}

func GetPrograms(station string, app *tview.Application) *Programs {
    res, err := GetResponse(fmt.Sprintf(programUrl, time.Now().Format("20060102"), station))
    if err != nil {
        app.Stop()
        log.Fatal(err)
    }
    p := Programs{Day: time.Now().Day()}
    err = xml.Unmarshal(res, &p)
    if err != nil {
        app.Stop()
        log.Fatal(err)
    }
    return &p
}