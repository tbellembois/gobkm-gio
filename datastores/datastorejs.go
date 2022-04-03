//go:build js && wasm

package datastores

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"github.com/tbellembois/gobkm-gio/globals"
)

type cookiedatastore struct {
	CServerURL      string
	CServerUsername string
	CHistorySize    string
}

func NewDatastore() *cookiedatastore {
	return &cookiedatastore{}
}

func (cd *cookiedatastore) InitDatastore() error {

	var (
		e error
		p Preferences
	)

	cd.CServerURL = "gobkmServerURL"
	cd.CServerUsername = "gobkmUsername"
	cd.CHistorySize = "gobkmHistorySize"

	if p, e = cd.LoadPreferences(); e != nil {
		return e
	}

	if p.ServerURL == "" {
		p.ServerURL = globals.DEFAULT_URL
		p.HistorySize = globals.DEFAULT_HISTORY

		if e := cd.SavePreferences(p); e != nil {
			return e
		}
	}

	return nil

}

func (cd *cookiedatastore) LoadPreferences() (Preferences, error) {

	var (
		e error
		p Preferences
		i int
	)

	c := js.Global().Get("document").Get("cookie").String()

	reurl := regexp.MustCompile(`serverURL=(\S+);{0,1}`)
	reusername := regexp.MustCompile(`username=(\S+);{0,1}`)
	rehistory := regexp.MustCompile(`historySize=(\d+);{0,1}`)

	if f := reurl.FindStringSubmatch(c); f != nil {
		p.ServerURL = strings.TrimSuffix(f[1], ";")
	}
	if f := reusername.FindStringSubmatch(c); f != nil {
		p.ServerUsername = strings.TrimSuffix(f[1], ";")
	}
	if f := rehistory.FindStringSubmatch(c); f != nil {
		h := strings.TrimSuffix(f[1], ";")
		if i, e = strconv.Atoi(h); e != nil {
			return Preferences{}, e
		}
		p.HistorySize = i
	}

	return p, nil

}

func (cd *cookiedatastore) SavePreferences(p Preferences) error {

	expires := time.Now().Local().AddDate(1, 0, 0).Format("Mon, 02-Jan-2006 15:04:05 MST")

	cookie := fmt.Sprintf("serverURL=%s; expires=%v; SameSite=None; Secure", p.ServerURL, expires)
	js.Global().Get("document").Set("cookie", cookie)
	cookie = fmt.Sprintf("username=%s; expires=%v; SameSite=None; Secure", p.ServerUsername, expires)
	js.Global().Get("document").Set("cookie", cookie)
	cookie = fmt.Sprintf("historySize=%s; expires=%v; SameSite=None; Secure", strconv.Itoa(p.HistorySize), expires)
	js.Global().Get("document").Set("cookie", cookie)

	return nil

}
