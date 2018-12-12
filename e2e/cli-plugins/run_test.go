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

func TestHelpNonexisting(t *testing.T) {
	res := icmd.RunCmd(icmd.Command("docker", "help", "nonexistent"))
	res.Assert(t, icmd.Expected{
		ExitCode: 1,
		Err:      "unknown help topic: nonexistent",
	})

	res = icmd.RunCmd(icmd.Command("docker", "nonexistent", "--help"))
	res.Assert(t, icmd.Expected{
		ExitCode: 0,
		Out: "Usage:	docker [OPTIONS] COMMAND",
	})
}

func TestRunBad(t *testing.T) {
	res := icmd.RunCmd(icmd.Command("docker", "badmeta"))
	res.Assert(t, icmd.Expected{
		ExitCode: 1,
		Err:      "docker: 'badmeta' is not a docker command.\nSee 'docker --help'",
	})
}

func TestHelpBad(t *testing.T) {
	res := icmd.RunCmd(icmd.Command("docker", "help", "badmeta"))
	res.Assert(t, icmd.Expected{
		ExitCode: 1,
		Err:      "unknown help topic: badmeta",
	})

	res = icmd.RunCmd(icmd.Command("docker", "badmeta", "--help"))
	res.Assert(t, icmd.Expected{
		ExitCode: 0,
		Out: "Usage:	docker [OPTIONS] COMMAND",
	})
}

func TestRunGood(t *testing.T) {
	res := icmd.RunCmd(icmd.Command("docker", "helloworld"))
	res.Assert(t, icmd.Expected{
		ExitCode: 0,
		Out:      "Hello World!",
	})
}

func TestHelpGood(t *testing.T) {
	res := icmd.RunCmd(icmd.Command("docker", "help", "helloworld"))
	res.Assert(t, icmd.Expected{
		ExitCode: 0,
		Out: "Usage:	docker helloworld",
	})
	res = icmd.RunCmd(icmd.Command("docker", "helloworld", "--help"))
	res.Assert(t, icmd.Expected{
		ExitCode: 0,
		Out: "Usage:	docker helloworld",
	})
}
