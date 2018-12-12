package main

import (
	"fmt"
	"math/rand"
	"time"
	"unicode"

	cliplugins "github.com/docker/cli/cli-plugins"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

func main() {
	plugin.Run(&plugin.Command{
		Command: cobra.Command{
			Use:   "kamoulox",
			Short: "Write down a useful kamoulox quote",
		},
		RunPlugin: func(cmd *cobra.Command, dockerCli command.Cli, args []string) {
			rand.Seed(time.Now().Unix())
			kamoulox := Kamoulox{}
			fmt.Fprintln(dockerCli.Out(), kamoulox.generate())
		},
	}, cliplugins.Metadata{
		Version: "0.1.0",
		Vendor:  "Docker Inc.",
	})
}

// Kamoulox generates down a famous kamoulox sentence
type Kamoulox struct{}

func (k *Kamoulox) generate() string {
	firstVerb := randomElementFrom(verbs)
	firstObject := randomElementFrom(objects)
	return fmt.Sprintf("%s %s et %s %s.", capitalizeFirst(firstVerb), firstObject, randomElementFromExcluding(verbs, firstVerb), randomElementFromExcluding(objects, firstObject))
}

func randomElementFrom(array []string) string {
	return array[rand.Intn(len(array))]
}

func capitalizeFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func randomElementFromExcluding(array []string, toExclude string) string {
	found := array[rand.Intn(len(array))]
	if found != toExclude {
		return found
	}
	return randomElementFromExcluding(array, toExclude)
}
