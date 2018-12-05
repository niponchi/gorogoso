package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/zapkub/gorogoso/runner"
)

func main() {
	watchPattern := flag.String("watch", "*.go", "Glob pattern you want to watch")
	entryfile := flag.String("entry", "main.go", "entryfile path")
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	watchPathPattern := path.Join(dir, *watchPattern)
	entryfilePath := path.Join(dir, *entryfile)

	pid := runner.Monit(watchPathPattern, entryfilePath)
	for {
		<-pid
	}

}
