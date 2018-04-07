package influx

import (
	"github.com/masami10/rush"
)

type InfluxSerializer struct {
}

func (s *InfluxSerializer) Serialize(m rush.Metric) ([]byte, error) {
	return m.Serialize(), nil
}
