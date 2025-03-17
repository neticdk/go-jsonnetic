package jsonneticcli

import "github.com/neticdk/go-common/pkg/cli/cmd"

type Context struct {
	EC *cmd.ExecutionContext
}

func NewContext(ec *cmd.ExecutionContext) *Context {
	return &Context{
		EC: ec,
	}
}
