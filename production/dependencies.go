package production

import "github.com/jakewan/distill/cmd"

func NewDependencies() cmd.Dependencies {
	return &dependencies{}
}

type dependencies struct{}
