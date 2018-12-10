package main

import (
	"fmt"

	cliplugins "github.com/docker/cli/cli-plugins"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

func main() {
	plugin.Run(&plugin.Command{
		Command: cobra.Command{
			Use:   "helloworld",
			Short: "A basic Hello World plugin for tests",
		},
		RunPlugin: func(cmd *cobra.Command, dockerCli command.Cli, args []string) {
			fmt.Fprintln(dockerCli.Out(), "Hello World!")
		},
	}, cliplugins.Metadata{
		Version: "0.1.0",
		Vendor:  "Docker Inc.",
	})
}
