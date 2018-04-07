package printer

import (
	"fmt"

	"github.com/masami10/rush"
	"github.com/masami10/rush/plugins/processors"
)

type Printer struct {
}

var sampleConfig = `
`

func (p *Printer) SampleConfig() string {
	return sampleConfig
}

func (p *Printer) Description() string {
	return "Print all metrics that pass through this filter."
}

func (p *Printer) Apply(in ...rush.Metric) []rush.Metric {
	for _, metric := range in {
		fmt.Println(metric.String())
	}
	return in
}

func init() {
	processors.Add("printer", func() rush.Processor {
		return &Printer{}
	})
}
