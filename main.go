package main

import (
	"flag"
	"fmt"
)

var (
	helpCommand bool

	serveCommand string

	explainCommand string
	analyseCommand string
	extractCommand string

	explainTarget string
)

func main() {
	flag.Parse()

	if helpCommand {
		flag.Usage()
	} else if analyseCommand != "" {
		analyseCliHandler(analyseCommand)
	} else if extractCommand != "" {
		extractCliHandler(extractCommand)
	} else if explainCommand != "" && explainTarget != "" {
		explainCliHandler()
	} else if serveCommand != "" {
		serveMain(serveCommand)
	} else {
		flag.Usage()
	}
}

func init() {
	flag.BoolVar(&helpCommand, "h", false, "Print help text")
	flag.StringVar(&explainTarget, "t", "", "Set the `target`, used with -e")
	flag.StringVar(&serveCommand, "s", "", "Serve on given `address`")
	flag.StringVar(&analyseCommand, "a", "", "Analyse given `file`, and save with .ana suffix")
	flag.StringVar(&extractCommand, "ex", "", "Extract given `file`, and save with .ext suffix")
	flag.StringVar(&explainCommand, "e", "", "Explain given `file`, and save with .exp suffix, used with -t")

	flag.Usage = usage
}

func usage() {
	_, _ = fmt.Println(`EDL Explainer Usage:`)

	flag.PrintDefaults()
}
