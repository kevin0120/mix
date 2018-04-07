# NATS Consumer Input Plugin

The [NATS](http://www.nats.io/about/) consumer plugin reads from
specified NATS subjects and adds messages to InfluxDB. The plugin expects messages
in the [Rush Input Data Formats](https://github.com/masami10/rush/blob/master/docs/DATA_FORMATS_INPUT.md).
A [Queue Group](http://www.nats.io/documentation/concepts/nats-queueing/)
is used when subscribing to subjects so multiple instances of rush can read
from a NATS cluster in parallel.

## Configuration

```toml
# Read metrics from NATS subject(s)
[[inputs.nats_consumer]]
  ## urls of NATS servers
  servers = ["nats://localhost:4222"]
  ## Use Transport Layer Security
  secure = false
  ## subject(s) to consume
  subjects = ["rush"]
  ## name a queue group
  queue_group = "rush_consumers"
  ## Maximum number of metrics to buffer between collection intervals
  metric_buffer = 100000

  ## Data format to consume. 

  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/masami10/rush/blob/master/docs/DATA_FORMATS_INPUT.md
  data_format = "influx"
```
