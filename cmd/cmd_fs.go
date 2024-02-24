package cmd

import "context"

type cmdFS struct {
	r    runner
	deps Dependencies
}

func newCmdFS(deps Dependencies) runner {
	return &cmdFS{deps: deps}
}

func (cmd *cmdFS) name() string {
	return "fs"
}

func (cmd *cmdFS) init(args []string) error {
	runners := []runner{
		newCmdFSDirsize(cmd.deps),
	}
	if len(args) < 1 {
		return newSubcommandExpectedError(runners)
	}
	return processSubcommand(&cmd.r, args, runners)
}

func (cmd *cmdFS) run(ctx context.Context) error {
	return cmd.r.run(ctx)
}
