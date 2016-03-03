package links

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/mitchellh/cli"
)

// Command is a Command implementation that runs a Consul agent.
// The command will not end unless a shutdown message is sent on the
// ShutdownCh. If two messages are sent on the ShutdownCh it will forcibly
// exit.
type Command struct {
	Revision          string
	Version           string
	VersionPrerelease string
	Ui                cli.Ui
}

func (c *Command) Run(args []string) int {
	cmd := exec.Command("./agent", "arguments")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	return 0
}

func (c *Command) Synopsis() string {
	return "Run a Open-Falcon agent"
}

func (c *Command) Help() string {
	helpText := `
Usage: Open-Falcon agent [options]

  Starts the Consul agent and runs until an interrupt is received. The
  agent represents a single node in a cluster.

Options:

  -bind=0.0.0.0            Sets the bind address for cluster communication
  -http-port=8500          Sets the HTTP API port to listen on
  -bootstrap-expect=0      Sets server to expect bootstrap mode.
  -client=127.0.0.1        Sets the address to bind for client access.
                           This includes RPC, DNS, HTTP and HTTPS (if configured)
  -config-file=foo         Path to a JSON file to read configuration from.
                           This can be specified multiple times.
  -config-dir=foo          Path to a directory to read configuration files
                           from. This will read every file ending in ".json"
                           as configuration in this directory in alphabetical
                           order. This can be specified multiple times.
 `
	return strings.TrimSpace(helpText)
}
