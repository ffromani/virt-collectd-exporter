FROM centos:7.4.1708
ENV container docker

LABEL maintainer Francesco Romani <fromani@redhat.com>

RUN cd /lib/systemd/system/sysinit.target.wants/; ls | grep -v systemd-tmpfiles-setup | xargs rm -f $1 \
rm -f /lib/systemd/system/multi-user.target.wants/*;\
rm -f /etc/systemd/system/*.wants/*;\
rm -f /lib/systemd/system/local-fs.target.wants/*; \
rm -f /lib/systemd/system/sockets.target.wants/*udev*; \
rm -f /lib/systemd/system/sockets.target.wants/*initctl*; \
rm -f /lib/systemd/system/basic.target.wants/*;\
rm -f /lib/systemd/system/anaconda.target.wants/*; \
rm -f /lib/systemd/system/plymouth*; \
rm -f /lib/systemd/system/systemd-update-utmp*;
RUN systemctl set-default multi-user.target
ENV init /lib/systemd/systemd

#ADD image/repos/opstools.repo /etc/yum.repos.d
RUN yum -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm && \
yum-config-manager --enable epel-testing

# HACK around missing PassEnvironment. We need binutils
RUN yum install -y collectd collectd-virt binutils && yum clean all

ADD image/configs/collectd.conf /etc/collectd.conf
ADD image/configs/virt.conf /etc/collectd.d/virt.conf
ADD image/configs/network.conf /etc/collectd.d/network.conf

RUN mkdir /etc/systemd/system/collectd.service.d/
ADD image/scripts/collectd-deps.conf /etc/systemd/system/collectd.service.d/
ADD image/scripts/collectd-deploy.service /etc/systemd/system/
ADD image/scripts/collectd-deploy.sh /usr/libexec/

ADD _output/virt-collectd-exporter /usr/sbin/
ADD image/scripts/virt-collectd-exporter.service /etc/systemd/system/

RUN systemctl enable collectd-deploy.service collectd.service virt-collectd-exporter.service

# https://developers.redhat.com/blog/2016/09/13/running-systemd-in-a-non-privileged-container/
STOPSIGNAL SIGRTMIN+3

ENTRYPOINT ["/sbin/init"]
