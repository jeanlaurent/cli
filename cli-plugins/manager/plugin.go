package manager

import (
	"encoding/json"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	cliplugins "github.com/docker/cli/cli-plugins"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	pluginNameRe = regexp.MustCompile("^[a-z][a-z0-9]*$")
)

type Plugin struct {
	cliplugins.Metadata

	Name string
	Path string

	// Err is non-nil if the plugin failed one of the candidate tests.
	Err error `json:",omitempty"`

	ShadowedPaths []string `json:",omitempty"`
}

// NewPlugin determines if the given candidate is valid and returns a
// Plugin.  If the candidate fails one of the tests then `Plugin.Err`
// is set, but the `Plugin` is still returned with no error. An error
// is only returned due to a non-recoverable error.
func NewPlugin(c Candidate, rootcmd *cobra.Command) (Plugin, error) {
	path := c.Path()
	if path == "" {
		return Plugin{}, errors.New("plugin candidate path cannot be empty")
	}

	// The candidate listing process should have skipped anything
	// which would fail here, so there are all real errors.
	fullname := filepath.Base(path)
	if fullname == "." {
		return Plugin{}, errors.Errorf("unable to determine basename of plugin candidate %q", path)
	}
	if runtime.GOOS == "windows" {
		exe := ".exe"
		if !strings.HasSuffix(fullname, exe) {
			return Plugin{}, errors.Errorf("plugin candidate %q lacks required %q suffix", path, exe)
		}
		fullname = strings.TrimSuffix(fullname, exe)
	}
	if !strings.HasPrefix(fullname, cliplugins.NamePrefix) {
		return Plugin{}, errors.Errorf("plugin candidate %q does not have %q prefix", path, cliplugins.NamePrefix)
	}

	p := Plugin{
		Name: strings.TrimPrefix(fullname, cliplugins.NamePrefix),
		Path: path,
	}

	// Now apply the candidate tests, so these update p.Err.
	if !pluginNameRe.MatchString(p.Name) {
		p.Err = errors.Errorf("plugin candidate %q did not match %q", p.Name, pluginNameRe.String())
		return p, nil
	}

	if rootcmd != nil {
		for _, cmd := range rootcmd.Commands() {
			if cmd.Name() == p.Name || cmd.HasAlias(p.Name) {
				p.Err = errors.New("plugin duplicates builtin command")
				return p, nil
			}
		}
	}

	// We are supposed to check for relevant execute permissions here. Instead we rely on an attempt to execute.
	meta, err := c.Metadata()
	if err != nil {
		p.Err = errors.Wrap(err, "failed to fetch metadata")
		return p, nil
	}

	if err := json.Unmarshal(meta, &p.Metadata); err != nil {
		p.Err = errors.Wrap(err, "invalid metadata")
		return p, nil
	}

	return p, nil
}

func (p *Plugin) IsValid() error {
	return p.Err
}
