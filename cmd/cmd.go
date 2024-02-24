package cmd

import (
	"context"
)

type Dependencies interface{}

type argName string

const (
	argNameStartingDir argName = "startingdir"
	argNameVerbose     argName = "verbose"
)

type argUsage string

const (
	argUsageStartingDir argUsage = "The starting directory of the operation"
	argUsageVerbose     argUsage = "Verbose output"
)

type runner interface {
	name() string
	init([]string) error
	run(ctx context.Context) error
}

// Execute takes the command line arguments and a set of injected dependencies
// and processes them to conduct the program's business logic.
func Execute(args []string, deps Dependencies) error {
	return rootHandler(args, deps)
}

func rootHandler(args []string, deps Dependencies) error {
	runners := []runner{
		newCmdFS(deps),
	}
	if len(args) < 1 {
		return newSubcommandExpectedError(runners)
	}
	subcommand := args[0]
	for _, r := range runners {
		if r.name() == subcommand {
			if err := r.init(args[1:]); err != nil {
				return err
			}
			if err := r.run(context.Background()); err != nil {
				return err
			}
			return nil
		}
	}
	return newUnexpectedSubcommandError(subcommand)
}

func processSubcommand(childRunner *runner, args []string, childRunners []runner) error {
	subcommand := args[0]
	for _, r := range childRunners {
		if r.name() == subcommand {
			if err := r.init(args[1:]); err != nil {
				return err
			}
			*childRunner = r
			return nil
		}
	}
	return newUnexpectedSubcommandError(subcommand)
}
