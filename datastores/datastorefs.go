// +build !js

package datastores

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/mitchellh/go-homedir"
	"github.com/tbellembois/gobkm-gio/globals"
)

type fsdatastore struct {
	PreferenceFilePath string
}

func NewDatastore() *fsdatastore {
	return &fsdatastore{}
}

func (fsd *fsdatastore) InitDatastore() error {

	var (
		e error
		d string
	)

	ros := runtime.GOOS

	// Getting homedir.
	switch ros {
	case "android":
		d = path.Join("data", "data")
	default:
		if d, e = homedir.Dir(); e != nil {
			return e
		}
	}

	// Building per os config directory path.
	switch ros {
	case "linux", "darwin":
		d = path.Join(d, ".config", "gobkm-gio")
	case "windows":
		d = path.Join(d, "gobkm-gio")
	case "android":
		d = path.Join(d, "com.github.gobkm_gio")
	default:
		panic(fmt.Errorf("architecture not supported %s", ros))
	}

	// Creating app config directory if needed.
	if _, e = os.Stat(d); os.IsNotExist(e) {
		err := os.Mkdir(d, 0770)
		if err != nil {
			return err
		}
	} else if e != nil && !os.IsNotExist(e) {
		return e
	}

	// Setting full preferences file path.
	fsd.PreferenceFilePath = path.Join(d, "preferences.json")

	// Creating app config file if needed.
	if _, e = os.Stat(fsd.PreferenceFilePath); os.IsNotExist(e) {
		file, _ := json.MarshalIndent(Preferences{ServerURL: globals.DEFAULT_URL, HistorySize: globals.DEFAULT_HISTORY}, "", " ")
		if e = ioutil.WriteFile(fsd.PreferenceFilePath, file, 0644); e != nil {
			return e
		}
	} else if e != nil && !os.IsNotExist(e) {
		return e
	}

	return nil
}

func (fsd *fsdatastore) LoadPreferences() (Preferences, error) {

	var (
		e error
		f *os.File
		b []byte
		p Preferences
	)

	// Opening preferences file.
	if f, e = os.Open(fsd.PreferenceFilePath); e != nil {
		return Preferences{}, e
	}
	defer f.Close()

	// Loading into JSON.
	if b, e = ioutil.ReadAll(f); e != nil {
		return Preferences{}, e
	}
	if e = json.Unmarshal(b, &p); e != nil {
		return Preferences{}, e
	}

	return p, nil

}

func (fsd *fsdatastore) SavePreferences(p Preferences) error {

	var (
		e error
		b []byte
	)

	if b, e = json.MarshalIndent(p, "", " "); e != nil {
		return e
	}

	if e = ioutil.WriteFile(fsd.PreferenceFilePath, b, 0644); e != nil {
		return e
	}

	return nil

}
