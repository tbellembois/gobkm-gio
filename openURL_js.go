package main

import (
	"syscall/js"
)

func openURL(url string) error {

	js.Global().Get("window").Call("open", url)
	return nil

}
