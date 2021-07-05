package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func openURL(url string) error {

	var cmd *exec.Cmd
	ros := runtime.GOOS
	switch ros {
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
