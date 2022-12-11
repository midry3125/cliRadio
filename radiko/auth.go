package radiko

import (
    "log"
    "strconv"
    "encoding/base64"
    "net/http"
    //"os/exec"
)

const (
    StationIDUrl string = "https://radiko.jp/v3/station/region/full.xml"
    Auth1_Url    string = "https://radiko.jp/v2/api/auth1"
    Auth2_Url    string = "https://radiko.jp/v2/api/auth2"
    AuthKey      string = "bcd151073c03b352e1ef2fd66c32209da9ca0afa"
)

func Auth() string {
    req, _ := http.NewRequest("GET", Auth1_Url, nil)
    req.Header.Set("X-Radiko-App", "pc_html5")
    req.Header.Set("X-Radiko-App-Version", "0.0.1")
    req.Header.Set("X-Radiko-Device", "pc")
    req.Header.Set("X-Radiko-User", "dummy_user")
    client := new(http.Client)
    res, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()
    token := res.Header.Get("x-radiko-authtoken")
    offset, _ := strconv.Atoi(res.Header.Get("x-radiko-keyoffset"))
    length, _ := strconv.Atoi(res.Header.Get("x-radiko-keylength"))
    partialKey := base64.StdEncoding.EncodeToString([]byte(AuthKey[offset:offset+length]))
    req, _ = http.NewRequest("GET", Auth2_Url, nil)
    req.Header.Set("X-Radiko-Device", "pc")
    req.Header.Set("X-Radiko-User", "dummy_user")
    req.Header.Set("X-Radiko-AuthToken", token)
    req.Header.Set("X-Radiko-PartialKey", partialKey)
    res, err = client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()
    if res.Status != "200 OK" {
        log.Fatal("Failed to auth!")
    }
    return token
}

/*
func main() {
    id := flag.String("id", "", "Station ID")
    showIDs := flag.Bool("show-id", false, "Show all station id")
    flag.Parse()
    if *showIDs {
        fmt.Println("Name\tID")
        for _, s := range ParseStationIDPage().Stations {
            fmt.Println(s.Name+"\t"+s.ID)
        }
    } else {
        /*
        cmd := exec.Command("ffplay", "-nodisp", "-i", "C:\\Users\\minom\\Downloads\\LoveStruck.mp3")
        stdout, _ := cmd.StdoutPipe()
        stderr, _ := cmd.StderrPipe()
        buff_stdout := make([]byte, 1024)
        buff_stderr := make([]byte, 1024)
        var m int
        cmd.Start()
        for n, err := stdout.Read(buff_stdout); err != io.EOF; {
            if 0 < n {
                fmt.Print(string(buff_stdout[:n]))
            }
            if 0 < m {
                fmt.Fprint(os.Stderr, string(buff_stderr[:m]))
            }
            n, err = stdout.Read(buff_stdout)
            m, _ = stderr.Read(buff_stderr)
        }
        os.Exit(0)
        token := Auth()
        StartStream(*id, token)
    }
}
*/