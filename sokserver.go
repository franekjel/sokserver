package main

import (
	flag "github.com/spf13/pflag"
	"os"
	"path/filepath"

	. "github.com/franekjel/sokserver/logger"
	"github.com/franekjel/sokserver/server"
)

func main() {
	sokPath := flag.StringP("path", "p", "sok executable folder", "specify the SOK data path")
	debug := flag.BoolP("debug", "d", false, "enable debug output")
	exit := flag.BoolP("exit-on-error", "e", false, "exit on error")
	flag.Parse()

	InitLogger(*exit, *debug)
	Log(INFO, "---Initializing SOK")

	if *sokPath == "sok executable folder" {
		ex, err := os.Executable()
		if err != nil {
			Log(FATAL, "Cannot get SOK localization, %s", err.Error())
		}
		*sokPath = filepath.Dir(ex)
	}

	Log(INFO, "SOK folder %s", *sokPath)

	server.InitServer(*sokPath)
}
