package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"syscall/js"
)

// isValidUrl tests a string to determine if it is a well-structured url or not.
func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func openURL(url string) error {
	var cmd *exec.Cmd
	ros := runtime.GOOS
	switch ros {
	case "js":
		js.Global().Get("window").Call("open", url)
		return nil
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "android":
		cmd = exec.Command("/system/bin/am", "start", "-a", "android.intent.action.VIEW", "--user", "0", "-d", url)
	default:
		fmt.Printf("ros: %s.\n", ros)
		return nil
	}

	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

func basicAuth(username, password string) string {

	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))

}
