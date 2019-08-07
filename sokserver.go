package main

import (
	flag "github.com/spf13/pflag"
	"os"
	"path/filepath"

	log "github.com/franekjel/sokserver/logger"
	"github.com/franekjel/sokserver/server"
)

func main() {
	sokPath := flag.StringP("path", "p", "sok executable folder", "specify the SOK data path")
	debug := flag.BoolP("debug", "d", false, "enable debug output")
	exit := flag.BoolP("exit-on-error", "e", false, "exit on error")
	flag.Parse()

	log.InitLogger(*exit, *debug)
	log.Info("---Initializing SOK")

	if *sokPath == "sok executable folder" {
		ex, err := os.Executable()
		if err != nil {
			log.Fatal("Cannot get SOK localization, %s", err.Error())
		}
		*sokPath = filepath.Dir(ex)
	}

	log.Info("SOK folder %s", *sokPath)

	server.InitServer(*sokPath)
}
