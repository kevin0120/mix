package xml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/masami10/rush"
	"github.com/masami10/rush/metric"
)

type XMLParser struct {
	MetricName  string
	TagKeys     []string
	DefaultTags map[string]string
}

func (p *XMLParser) parseArray(buf []byte) ([]rush.Metric, error) {
	metrics := make([]rush.Metric, 0)

	var xmlOut []map[string]interface{}
	err := xml.Unmarshal(buf, &xmlOut)
	if err != nil {
		err = fmt.Errorf("unable to parse out as XML Array, %s", err)
		return nil, err
	}
	for _, item := range xmlOut {
		metrics, err = p.parseObject(metrics, item)
	}
	return metrics, nil
}

func (p *XMLParser) parseObject(metrics []rush.Metric, xmlOut map[string]interface{}) ([]rush.Metric, error) {

	tags := make(map[string]string)
	for k, v := range p.DefaultTags {
		tags[k] = v
	}

	for _, tag := range p.TagKeys {
		switch v := xmlOut[tag].(type) {
		case string:
			tags[tag] = v
		case bool:
			tags[tag] = strconv.FormatBool(v)
		case float64:
			tags[tag] = strconv.FormatFloat(v, 'f', -1, 64)
		}
		delete(xmlOut, tag)
	}

	f := XMLFlattener{}
	err := f.FlattenXML("", xmlOut)
	if err != nil {
		return nil, err
	}

	metric, err := metric.New(p.MetricName, tags, f.Fields, time.Now().UTC())

	if err != nil {
		return nil, err
	}
	return append(metrics, metric), nil
}

func (p *XMLParser) Parse(buf []byte) ([]rush.Metric, error) {
	buf = bytes.TrimSpace(buf)
	if len(buf) == 0 {
		return make([]rush.Metric, 0), nil
	}

	if !isarray(buf) {
		metrics := make([]rush.Metric, 0)
		var xmlOut map[string]interface{}
		err := xml.Unmarshal(buf, &xmlOut)
		if err != nil {
			err = fmt.Errorf("unable to parse out as XML, %s", err)
			return nil, err
		}
		return p.parseObject(metrics, xmlOut)
	}
	return p.parseArray(buf)
}

func (p *XMLParser) ParseLine(line string) (rush.Metric, error) {
	metrics, err := p.Parse([]byte(line + "\n"))

	if err != nil {
		return nil, err
	}

	if len(metrics) < 1 {
		return nil, fmt.Errorf("Can not parse the line: %s, for data format: influx ", line)
	}

	return metrics[0], nil
}

func (p *XMLParser) SetDefaultTags(tags map[string]string) {
	p.DefaultTags = tags
}

type XMLFlattener struct {
	Fields map[string]interface{}
}

// FlattenXML flattens nested maps/interfaces into a fields map (ignoring bools and string)
func (f *XMLFlattener) FlattenXML(
	fieldname string,
	v interface{}) error {
	if f.Fields == nil {
		f.Fields = make(map[string]interface{})
	}
	return f.FullFlattenXML(fieldname, v, false, false)
}

// FullFlattenXML flattens nested maps/interfaces into a fields map (including bools and string)
func (f *XMLFlattener) FullFlattenXML(
	fieldname string,
	v interface{},
	convertString bool,
	convertBool bool,
) error {
	if f.Fields == nil {
		f.Fields = make(map[string]interface{})
	}
	fieldname = strings.Trim(fieldname, "_")
	switch t := v.(type) {
	case map[string]interface{}:
		for k, v := range t {
			err := f.FullFlattenXML(fieldname+"_"+k+"_", v, convertString, convertBool)
			if err != nil {
				return err
			}
		}
	case []interface{}:
		for i, v := range t {
			k := strconv.Itoa(i)
			err := f.FullFlattenXML(fieldname+"_"+k+"_", v, convertString, convertBool)
			if err != nil {
				return nil
			}
		}
	case float64:
		f.Fields[fieldname] = t
	case string:
		if convertString {
			f.Fields[fieldname] = v.(string)
		} else {
			return nil
		}
	case bool:
		if convertBool {
			f.Fields[fieldname] = v.(bool)
		} else {
			return nil
		}
	case nil:
		return nil
	default:
		return fmt.Errorf("XML Flattener: got unexpected type %T with value %v (%s)",
			t, t, fieldname)
	}
	return nil
}

func isarray(buf []byte) bool {
	ia := bytes.IndexByte(buf, '[')
	ib := bytes.IndexByte(buf, '{')
	if ia > -1 && ia < ib {
		return true
	} else {
		return false
	}
}
