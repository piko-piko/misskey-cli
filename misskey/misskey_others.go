//go:build !windows
// +build !windows

package misskey

import (
	"fmt"
	"os"
	"golang.org/x/crypto/ssh/terminal"
)

func terminalWidth() int {
	width, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %+v", err)
		os.Exit(1)
	}
	return width
}
