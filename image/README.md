= collectd =

```
docker run \
	${OPTIONS} \
	--tmpfs /run \
	--tmpfs /tmp \
	-v /sys/fs/cgroup:/sys/fs/cgroup:ro \
	-v /var/run/libvirt/libvirt-sock-ro:/var/lib/libvirt/libvirt-sock-ro:ro \
	-e COLLECTD_GR_LOC=myhost \
	-e COLLECTD_GR_HOST=carbon-host \
	-e COLLECTD_GR_PORT=2003 \
	centos/collectd
```
