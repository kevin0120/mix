package aggregators

import "github.com/masami10/rush"

type Creator func() rush.Aggregator

var Aggregators = map[string]Creator{}

func Add(name string, creator Creator) {
	Aggregators[name] = creator
}
