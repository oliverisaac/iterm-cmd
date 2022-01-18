package generate

import (
	"bytes"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/alessio/shellescape"
	"github.com/oliverisaac/go-levellogger"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"crypto/md5"
)

type CommandGenerator struct {
	logger    levellogger.Logger
	directory string
}

type FileCommandGenerator interface {
	CommandFile(exe string, args ...string) (filename string, err error)
}

// NewDefaultCommandGenerator creates a new command generator using the defaults
func NewDefaultCommandGenerator() FileCommandGenerator {
	return NewCommandGenerator(nil, "")
}

// NewCommandGenerator will create a comand generator pointing at the desired directory
// If the directory does not exist, it will be created on the first CommandFile call
// If the directory is an empty string, the default temp directory will be used
// If the logger is nil, logrus will be used
// If you don't want to specify either arg, just call NewDefaultCommandGenerator()
func NewCommandGenerator(logger levellogger.Logger, directory string) FileCommandGenerator {
	if directory == "" {
		directory = path.Join(os.TempDir(), "iT2cmd")
	}
	if logger == nil {
		logger = logrus.StandardLogger()
	}

	return &CommandGenerator{
		logger:    logger,
		directory: directory,
	}
}

// CommandFile generates a file in the tmp directory in the format of: /tmp/it2c/[a-f0-9]
// The filename is a hash of the file contents with as few characters as possible to be unique
func (cg *CommandGenerator) CommandFile(exe string, args ...string) (string, error) {
	log := cg.logger

	fullCommand := make([]string, 0, len(args)+1)
	fullCommand = append(fullCommand, exe)
	fullCommand = append(fullCommand, args...)

	err := os.MkdirAll(cg.directory, 0700)
	if err != nil {
		return "", errors.Wrap(err, "Ensuring directory exists")
	}

	quoted := shellescape.QuoteCommand(fullCommand)
	log.Tracef("Quoted command is: %s", quoted)

	commandHash := hashString(quoted)
	log.Tracef("Hash is: %s", commandHash)

	var filename string
	for i := 1; i <= len(commandHash); i++ {
		filename = path.Join(cg.directory, commandHash[:i])
		log.Tracef("Going to check file: %s", filename)

		exists, err := fileExists(filename)
		if err != nil {
			return "", errors.Wrap(err, "Checking if file exists")
		}

		// If it doesn't exist, then we will just write to that file
		if !exists {
			log.Debugf("Going to write command to new file: %s", filename)
			err = ioutil.WriteFile(filename, []byte(quoted), 0700)
			return filename, errors.Wrapf(err, "Writing to file %s", filename)
		}

		// If it does exist, then we need to see if the contents of this file match the hash we have
		fileHash, err := hashFile(filename)
		if fileHash == commandHash {
			log.Debugf("Command already written to existing file: %s", filename)
			return filename, nil
		}
	}

	return "", errors.New("Failed to find a unique filename for the given command")
}

// https://stackoverflow.com/a/12518877
func fileExists(filepath string) (bool, error) {
	if _, err := os.Stat(filepath); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, errors.Wrap(err, "Calling os.stat on file")
	}
}

func hashString(in string) (hashValue string) {
	r := bytes.NewReader([]byte(in))
	return hashReader(r)
}

func hashFile(filename string) (hashValue string, err error) {
	r, err := os.Open(filename)
	if err != nil {
		return "", errors.Wrap(err, "Opening file to try hashing it")
	}
	defer r.Close()
	return hashReader(r), nil
}

func hashReader(in io.Reader) (hashValue string) {
	// We are using md5 because it is fast and this doesn't need to be cryptographically secure
	hasher := md5.New()
	io.Copy(hasher, in)
	return hex.EncodeToString(hasher.Sum([]byte{}))
}
