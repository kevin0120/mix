package discard

import (
	"github.com/masami10/rush"
	"github.com/masami10/rush/plugins/outputs"
)

type Discard struct{}

func (d *Discard) Connect() error                        { return nil }
func (d *Discard) Close() error                          { return nil }
func (d *Discard) SampleConfig() string                  { return "" }
func (d *Discard) Description() string                   { return "Send metrics to nowhere at all" }
func (d *Discard) Write(metrics []rush.Metric) error { return nil }

func init() {
	outputs.Add("discard", func() rush.Output { return &Discard{} })
}
