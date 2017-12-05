# virt-collectd-exporter

A prometheus exporter for [collectd](https://collectd.org/).
Like [another exporter](https://github.com/prometheus/collectd_exporter)
It accepts collectd's
[binary network protocol](https://collectd.org/wiki/index.php/Binary_protocol)
or
[metrics in JSON format](https://collectd.org/wiki/index.php/Plugin:Write_HTTP),
as source input, and expose the metrics for consumption by Prometheus server.

## motivation

The [collectd_exporter](https://github.com/prometheus/collectd_exporter) exposes the collectd
metrics to prometheus, with a simple mapping of collectd metric attributes to prometheus names
and labels. For a large set of use cases, this isn't an issue.
However, if you happen to use collectd to gather metrics and if you don't want to expose this detail,
this design choice becomes a limit.

The collectd [write_prometheus plugin](https://collectd.org/wiki/index.php/Plugin:Write_Prometheus)
intentionally follows the same naming rules of the `collectd_exporter`, in order to be a drop-in replacement.

The purpose of `virt-collectd-exporter` is to fill this very specific gap, and to allow flexibility
about the mapping between the name of collectd metrics and prometheus metrics.

## status

pre-alpha. All the infrastrucure is set, and we are now focusing on the name mapping module, the
very core of the `virt-collectd-exporter`.


## next steps

* nameconv API (figure out how to configure)
* comprehensive test suite


## alternatives

If you don't want or don't need the added flexibility that `virt-collectd-exporter` provides, you
may want to use the
Promethues' [collectd_exporter](https://github.com/prometheus/collectd_exporter)
or maybe the
The collectd [write_prometheus plugin](https://collectd.org/wiki/index.php/Plugin:Write_Prometheus)


## run it
```
docker run \
	${OPTIONS} \
	-p 9103:9103 \
	--tmpfs /run \
	--tmpfs /tmp \
	-v /sys/fs/cgroup:/sys/fs/cgroup:ro \
	-v /var/run/libvirt/libvirt-sock-ro:/var/lib/libvirt/libvirt-sock-ro:ro \
	${IMAGE_NAME}
```
