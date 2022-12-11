package radiko

import (
    "encoding/xml"
)

type Programs struct {
    ID    string
    Extra map[string]string
}

func (p *Programs) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    p.Extra = make(map[string]string)
    for _, attr := range start.Attr {
        if attr.Name.Local == "id" {
            p.ID = attr.Value
        }
    }
    for {
        token, err := d.Token()
        if token == nil {
            break
        }
        if err != nil {
            return err
        }
        if t, ok := token.(xml.StartElement); ok {
            var data string
            if err := d.DecodeElement(&data, &t); err != nil {
                return err
            }
            p.Extra[t.Name.Local] = data
        }
    }
    return nil
}

func GetPrograms(area, station string) {
}