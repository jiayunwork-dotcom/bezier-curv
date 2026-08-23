// bezier-curv 是三次 Bézier 弧长、曲率与法向等距偏移核算命令行工具。
// 求解主逻辑在 internal/，本文件只做子命令分发接线。
package main

import (
	"os"

	"bezier-curv/internal/cli"
)

func main() {
	env := cli.Env{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
	os.Exit(cli.Run(os.Args[1:], env))
}
