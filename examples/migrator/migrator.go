package main

import (
	"context"
	"fmt"
	"os"

	"go.segfaultmedaddy.com/komandir"
)

type FlagSet struct {
	MigrationDirectory string `name:"migration-dir" alias:"dir" desc:"migration directory"`
	Verbose            bool   `name:"verbose" alias:"v" desc:"verbose output"`
	Version            bool   `name:"version" alias:"V" desc:"version information"`
}

type ArgSet struct {
	Direction string `name:"direction" desc:"migration direction"`
}

var cmd = &komandir.Command[FlagSet, ArgSet]{
	Name: "migrate",
	Action: func(ctx context.Context, cmd *komandir.Command[FlagSet, ArgSet]) error {
		fmt.Println("verbose", cmd.Flags.Verbose) // Verbose is a boolean flag indicating whether verbose output should be enabled.
		fmt.Println("direction", cmd..Direction)

		return nil
	},
}

func init() {
	cmd.AddCommand(cmd)
}

func main() {
	ctx := context.Background()

	err := cmd.Exec(ctx, os.Args[1:]...)
	if err != nil {
		panic(err)
	}
}
