package jolokia2

import (
	"github.com/masami10/rush"
	"github.com/masami10/rush/plugins/inputs"
)

func init() {
	inputs.Add("jolokia2_agent", func() rush.Input {
		return &JolokiaAgent{
			Metrics:               []MetricConfig{},
			DefaultFieldSeparator: ".",
		}
	})
	inputs.Add("jolokia2_proxy", func() rush.Input {
		return &JolokiaProxy{
			Metrics:               []MetricConfig{},
			DefaultFieldSeparator: ".",
		}
	})
}
