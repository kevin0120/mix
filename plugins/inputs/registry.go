package inputs

import "github.com/masami10/rush"

type Creator func() rush.Input

var Inputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Inputs[name] = creator
}
