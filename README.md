# iTerm Command

This repo has a tool which does two things:

1. Handle clicks from iterm
2. Generate files which the iterm click-handler can understand

# Configuration

The only configuration option is `ITERM2_CMD_DIR` which sets where the files are created. This defaults to `${TMPDIR}/it2cmd` but can be set to any directory.

**I would suggest setting `ITERM2_CMD_DIR=/tmp/it2cmd`**

# Installation

### 1. Clone the Repo
Clone this repo into an appropriate directory:

```
mkdir -p "$GOPATH"/src/github.com/oliverisaac/iterm-cmd
git clone git@github.com:oliverisaac/iterm-cmd.git "$GOPATH"/src/github.com/oliverisaac/iterm-cmd
```

### 2. Build and install `iterm-cmd`

```
cd "$GOPATH"/src/github.com/oliverisaac/iterm-cmd
go install .
```

### 3. Edit your iTerm2 Settings to use the handle-click script:

1. Get the path to the `handle-click` script: 
```bash
echo "$GOPATH"/src/github.com/oliverisaac/iterm-cmd/handle-click
```

2. In iTerm2, go to the application preferences (`cmd` + `comma`)

3. Select Profiles -> Default -> Advanced

4. Under "Semantic History" select "Run Coprocess..."

5. In the text box that appears, paste in the path to the `handle-click` script and add: `'\1' '\2'` at the end. It should look something like: 
```
/Users/example/go/src/github.com/oliverisaac/iterm-cmd/handle-click '\1' '\2'
```

6. If you are using a custom `ITERM2_CMD_DIR` then you will want to prefix the command with that definition. You will end up with something like
```
ITERM2_CMD_DIR=/tmp/it2cmd $HOME/go/src/github.com/oliverisaac/iterm-cmd/handle-click '\1' '\2'
```

### 4. Use the script

You can now use `iterm-cmd` to generate a file which contains a command that iTerm2 can then execute. For example, run:

```
iterm-cmd echo "Hello world"
```

Then `cmd-click` on the output filename.
