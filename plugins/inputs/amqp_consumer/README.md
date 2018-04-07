# AMQP Consumer Input Plugin

This plugin provides a consumer for use with AMQP 0-9-1, a promenent implementation of this protocol being [RabbitMQ](https://www.rabbitmq.com/).

Metrics are read from a topic exchange using the configured queue and binding_key.

Message payload should be formatted in one of the [Rush Data Formats](https://github.com/masami10/rush/blob/master/docs/DATA_FORMATS_INPUT.md).

For an introduction to AMQP see:
- https://www.rabbitmq.com/tutorials/amqp-concepts.html
- https://www.rabbitmq.com/getstarted.html

The following defaults are known to work with RabbitMQ:

```toml
# AMQP consumer plugin
[[inputs.amqp_consumer]]
  ## AMQP url
  url = "amqp://localhost:5672/influxdb"
  ## AMQP exchange
  exchange = "rush"
  ## AMQP queue name
  queue = "rush"
  ## Binding Key
  binding_key = "#"

  ## Controls how many messages the server will try to keep on the network
  ## for consumers before receiving delivery acks.
  #prefetch_count = 50

  ## Auth method. PLAIN and EXTERNAL are supported.
  ## Using EXTERNAL requires enabling the rabbitmq_auth_mechanism_ssl plugin as
  ## described here: https://www.rabbitmq.com/plugins.html
  # auth_method = "PLAIN"
  ## Optional SSL Config
  # ssl_ca = "/etc/rush/ca.pem"
  # ssl_cert = "/etc/rush/cert.pem"
  # ssl_key = "/etc/rush/key.pem"
  ## Use SSL but skip chain & host verification
  # insecure_skip_verify = false

  ## Data format to consume.
  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/masami10/rush/blob/master/docs/DATA_FORMATS_INPUT.md
  data_format = "influx"
```
