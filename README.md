# rush
MasterPC 控制器与后端数据接入层
# Rush [![Circle CI](https://circleci.com/gh/masami10/rush.svg?style=svg)](https://circleci.com/gh/masami10/rush) [![Docker pulls](https://img.shields.io/docker/pulls/library/rush.svg)](https://hub.docker.com/_/rush/)

Rush is an agent written in Go for collecting, processing, aggregating,
and writing metrics.

Design goals are to have a minimal memory footprint with a plugin system so
that developers in the community can easily add support for collecting metrics
from local or remote services.

Rush is plugin-driven and has the concept of 4 distinct plugins:

1. [Input Plugins](#input-plugins) collect metrics from the system, services, or 3rd party APIs
2. [Processor Plugins](#processor-plugins) transform, decorate, and/or filter metrics
3. [Aggregator Plugins](#aggregator-plugins) create aggregate metrics (e.g. mean, min, max, quantiles, etc.)
4. [Output Plugins](#output-plugins) write metrics to various destinations

For more information on Processor and Aggregator plugins please [read this](./docs/AGGREGATORS_AND_PROCESSORS.md).

New plugins are designed to be easy to contribute,
we'll eagerly accept pull
requests and will manage the set of plugins that Rush supports.

## Contributing

There are many ways to contribute:
- Fix and [report bugs](https://github.com/masami10/rush/issues/new)
- [Improve documentation](https://github.com/masami10/rush/issues?q=is%3Aopen+label%3Adocumentation)
- [Review code and feature proposals](https://github.com/masami10/rush/pulls)
- Answer questions on github and on the [Community Site](https://community.influxdata.com/)
- [Contribute plugins](CONTRIBUTING.md)

## Installation:

You can download the binaries directly from the [downloads](https://www.influxdata.com/downloads) page
or from the [releases](https://github.com/masami10/rush/releases) section.

### Ansible Role:

Ansible role: https://github.com/rossmcdonald/rush

### From Source:

Rush requires golang version 1.8+, the Makefile requires GNU make.

Dependencies are managed with [gdm](https://github.com/sparrc/gdm),
which is installed by the Makefile if you don't have it already.

1. [Install Go](https://golang.org/doc/install)
2. [Setup your GOPATH](https://golang.org/doc/code.html#GOPATH)
3. Run `go get -d github.com/masami10/rush`
4. Run `cd $GOPATH/src/github.com/masami10/rush`
5. Run `make`

### Nightly Builds

These builds are generated from the master branch:
- [rush_nightly_amd64.deb](https://dl.influxdata.com/rush/nightlies/rush_nightly_amd64.deb)
- [rush_nightly_arm64.deb](https://dl.influxdata.com/rush/nightlies/rush_nightly_arm64.deb)
- [rush-nightly.arm64.rpm](https://dl.influxdata.com/rush/nightlies/rush-nightly.arm64.rpm)
- [rush_nightly_armel.deb](https://dl.influxdata.com/rush/nightlies/rush_nightly_armel.deb)
- [rush-nightly.armel.rpm](https://dl.influxdata.com/rush/nightlies/rush-nightly.armel.rpm)
- [rush_nightly_armhf.deb](https://dl.influxdata.com/rush/nightlies/rush_nightly_armhf.deb)
- [rush-nightly.armv6hl.rpm](https://dl.influxdata.com/rush/nightlies/rush-nightly.armv6hl.rpm)
- [rush-nightly_freebsd_amd64.tar.gz](https://dl.influxdata.com/rush/nightlies/rush-nightly_freebsd_amd64.tar.gz)
- [rush-nightly_freebsd_i386.tar.gz](https://dl.influxdata.com/rush/nightlies/rush-nightly_freebsd_i386.tar.gz)
- [rush_nightly_i386.deb](https://dl.influxdata.com/rush/nightlies/rush_nightly_i386.deb)
- [rush-nightly.i386.rpm](https://dl.influxdata.com/rush/nightlies/rush-nightly.i386.rpm)
- [rush-nightly_linux_amd64.tar.gz](https://dl.influxdata.com/rush/nightlies/rush-nightly_linux_amd64.tar.gz)
- [rush-nightly_linux_arm64.tar.gz](https://dl.influxdata.com/rush/nightlies/rush-nightly_linux_arm64.tar.gz)
- [rush-nightly_linux_armel.tar.gz](https://dl.influxdata.com/rush/nightlies/rush-nightly_linux_armel.tar.gz)
- [rush-nightly_linux_armhf.tar.gz](https://dl.influxdata.com/rush/nightlies/rush-nightly_linux_armhf.tar.gz)
- [rush-nightly_linux_i386.tar.gz](https://dl.influxdata.com/rush/nightlies/rush-nightly_linux_i386.tar.gz)
- [rush-nightly_linux_s390x.tar.gz](https://dl.influxdata.com/rush/nightlies/rush-nightly_linux_s390x.tar.gz)
- [rush_nightly_s390x.deb](https://dl.influxdata.com/rush/nightlies/rush_nightly_s390x.deb)
- [rush-nightly.s390x.rpm](https://dl.influxdata.com/rush/nightlies/rush-nightly.s390x.rpm)
- [rush-nightly_windows_amd64.zip](https://dl.influxdata.com/rush/nightlies/rush-nightly_windows_amd64.zip)
- [rush-nightly_windows_i386.zip](https://dl.influxdata.com/rush/nightlies/rush-nightly_windows_i386.zip)
- [rush-nightly.x86_64.rpm](https://dl.influxdata.com/rush/nightlies/rush-nightly.x86_64.rpm)
- [rush-static-nightly_linux_amd64.tar.gz](https://dl.influxdata.com/rush/nightlies/rush-static-nightly_linux_amd64.tar.gz)

## How to use it:

See usage with:

```
./rush --help
```

#### Generate a rush config file:

```
./rush config > rush.conf
```

#### Generate config with only cpu input & influxdb output plugins defined:

```
./rush --input-filter cpu --output-filter influxdb config
```

#### Run a single rush collection, outputing metrics to stdout:

```
./rush --config rush.conf --test
```

#### Run rush with all plugins defined in config file:

```
./rush --config rush.conf
```

#### Run rush, enabling the cpu & memory input, and influxdb output plugins:

```
./rush --config rush.conf --input-filter cpu:mem --output-filter influxdb
```


## Configuration

See the [configuration guide](docs/CONFIGURATION.md) for a rundown of the more advanced
configuration options.

## Input Plugins

* [aerospike](./plugins/inputs/aerospike)
* [amqp_consumer](./plugins/inputs/amqp_consumer) (rabbitmq)
* [apache](./plugins/inputs/apache)
* [aws cloudwatch](./plugins/inputs/cloudwatch)
* [bcache](./plugins/inputs/bcache)
* [bond](./plugins/inputs/bond)
* [cassandra](./plugins/inputs/cassandra)
* [ceph](./plugins/inputs/ceph)
* [cgroup](./plugins/inputs/cgroup)
* [chrony](./plugins/inputs/chrony)
* [consul](./plugins/inputs/consul)
* [conntrack](./plugins/inputs/conntrack)
* [couchbase](./plugins/inputs/couchbase)
* [couchdb](./plugins/inputs/couchdb)
* [DC/OS](./plugins/inputs/dcos)
* [disque](./plugins/inputs/disque)
* [dmcache](./plugins/inputs/dmcache)
* [dns query time](./plugins/inputs/dns_query)
* [docker](./plugins/inputs/docker)
* [dovecot](./plugins/inputs/dovecot)
* [elasticsearch](./plugins/inputs/elasticsearch)
* [exec](./plugins/inputs/exec) (generic executable plugin, support JSON, influx, graphite and nagios)
* [fail2ban](./plugins/inputs/fail2ban)
* [filestat](./plugins/inputs/filestat)
* [fluentd](./plugins/inputs/fluentd)
* [graylog](./plugins/inputs/graylog)
* [haproxy](./plugins/inputs/haproxy)
* [hddtemp](./plugins/inputs/hddtemp)
* [http_response](./plugins/inputs/http_response)
* [httpjson](./plugins/inputs/httpjson) (generic JSON-emitting http service plugin)
* [internal](./plugins/inputs/internal)
* [influxdb](./plugins/inputs/influxdb)
* [interrupts](./plugins/inputs/interrupts)
* [ipmi_sensor](./plugins/inputs/ipmi_sensor)
* [iptables](./plugins/inputs/iptables)
* [jolokia](./plugins/inputs/jolokia) (deprecated, use [jolokia2](./plugins/inputs/jolokia2))
* [jolokia2](./plugins/inputs/jolokia2)
* [kapacitor](./plugins/inputs/kapacitor)
* [kubernetes](./plugins/inputs/kubernetes)
* [leofs](./plugins/inputs/leofs)
* [lustre2](./plugins/inputs/lustre2)
* [mailchimp](./plugins/inputs/mailchimp)
* [memcached](./plugins/inputs/memcached)
* [mesos](./plugins/inputs/mesos)
* [minecraft](./plugins/inputs/minecraft)
* [mongodb](./plugins/inputs/mongodb)
* [mysql](./plugins/inputs/mysql)
* [net_response](./plugins/inputs/net_response)
* [nginx](./plugins/inputs/nginx)
* [nginx_plus](./plugins/inputs/nginx_plus)
* [nsq](./plugins/inputs/nsq)
* [nstat](./plugins/inputs/nstat)
* [ntpq](./plugins/inputs/ntpq)
* [openldap](./plugins/inputs/openldap)
* [opensmtpd](./plugins/inputs/opensmtpd)
* [pf](./plugins/inputs/pf)
* [phpfpm](./plugins/inputs/phpfpm)
* [phusion passenger](./plugins/inputs/passenger)
* [ping](./plugins/inputs/ping)
* [postfix](./plugins/inputs/postfix)
* [postgresql_extensible](./plugins/inputs/postgresql_extensible)
* [postgresql](./plugins/inputs/postgresql)
* [powerdns](./plugins/inputs/powerdns)
* [procstat](./plugins/inputs/procstat)
* [prometheus](./plugins/inputs/prometheus) (can be used for [Caddy server](./plugins/inputs/prometheus/README.md#usage-for-caddy-http-server))
* [puppetagent](./plugins/inputs/puppetagent)
* [rabbitmq](./plugins/inputs/rabbitmq)
* [raindrops](./plugins/inputs/raindrops)
* [redis](./plugins/inputs/redis)
* [rethinkdb](./plugins/inputs/rethinkdb)
* [riak](./plugins/inputs/riak)
* [salesforce](./plugins/inputs/salesforce)
* [sensors](./plugins/inputs/sensors)
* [smart](./plugins/inputs/smart)
* [snmp](./plugins/inputs/snmp)
* [snmp_legacy](./plugins/inputs/snmp_legacy)
* [solr](./plugins/inputs/solr)
* [sql server](./plugins/inputs/sqlserver) (microsoft)
* [teamspeak](./plugins/inputs/teamspeak)
* [tomcat](./plugins/inputs/tomcat)
* [twemproxy](./plugins/inputs/twemproxy)
* [unbound](./plugins/input/unbound)
* [varnish](./plugins/inputs/varnish)
* [zfs](./plugins/inputs/zfs)
* [zookeeper](./plugins/inputs/zookeeper)
* [win_perf_counters](./plugins/inputs/win_perf_counters) (windows performance counters)
* [win_services](./plugins/inputs/win_services)
* [sysstat](./plugins/inputs/sysstat)
* [system](./plugins/inputs/system)
    * cpu
    * mem
    * net
    * netstat
    * disk
    * diskio
    * swap
    * processes
    * kernel (/proc/stat)
    * kernel (/proc/vmstat)
    * linux_sysctl_fs (/proc/sys/fs)

Rush can also collect metrics via the following service plugins:

* [http_listener](./plugins/inputs/http_listener)
* [kafka_consumer](./plugins/inputs/kafka_consumer)
* [mqtt_consumer](./plugins/inputs/mqtt_consumer)
* [nats_consumer](./plugins/inputs/nats_consumer)
* [nsq_consumer](./plugins/inputs/nsq_consumer)
* [logparser](./plugins/inputs/logparser)
* [statsd](./plugins/inputs/statsd)
* [socket_listener](./plugins/inputs/socket_listener)
* [tail](./plugins/inputs/tail)
* [tcp_listener](./plugins/inputs/socket_listener)
* [udp_listener](./plugins/inputs/socket_listener)
* [webhooks](./plugins/inputs/webhooks)
  * [filestack](./plugins/inputs/webhooks/filestack)
  * [github](./plugins/inputs/webhooks/github)
  * [mandrill](./plugins/inputs/webhooks/mandrill)
  * [papertrail](./plugins/inputs/webhooks/papertrail)
  * [particle](./plugins/inputs/webhooks/particle)
  * [rollbar](./plugins/inputs/webhooks/rollbar)
* [zipkin](./plugins/inputs/zipkin)

Rush is able to parse the following input data formats into metrics, these
formats may be used with input plugins supporting the `data_format` option:

* [InfluxDB Line Protocol](./docs/DATA_FORMATS_INPUT.md#influx)
* [JSON](./docs/DATA_FORMATS_INPUT.md#json)
* [Graphite](./docs/DATA_FORMATS_INPUT.md#graphite)
* [Value](./docs/DATA_FORMATS_INPUT.md#value)
* [Nagios](./docs/DATA_FORMATS_INPUT.md#nagios)
* [Collectd](./docs/DATA_FORMATS_INPUT.md#collectd)

## Processor Plugins

* [printer](./plugins/processors/printer)

## Aggregator Plugins

* [basicstats](./plugins/aggregators/basicstats)
* [minmax](./plugins/aggregators/minmax)
* [histogram](./plugins/aggregators/histogram)

## Output Plugins

* [influxdb](./plugins/outputs/influxdb)
* [amon](./plugins/outputs/amon)
* [amqp](./plugins/outputs/amqp) (rabbitmq)
* [aws kinesis](./plugins/outputs/kinesis)
* [aws cloudwatch](./plugins/outputs/cloudwatch)
* [cratedb](./plugins/outputs/cratedb)
* [datadog](./plugins/outputs/datadog)
* [discard](./plugins/outputs/discard)
* [elasticsearch](./plugins/outputs/elasticsearch)
* [file](./plugins/outputs/file)
* [graphite](./plugins/outputs/graphite)
* [graylog](./plugins/outputs/graylog)
* [instrumental](./plugins/outputs/instrumental)
* [kafka](./plugins/outputs/kafka)
* [librato](./plugins/outputs/librato)
* [mqtt](./plugins/outputs/mqtt)
* [nats](./plugins/outputs/nats)
* [nsq](./plugins/outputs/nsq)
* [opentsdb](./plugins/outputs/opentsdb)
* [prometheus](./plugins/outputs/prometheus_client)
* [riemann](./plugins/outputs/riemann)
* [riemann_legacy](./plugins/outputs/riemann_legacy)
* [socket_writer](./plugins/outputs/socket_writer)
* [tcp](./plugins/outputs/socket_writer)
* [udp](./plugins/outputs/socket_writer)
* [wavefront](./plugins/outputs/wavefront)
