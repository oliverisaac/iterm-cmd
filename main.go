package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/oliverisaac/iterm-cmd/generate"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var version string = "unset-version"
var commit string = "unset-commit"

func main() {
	workDir := path.Join(os.TempDir(), "it2cmd")
	if dir, isSet := os.LookupEnv("ITERM_CMD_DIR"); isSet {
		workDir = dir
	}

	if len(os.Args) <= 1 {
		log.Fatal("You must pass in at least one argument")
	}

	firstArg := os.Args[1]

	// Check if user is asking for help
	for _, h := range []string{"-h", "--help", "help", "--h"} {
		if strings.EqualFold(firstArg, h) {
			printHelp()
			return
		}
	}

	var err error
	if strings.EqualFold(firstArg, "handle") {
		err = handleClick(workDir)
	} else if strings.EqualFold(firstArg, "version") {
		fmt.Printf("iterm-cmd version: %s (%s)\n", version, commit)
		return
	} else {
		err = handleCreate(workDir)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func handleClick(workDir string) error {
	if len(os.Args) < 3 {
		return errors.New("You must pass a filename when calling with `handle`")
	}

	// "handle" is arg 1
	pathArg := os.Args[2]

	var lineNumber int
	var err error
	if len(os.Args) > 3 {
		lineNumberString := os.Args[3]
		if lineNumberString != "" {
			lineNumber, err = strconv.Atoi(os.Args[3])
			if err != nil {
				return errors.Wrap(err, "Second argument should be line number")
			}
		}
	}

	isDir, err := pathIsDirectory(pathArg)
	if err != nil {
		return errors.Wrap(err, "Checking if path is directory")
	}

	// Handle the scenario where we passed af file in the iterm_cmd dir
	if !isDir && strings.HasPrefix(pathArg, workDir) {
		logrus.Trace("fileArg starts with Workdir and is not a directory, going to treate as iterm-cmd")
		src, err := os.Open(pathArg)
		if err != nil {
			return errors.Wrap(err, "Failed to open file for reading")
		}
		defer src.Close()

		io.Copy(os.Stdout, src)
		return nil
	} else if isDir {
		logrus.Trace("fileArg is a directory, going to CD to it")
		cmd := []string{}

		if lookupBoolEnv("ITERM_CMD_PRINT_EASY_NAV", true) {
			easyNav := fmt.Sprintf(`echo "< $PWD ^ $( dirname %q )"`, pathArg)
			cmd = append(cmd, easyNav)
		}

		cmd = append(cmd, fmt.Sprintf("cd %q", pathArg))

		if lookupBoolEnv("ITERM_CMD_LS_AFTER_CD", true) {
			cmd = append(cmd, "ls")
		}

		fmt.Println(strings.Join(cmd, "; "))
		return nil
	} else {
		logrus.Trace("fileArg is a file, going to open it in EDITOR")
		editor, editorSet := os.LookupEnv("EDITOR")
		if !editorSet {
			logrus.Trace("EDITOR is not set, defaulting to vim")
			editor = "vim"
		}

		cmd := []string{editor}

		if lineNumber > 0 {
			if strings.HasSuffix(editor, "vim") {
				cmd = append(cmd, fmt.Sprintf("+%d", lineNumber))
			} else if strings.HasSuffix(editor, "code") {
				cmd = append(cmd, fmt.Sprintf("--goto=%s:%d", pathArg, lineNumber))
			}
		}

		cmd = append(cmd, pathArg)
		for _, a := range cmd {
			fmt.Printf("%q ", a)
		}
		fmt.Printf("\n")
	}

	return nil
}

func lookupBoolEnv(envVar string, defaultValue bool) bool {
	retVal := defaultValue
	if val, ok := os.LookupEnv(envVar); ok {
		var err error
		retVal, err = strconv.ParseBool(val)
		if err != nil {
			log.Error(errors.Wrapf(err, "Parsing boolean option %s", envVar))
		}
	}
	return retVal
}

// isDirectory determines if a file represented
// by `path` is a directory or not
// https://freshman.tech/snippets/go/check-if-file-is-dir/
func pathIsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, errors.Wrap(err, "Getting stat of possible directory")
	}

	return fileInfo.IsDir(), nil
}

func handleCreate(workDir string) error {
	exe := os.Args[1]
	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	file, err := generate.NewCommandGenerator(nil, workDir).CommandFile(exe, args...)
	if err != nil {
		return errors.Wrap(err, "Creating command file")
	}
	fmt.Printf("%s", file)
	return nil
}

func printHelp() {
	help := `
iterm-cmd takes an arbitrary command as an argument and then generates a file in ITERM_CMD_DIR named after the hash of that command. The hash is truncated to create unique filenames but also be as short as possible.

This is meant to be used in conjunction with the iTerm2 cmd-click functionality.

More details are here: github.com/oliverisaac/iterm-cmd/

If this command is called with 'handle' as its first argument, then it will do the handle portion of handling clicks.
`
	fmt.Println(help)
}
