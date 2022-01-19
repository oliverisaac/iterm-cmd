package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/oliverisaac/iterm-cmd/generate"
	"github.com/pkg/errors"
)

func main() {
	workDir := path.Join(os.TempDir(), "it2cmd")

	if dir, isSet := os.LookupEnv("ITERM2_CMD_DIR"); isSet {
		workDir = dir
	}

	if len(os.Args) <= 1 {
		log.Fatal("You must pass in at least one argument")
	}

	exe := os.Args[1]
	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	file, err := generate.NewCommandGenerator(nil, workDir).CommandFile(exe, args...)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Creating command file"))
	}

	fmt.Printf("%s", file)
}
