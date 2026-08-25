package cli

import (
	"fmt"
	"io"
)

const (
	ExitOK    = 0
	ExitError = 1
	ExitUsage = 2
)

type Env struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func Run(args []string, env Env) int {
	if env.Stdout == nil {
		env.Stdout = io.Discard
	}
	if env.Stderr == nil {
		env.Stderr = io.Discard
	}
	if env.Stdin == nil {
		env.Stdin = nilReader{}
	}
	if len(args) == 0 {
		writeUsage(env.Stderr)
		return ExitUsage
	}
	cmd := args[0]
	rest := args[1:]
	switch cmd {
	case "sample":
		return cmdSample(rest, env)
	case "arclength", "arc-length", "len":
		return cmdArcLength(rest, env)
	case "curvature", "kappa":
		return cmdCurvature(rest, env)
	case "offset":
		return cmdOffset(rest, env)
	case "validate":
		return cmdValidate(rest, env)
	case "check", "verify":
		return cmdCheck(rest, env)
	case "help", "--help", "-h":
		writeUsage(env.Stdout)
		return ExitOK
	case "version", "--version":
		fmt.Fprintln(env.Stdout, "bezier-curv 1.0.0")
		return ExitOK
	default:
		fmt.Fprintf(env.Stderr, "bezier-curv: unknown command %q\n", cmd)
		writeUsage(env.Stderr)
		return ExitUsage
	}
}

func fail(env Env, format string, a ...any) int {
	fmt.Fprintf(env.Stderr, "bezier-curv: error: "+format+"\n", a...)
	return ExitError
}

type nilReader struct{}

func (nilReader) Read([]byte) (int, error) { return 0, io.EOF }
