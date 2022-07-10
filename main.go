package main

import (
	"flag"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/croyleje/gorm/cmd"
	"github.com/croyleje/gorm/ui"
)

func initState(debug bool) (ui.Model, *os.File) {
	var loggerFile *os.File
	var err error

	if debug {
		loggerFile, err = tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatalf("error initializing logger: %v", err)
		}
	}

	return ui.InitialModel(), loggerFile
}

func main() {
	debug := flag.Bool("debug", false, "enables error log ./debug.log")
	list := flag.Bool("list", false, "cli list of item")
	put := flag.Bool("put", false, "move file to trash")

	flag.Parse()

	model, logger := initState(*debug)

	switch {
	case *list:
		cmd.CliTrashList()

	case *put:
		cmd.TrashPut(flag.Args())

	default:
		if logger != nil {
			defer logger.Close()
		}

		p := tea.NewProgram(
			model,
			tea.WithAltScreen(),
		)
		if err := p.Start(); err != nil {
			log.Fatal(err)
		}
	}

}
