// +build !linux

package cgroup

import (
	"github.com/masami10/rush"
)

func (g *CGroup) Gather(acc rush.Accumulator) error {
	return nil
}
