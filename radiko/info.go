package radiko

import (
    "fmt"
    "io"
    "log"
    "regexp"
    "encoding/xml"
    "net/http"
)

const (
    areaUrl    string = "https://radiko.jp/area"
    stationUrl string = "https://radiko.jp/v3/station/list/%s.xml"
)

var (
    re = regexp.MustCompile("class=\"(.+)\"")
)

type Stations struct {
    XMLName  xml.Name      `xml:"stations"`
    Stations []StationInfo `xml:"station"`
}

type StationInfo struct {
    XMLName xml.Name `xml:"station"`
    Name    string   `xml:"name"`
    ID      string   `xml:"id"`
}

func GetResponse(url string) (response []byte, err error) {
    res, err := http.Get(url)
    if err != nil {
        return response, err
    }
    defer res.Body.Close()
    response, _ = io.ReadAll(res.Body)
    return response, nil
}

func ParseStationIDPage() ([]StationInfo, string) {
    res, err := GetResponse(areaUrl)
    if err != nil {
        log.Fatal(err)
    }
    areaCode := string(re.FindSubmatch(res)[1])
    res, err = GetResponse(fmt.Sprintf(stationUrl, areaCode))
    if err != nil {
        log.Fatal(err)
    }
    var stations Stations
    err = xml.Unmarshal(res, &stations)
    if err != nil {
        log.Fatal(err)
    }
    return stations.Stations, areaCode
}