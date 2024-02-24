package cmd

import (
	"fmt"
	"slices"
	"strings"
)

func newRequiredArgumentMissingError(n argName) error {
	arg := string(n)
	if !strings.HasPrefix(arg, "-") {
		arg = "-" + arg
	}
	return requiredArgumentMissingError{arg}
}

type requiredArgumentMissingError struct {
	arg string
}

func (e requiredArgumentMissingError) Error() string {
	return fmt.Sprintf("required argument missing %s", e.arg)
}

func newSubcommandExpectedError(runners []runner) error {
	names := []string{}
	for _, r := range runners {
		names = append(names, r.name())
	}
	slices.Sort(names)
	return subcommandExpectedError{oneOf: names}
}

type subcommandExpectedError struct {
	oneOf []string
}

// Error implements error.
func (e subcommandExpectedError) Error() string {
	if len(e.oneOf) < 1 {
		panic("no subcommands programmed")
	}
	return fmt.Sprintf("subcommand expected (%s)", strings.Join(e.oneOf, ","))
}

func newUnexpectedSubcommandError(subcommand string) error {
	return unexpectedSubcommandError{subcommand}
}

type unexpectedSubcommandError struct {
	unexpectedSubcommand string
}

// Error implements error.
func (e unexpectedSubcommandError) Error() string {
	return fmt.Sprintf("unexpected subcommand %s", e.unexpectedSubcommand)
}
