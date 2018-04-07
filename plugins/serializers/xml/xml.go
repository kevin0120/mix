package xml

import (
	exml "encoding/xml"
	"time"

	"github.com/masami10/rush"
)

type XmlSerializer struct {
	TimestampUnits time.Duration
}

func (s *XmlSerializer) Serialize(metric rush.Metric) ([]byte, error) {
	m := make(map[string]interface{})
	units_nanoseconds := s.TimestampUnits.Nanoseconds()
	// if the units passed in were less than or equal to zero,
	// then serialize the timestamp in seconds (the default)
	if units_nanoseconds <= 0 {
		units_nanoseconds = 1000000000
	}
	m["tags"] = metric.Tags()
	m["fields"] = metric.Fields()
	m["name"] = metric.Name()
	m["timestamp"] = metric.UnixNano() / units_nanoseconds
	serialized, err := exml.Marshal(m)
	if err != nil {
		return []byte{}, err
	}
	serialized = append(serialized, '\n')

	return serialized, nil
}
