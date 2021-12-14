package app

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	commandPrefix = ".rt"
)

type tCommand []string

func toCommand(s string) (tCommand, bool) {
	if !strings.HasPrefix(s, commandPrefix) {
		return nil, false
	}
	ss := strings.Fields(s)
	if len(ss) <= 1 || ss[0] != commandPrefix {
		return nil, false
	}
	return tCommand(ss[1:]), true
}

func (command tCommand) match(expect ...string) bool {
	if len(command) != len(expect) {
		return false
	}
	for i := range command {
		if command[i] != expect[i] && expect[i] != "*" {
			return false
		}
	}
	return true
}

func (command tCommand) fetch(i int) string {
	if len(command) <= i {
		slackWarning.Post(errors.New(fmt.Sprintf(
			"slice size is over\ncommand: %+v\ni: %v",
			command,
			i,
		)))
		return ""
	}
	return command[i]
}
