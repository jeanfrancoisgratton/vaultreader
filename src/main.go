package main

import (
	"fmt"
	"os"
	"path/filepath"
	"vaultreader/cmd"
)

func main() {
	// Create the vaultreader logdir
	base := os.Getenv("XDG_STATE_HOME")
	if base == "" {
		base = filepath.Join(os.Getenv("HOME"), ".local", "state")
	}
	if err := os.MkdirAll(filepath.Join(base, "vaultreader"), 0755); err != nil {
		fmt.Println(err)
		os.Exit(11)
	}

	// "Touch" the logfile so we avoid an error at the very first use of the tool
	if f, e := os.OpenFile(filepath.Join(base, "vaultreader", "vaultreader.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); e != nil {
		fmt.Println(e)
		os.Exit(12)
	} else {
		f.Close()
	}
	cmd.Execute()
}
