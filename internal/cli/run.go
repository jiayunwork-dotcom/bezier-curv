// Package cli 提供 bezier-curv 的命令行入口：子命令分发、输入读取与输出格式化。
// 所有非法输入（文件缺失、控制点不合法、尖点求曲率、偏移被拒绝）都写到 stderr，
// 并以非零退出码结束；成功结果写到 stdout。
package cli

import (
	"fmt"
	"io"
)

// 退出码约定：0 成功；1 数据/校验/运行时错误；2 用法错误。
const (
	ExitOK    = 0
	ExitError = 1
	ExitUsage = 2
)

// Env 是命令执行的 I/O 上下文，便于在测试中注入。
type Env struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Run 执行命令行 args（不含程序名）并返回进程退出码。
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

// fail 把错误写到 stderr 并返回错误退出码。
func fail(env Env, format string, a ...any) int {
	fmt.Fprintf(env.Stderr, "bezier-curv: error: "+format+"\n", a...)
	return ExitError
}

// nilReader 是永不返回数据的空 reader。
type nilReader struct{}

func (nilReader) Read([]byte) (int, error) { return 0, io.EOF }
