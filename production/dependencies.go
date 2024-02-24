package production

import "github.com/cbsinteractive/jakewan/distill/cmd"

func NewDependencies() cmd.Dependencies {
	return &dependencies{}
}

type dependencies struct{}
