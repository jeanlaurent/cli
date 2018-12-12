package manager

import (
	"fmt"
	"strings"
	"testing"

	cliplugins "github.com/docker/cli/cli-plugins"
	"github.com/spf13/cobra"
	"gotest.tools/assert"
)

var execError = fmt.Errorf("exec error")

type mockCandidate struct {
	path string
	exec bool
	meta string
}

func (c *mockCandidate) Path() string {
	return c.path
}

func (c *mockCandidate) Metadata() ([]byte, error) {
	if !c.exec {
		return nil, execError
	}
	return []byte(c.meta), nil
}

func TestValidateCandidate(t *testing.T) {
	var (
		goodPluginName = cliplugins.NamePrefix + "goodplugin"
		goodVersion    = "0.1.0"

		builtinName  = cliplugins.NamePrefix + "builtin"
		builtinAlias = cliplugins.NamePrefix + "alias"

		badPrefixPath  = "/usr/local/libexec/cli-plugins/wobble"
		badNamePath    = "/usr/local/libexec/cli-plugins/docker-123456"
		goodPluginPath = "/usr/local/libexec/cli-plugins/" + goodPluginName
	)

	fakeroot := &cobra.Command{Use: "docker"}
	fakeroot.AddCommand(&cobra.Command{
		Use: strings.TrimPrefix(builtinName, cliplugins.NamePrefix),
		Aliases: []string{
			strings.TrimPrefix(builtinAlias, cliplugins.NamePrefix),
		},
	})

	for _, tc := range []struct {
		c    *mockCandidate
		meta string

		// Either err or invalid may be non-empty, but not both (both can be empty for a good plugin).
		err     string
		invalid string
	}{
		/* Each failing one of the tests */
		{c: &mockCandidate{path: ""}, err: "plugin candidate path cannot be empty"},
		{c: &mockCandidate{path: badPrefixPath}, err: fmt.Sprintf("does not have %q prefix", cliplugins.NamePrefix)},
		{c: &mockCandidate{path: badNamePath}, invalid: "did not match"},
		{c: &mockCandidate{path: builtinName}, invalid: "plugin duplicates builtin command"},
		{c: &mockCandidate{path: builtinAlias}, invalid: "plugin duplicates builtin command"},
		{c: &mockCandidate{path: goodPluginPath, exec: false}, invalid: "exec error"},
		{c: &mockCandidate{path: goodPluginPath, exec: true, meta: `xyzzy`}, invalid: "invalid character"},
		// This one should work
		{c: &mockCandidate{path: goodPluginPath, exec: true, meta: fmt.Sprintf(`{"Version": %q}`, goodVersion)}},
	} {
		p, err := NewPlugin(tc.c, fakeroot)
		if tc.err != "" {
			assert.ErrorContains(t, err, tc.err)
		} else if tc.invalid != "" {
			assert.NilError(t, err)
			assert.ErrorContains(t, p.IsValid(), tc.invalid)
		} else {
			assert.NilError(t, err)
			assert.Equal(t, cliplugins.NamePrefix+p.Name, goodPluginName)
			assert.Equal(t, p.Version, goodVersion)
		}
	}
}
