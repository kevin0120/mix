package outputs

import (
	"github.com/masami10/rush"
)

type Creator func() rush.Output

var Outputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Outputs[name] = creator
}
