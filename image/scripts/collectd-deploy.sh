#!/bin/sh

# HACK around missing PassEnvironment
for line in $( /usr/bin/strings /proc/1/environ | grep COLLECTD_ ); do
	eval "export $line";
done
# HACK around missing PassEnvironment

# HACK around /run tmpfs mismounted
# it seems that /run is mounted out of order, or one time too much. Until we find out why, let's workaround it
if [ ! -d /var/run/libvirt ]; then
	/bin/mkdir -p /var/run/libvirt
fi
if [ ! -S /var/run/libvirt/libvirt-sock ] && [ -S /var/lib/libvirt/libvirt-sock ]; then
	ln -s /var/lib/libvirt/libvirt-sock /var/run/libvirt/libvirt-sock
fi
if [ ! -S /var/run/libvirt/libvirt-sock-ro ] && [ -S /var/lib/libvirt/libvirt-sock-ro ]; then
	ln -s /var/lib/libvirt/libvirt-sock-ro /var/run/libvirt/libvirt-sock-ro
fi
# HACK around /run tmpfs mismounted

echo "# disabled: [${COLLECTD_GR_LOC}] -> [${COLLECTD_GR_HOST}:${COLLECTD_GR_PORT}]" > /etc/collectd.d/write_graphite.conf

if [ -n "${COLLECTD_GR_LOC}" ] && [ -n "${COLLECTD_GR_HOST}" ] && [ -n "${COLLECTD_GR_PORT}" ]; then
	cat > /etc/collectd.d/write_graphite.conf <<EOF
LoadPlugin write_graphite
<Plugin write_graphite>
  <Node "${COLLECTD_GR_LOC}">
    Host "${COLLECTD_GR_HOST}"
    Port "${COLLECTD_GR_PORT}"
    Protocol "tcp"
    ReconnectInterval 0
    LogSendErrors true
    Prefix "collectd."
    EscapeCharacter "_"
  </Node>
</Plugin>
EOF
fi
