package flags

import (
	"flag"
	"fmt"
	"strings"
)

type TargetVariable[T any] interface {
	SetDefault(defaultValue T)
}

type targetVariable[T any] struct {
	name         string
	t            *T
	defaultValue T
	usage        string
}

// SetDefault implements TargetVariable.
func (t *targetVariable[T]) SetDefault(defaultValue T) {
	t.defaultValue = defaultValue
}

type FlagSet interface {
	StringVar(target *string, name string, usage string) TargetVariable[string]
	BoolVar(target *bool, name string, usage string) TargetVariable[bool]
	Parse(args []string) error
}

func NewFlagSet(name string, errorMode flag.ErrorHandling) FlagSet {
	return &flagSet{
		boolVars:   []TargetVariable[bool]{},
		stringVars: []TargetVariable[string]{},
	}
}

type flagSet struct {
	boolVars   []TargetVariable[bool]
	stringVars []TargetVariable[string]
}

// Parse implements FlagSet.
func (f *flagSet) Parse(args []string) error {
	fmt.Printf("Args: %s\n", args)
	for _, v := range f.stringVars {
		fmt.Printf("Flag: %v\n", v)
	}
	for _, v := range f.boolVars {
		fmt.Printf("Flag: %v\n", v)
	}
	fmt.Printf("Arg length: %d\n", len(args))
	for currentArgIdx, currentArg := range args {
		if strings.HasPrefix(currentArg, "-") {
			fmt.Printf("Current arg index: %d\n", currentArgIdx)
			currentArgName := strings.TrimLeft(currentArg, "-")
			fmt.Printf("Current arg: %s\n", currentArgName)
			if len(args) > currentArgIdx+1 {
				nextArg := args[currentArgIdx+1]
				fmt.Printf("Current arg value as string: %s\n", nextArg)
			}
		}
	}
	panic("unimplemented")
}

// BoolVar implements FlagSet.
func (f *flagSet) BoolVar(target *bool, name string, usage string) TargetVariable[bool] {
	t := &targetVariable[bool]{
		name:  name,
		t:     target,
		usage: usage,
	}
	f.boolVars = append(f.boolVars, t)
	return t
}

// StringVar implements FlagSet.
func (f *flagSet) StringVar(target *string, name string, usage string) TargetVariable[string] {
	t := &targetVariable[string]{
		name:  name,
		t:     target,
		usage: usage,
	}
	f.stringVars = append(f.stringVars, t)
	return t
}
