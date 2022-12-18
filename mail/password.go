package mail

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

var (
	directory  string = filepath.Join(os.Getenv("APPDATA"), "cliRadio")
	ConfigFile string = filepath.Join(directory, "mail.json")
)

type PasswordManager struct {
	Data map[string]string
}

func NewPasswordManager() PasswordManager {
	var (
		Default PasswordManager
		data    map[string]string
	)
	f, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return Default
	}
	err = json.Unmarshal(f, &data)
    return PasswordManager{data}
}

func (p PasswordManager) Get(addr string) (string, bool) {
	v, ok := p.Data[addr]
	if ok {
	    return v, true
	} else {
		return "", false
	}
}

func init() {
	if _, err := os.Stat(directory); err != nil {
		os.Mkdir(directory, 0777)
	}
}