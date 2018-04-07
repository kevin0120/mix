// +build !linux

package system

import (
	"github.com/masami10/rush"
	"github.com/masami10/rush/plugins/inputs"
)

type Kernel struct {
}

func (k *Kernel) Description() string {
	return "Get kernel statistics from /proc/stat"
}

func (k *Kernel) SampleConfig() string { return "" }

func (k *Kernel) Gather(acc rush.Accumulator) error {
	return nil
}

func init() {
	inputs.Add("kernel", func() rush.Input {
		return &Kernel{}
	})
}
