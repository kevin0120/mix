// +build !linux

package dmcache

import (
	"github.com/masami10/rush"
)

func (c *DMCache) Gather(acc rush.Accumulator) error {
	return nil
}

func dmSetupStatus() ([]string, error) {
	return []string{}, nil
}
