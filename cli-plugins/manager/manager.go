package manager

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"

	cliplugins "github.com/docker/cli/cli-plugins"
	"github.com/spf13/cobra"
)

type ErrPluginNotFound string

func (e ErrPluginNotFound) Error() string {
	return "Error: No such CLI plugin: " + string(e)
}

func (_ ErrPluginNotFound) NotFound() bool {
	return true
}

var (
	pluginDirs     []string
	pluginDirsOnce sync.Once
)

func PluginDirs() []string {
	pluginDirsOnce.Do(func() {
		// Mostly for test.
		if ds := os.Getenv("DOCKER_CLI_PLUGIN_EXTRA_DIRS"); ds != "" {
			pluginDirs = append(pluginDirs, strings.Split(ds, ":")...)
		}

		pluginDirs = append(pluginDirs,
			filepath.Join(os.Getenv("HOME"), ".docker/cli-plugins"),
			"/usr/local/lib/docker/cli-plugins", "/usr/local/libexec/docker/cli-plugins",
			"/usr/lib/docker/cli-plugins", "/usr/libexec/docker/cli-plugins",
		)
	})
	return pluginDirs
}

func addPluginCandidatesFromDir(res map[string][]string, d string) error {
	dentries, err := ioutil.ReadDir(d)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		// Portable? This seems to be as good as it gets :-/
		if serr, ok := err.(*os.SyscallError); ok && serr.Err == syscall.ENOTDIR {
			return nil
		}
		return err
	}
	for _, dentry := range dentries {
		switch dentry.Mode() & os.ModeType {
		case 0, os.ModeSymlink:
			// Regular file or symlink, keep going
		default:
			// Something else, ignore.
			continue
		}
		name := dentry.Name()
		if !strings.HasPrefix(name, cliplugins.NamePrefix) {
			continue
		}
		name = strings.TrimPrefix(name, cliplugins.NamePrefix)
		if runtime.GOOS == "windows" {
			exe := ".exe"
			if !strings.HasSuffix(name, exe) {
				continue
			}
			name = strings.TrimSuffix(name, exe)
		}
		res[name] = append(res[name], filepath.Join(d, dentry.Name()))
	}
	return nil
}

// listPluginCandidates allows the dirs to be specified for testing purposes.
func listPluginCandidates(dirs []string) (map[string][]string, error) {
	result := make(map[string][]string)
	for _, d := range dirs {
		if err := addPluginCandidatesFromDir(result, d); err != nil {
			return nil, err // Or return partial result?
		}
	}
	return result, nil
}

// ListPluginCandidates returns a map from plugin name to the list of (unvalidated) Candidates. The list is in descending order of priority.
func ListPluginCandidates() (map[string][]string, error) {
	return listPluginCandidates(PluginDirs())
}

func ListPlugins(rootcmd *cobra.Command) ([]Plugin, error) {
	candidates, err := ListPluginCandidates()
	if err != nil {
		return nil, err
	}

	var plugins []Plugin
	for _, paths := range candidates {
		if len(paths) == 0 {
			continue
		}
		c := &candidate{paths[0]}
		p, err := NewPlugin(c, rootcmd)
		if err != nil {
			return nil, err
		}
		p.ShadowedPaths = paths[1:]
		plugins = append(plugins, p)
	}

	return plugins, nil
}

// FindPlugin finds a valid plugin, if the first candidate is invalid then returns an error
func FindPlugin(name string, rootcmd *cobra.Command, includeShadowed bool) (Plugin, error) {
	if !pluginNameRe.MatchString(name) {
		// We treat this as "not found" so that callers will
		// fallback to their "invalid" command path.
		return Plugin{}, ErrPluginNotFound(name)
	}
	exename := cliplugins.NamePrefix + name
	if runtime.GOOS == "windows" {
		exename = exename + ".exe"
	}
	var plugin Plugin
	for _, d := range PluginDirs() {
		path := filepath.Join(d, exename)

		// We stat here rather than letting the exec tell us
		// ENOENT because the latter does not distinguish a
		// file not existing from its dynamic loader or one of
		// its libraries not existing.
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		if plugin.Path == "" {
			c := &candidate{path: path}
			var err error
			if plugin, err = NewPlugin(c, rootcmd); err != nil {
				return Plugin{}, err
			}
			if !includeShadowed {
				return plugin, nil
			}
		} else {
			plugin.ShadowedPaths = append(plugin.ShadowedPaths, path)
		}
	}
	if plugin.Path == "" {
		return Plugin{}, ErrPluginNotFound(name)
	}
	return plugin, nil
}

// runPluginCommand returns a Cmd which will run the named plugin.
// rootcmd is used to detect conficts with builtin commands.
// The error returned is an ErrPluginNotFound if no plugin was found or if the first candidate plugin was invalid somehow.
func runPluginCommand(name string, rootcmd *cobra.Command, args []string) (*exec.Cmd, error) {
	plugin, err := FindPlugin(name, rootcmd, false)
	if err != nil {
		return nil, err
	}
	if err := plugin.IsValid(); err != nil {
		return nil, ErrPluginNotFound(name)
	}
	return exec.Command(plugin.Path, args...), nil
}

func PluginRunCommand(name string, rootcmd *cobra.Command) (*exec.Cmd, error) {
	// This uses the full original args, not the args which may
	// have been provided by cobra to our caller. This is because
	// they lack e.g. global options which we must propagate here.
	return runPluginCommand(name, rootcmd, os.Args[1:])
}

func PluginHelpCommand(name string, rootcmd *cobra.Command) (*exec.Cmd, error) {
	return runPluginCommand(name, rootcmd, []string{"help", name})
}
