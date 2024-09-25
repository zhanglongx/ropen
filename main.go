package main

import (
	"flag"

	"github.com/zhanglongx/ropen/pkg"
)

func main() {
	cfg := flag.String("cfg", "", "config file path")
	debug := flag.Bool("debug", false, "config file path")
	port := flag.Int("port", -1, "port number")
	version := flag.Bool("version", false, "display version")

	flag.Parse()

	if *version {
		println(pkg.APP_NAME, pkg.APP_VERSION)
		return
	}

	if *debug {
		pkg.SetLevel(pkg.LevelDebug)
	}

	args := flag.Args()
	if len(args) < 1 {
		panic("missing path parameters")
	}

	path := args[0]

	app, err := pkg.NewApp(*cfg, *port)
	if err != nil {
		panic(err)
	}

	if err := app.Run(path); err != nil {
		panic(err)
	}
}
