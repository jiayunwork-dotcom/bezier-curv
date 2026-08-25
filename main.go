package main

import (
	"fmt"
	"os"

	"bezier-curv/internal/cli"
	"bezier-curv/internal/server"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		addr := ":8080"
		if len(os.Args) > 2 {
			addr = os.Args[2]
		}
		fmt.Fprintf(os.Stderr, "bezier-curv: listening on %s\n", addr)
		srv := server.New(addr)
		if err := srv.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "bezier-curv: %v\n", err)
			os.Exit(1)
		}
		return
	}
	env := cli.Env{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
	os.Exit(cli.Run(os.Args[1:], env))
}
