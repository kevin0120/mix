package processors

import "github.com/masami10/rush"

type Creator func() rush.Processor

var Processors = map[string]Creator{}

func Add(name string, creator Creator) {
	Processors[name] = creator
}
