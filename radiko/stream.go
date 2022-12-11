package radiko

import (
    "fmt"
    "os"
    "time"
    "syscall"
    "os/exec"
    "os/signal"
    "github.com/rivo/tview"
)

type Stream struct {
    ID        string
    Token     string
    Player    *exec.Cmd
    Canceller chan struct{}
}

type output struct {
    View *tview.TextView
    App  *tview.Application
}

var newline []byte = []byte("\n\r")

func (o *output) Write(b []byte) (int, error) {
    o.App.QueueUpdateDraw(func() {
        o.View.Write(b)
        o.View.Write(newline)
    })
    return len(b), nil
}

func (s *Stream) StartStream(app *tview.Application, view *tview.TextView) {
    url := fmt.Sprintf("https://f-radiko.smartstream.ne.jp/%s/_definst_/simul-stream.stream/playlist.m3u8", s.ID)
    _, err := exec.LookPath("ffplay")
    if err != nil {
        app.Stop()
        fmt.Println("[\033[31mError\033[0m] ffplay is not found")
        os.Exit(1)
    }
    s.Player = exec.Command("ffplay", "-nodisp", "-headers", fmt.Sprintf("X-Radiko-AuthToken: %s", s.Token), "-i", url)
    s.Player.Stdout = &output{View: view, App: app}
    s.Player.Stderr = &output{View: view, App: app}
    s.Player.Start()
    defer s.Player.Wait()
    defer s.Stop()
    s.Canceller = make(chan struct{})
    go func(s Stream) {
        trap := make(chan os.Signal, 1)
        signal.Notify(trap, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
        for {
            select {
            case <-trap:
                s.Stop()
                os.Exit(0)
            case <-s.Canceller:
                return
            }
        }
    }(*s)
    for {
        time.Sleep(time.Minute * 60)
    }
}

func (s Stream) Stop() {
    if s.Player != nil {
        s.Player.Process.Kill()
        close(s.Canceller)
    }
}