package radiko

import (
    "fmt"
    "os"
    "time"
    "os/exec"
    "github.com/rivo/tview"
)

type Stream struct {
    ID        string
    Token     string
    Area      string
    Player    *exec.Cmd
    Cancel    bool
}

func (s *Stream) StartStream(app *tview.Application, updateProgs func()) {
    url := fmt.Sprintf("https://si-f-radiko.smartstream.ne.jp/so/playlist.m3u8?station_id=%s&l=15&lsid=27bc2ff61afcd9d30da1be25a4731e14&type=b", s.ID)
    _, err := exec.LookPath("ffplay")
    if err != nil {
        app.Stop()
        fmt.Println("[\033[31mError\033[0m] ffplay is not found")
        os.Exit(1)
    }
    s.Player = exec.Command(
        "ffplay",
        "-nodisp",
        "-vn",
        "-headers",
        "X-Radiko-AuthToken: "+s.Token,
        "-i", url,
    )
    s.Player.Start()
    defer s.Player.Wait()
    defer s.Stop()
    for {
        if s.Cancel {
            return
        }
        app.QueueUpdateDraw(updateProgs)
        time.Sleep(time.Second * 10)
    }
}

func (s *Stream) Stop() {
    if s.Player != nil {
        s.Player.Process.Kill()
        s.Cancel = true
    }
}