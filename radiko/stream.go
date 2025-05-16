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
    url := fmt.Sprintf("https://f-radiko.smartstream.ne.jp/%s/_definst_/simul-stream.stream/playlist.m3u8", s.ID)
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
        fmt.Sprintf("X-Radiko-AuthToken: %s", s.Token),
        url,
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