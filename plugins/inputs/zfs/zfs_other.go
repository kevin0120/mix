// +build !linux,!freebsd

package zfs

import (
	"github.com/masami10/rush"
	"github.com/masami10/rush/plugins/inputs"
)

func (z *Zfs) Gather(acc rush.Accumulator) error {
	return nil
}

func init() {
	inputs.Add("zfs", func() rush.Input {
		return &Zfs{}
	})
}
