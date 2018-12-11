package cliplugins

import (
	"testing"

	"gotest.tools/icmd"
)

func WithPluginDir(d string) icmd.CmdOp {
	return icmd.WithEnv("DOCKER_CLI_PLUGIN_EXTRA_DIRS=" + d)
}

func TestRunNonexisting(t *testing.T) {
	res := icmd.RunCmd(icmd.Command("docker", "nonexistent"))
	res.Assert(t, icmd.Expected{
		ExitCode: 1,
		Err:      "docker: 'nonexistent' is not a docker command.\nSee 'docker --help'",
	})
}

func TestRunBad(t *testing.T) {
	res := icmd.RunCmd(icmd.Command("docker", "badmeta"))
	res.Assert(t, icmd.Expected{
		ExitCode: 1,
		Err:      "docker: 'badmeta' is not a docker command.\nSee 'docker --help'",
	})
}

func TestRunGood(t *testing.T) {
	res := icmd.RunCmd(icmd.Command("docker", "helloworld"))
	res.Assert(t, icmd.Expected{
		ExitCode: 0,
		Out:      "Hello World!",
	})
}
