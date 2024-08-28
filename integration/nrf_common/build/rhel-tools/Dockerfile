FROM armdocker.rnd.ericsson.se/rhel7.4 as rheltools-builder
MAINTAINER ezganlu/ericsson

# Copy reposiroty file into base image, so the tools can be installed.
COPY ./build/rhel-tools/ericsson-yum-server.repo /etc/yum.repos.d/ericsson-yum-server.repo

RUN yum install -y gcc-c++ autoconf automake libtool make openssl-devel libxml2-devel python-devel systemd-devel c-ares-devel jansson wget

RUN cd /tmp && \
    wget http://www.digip.org/jansson/releases/jansson-2.9.tar.gz -O /tmp/jansson.tar.gz && \
    mkdir -p /tmp/jansson && \
    tar -zxvf /tmp/jansson.tar.gz -C /tmp/jansson --strip-components=1 && \
    cd /tmp/jansson && ./configure --libdir=/usr/lib64 && \
    make && make install && \
    cd /tmp/ && \
    wget http://mirror.centos.org/centos/7/os/x86_64/Packages/libevent-2.0.21-4.el7.x86_64.rpm && \
    wget http://mirror.centos.org/centos/7/os/x86_64/Packages/libevent-devel-2.0.21-4.el7.x86_64.rpm && \
    wget http://dl.fedoraproject.org/pub/epel/7/x86_64/Packages/j/jemalloc-3.6.0-1.el7.x86_64.rpm && \
    wget http://dl.fedoraproject.org/pub/epel/7/x86_64/Packages/j/jemalloc-devel-3.6.0-1.el7.x86_64.rpm && \
    wget http://dl.fedoraproject.org/pub/epel/7/x86_64/Packages/l/libev-4.15-3.el7.x86_64.rpm && \
    wget http://dl.fedoraproject.org/pub/epel/7/x86_64/Packages/l/libev-devel-4.15-3.el7.x86_64.rpm && \
    rpm -hiv *.rpm

ENV NGHTTP2_version=1.28.0
RUN cd /tmp && \
    wget -O nghttp2.tar.gz https://github.com/nghttp2/nghttp2/releases/download/v${NGHTTP2_version}/nghttp2-${NGHTTP2_version}.tar.gz && \ 
    tar zxvf nghttp2.tar.gz && \
    cd nghttp2-${NGHTTP2_version} && \
    ./configure --enable-app && \
    make && make install  
    
RUN mkdir -p /tmp/usr/lib64 && mkdir -p /tmp/usr/bin && \
    cp /usr/local/lib/libnghttp2.so.14 /tmp/usr/lib64 && \
    cp /usr/lib64/libjemalloc.so.1 /tmp/usr/lib64 && \
    cp /usr/lib64/libev.so.4 /tmp/usr/lib64 && \
    cp /usr/lib64/libcares.so.2 /tmp/usr/lib64 && \
    cp /usr/lib64/libsystemd.so.0 /tmp/usr/lib64 && \
    cp /usr/lib64/libjansson.so.4 /tmp/usr/lib64 && \
    cp /usr/lib64/libdw.so.1 /tmp/usr/lib64 && \
    cp /usr/local/bin/nghttp /tmp/usr/bin && \
    cp /etc/yum.repos.d/ericsson-yum-server.repo /tmp/usr

RUN cp /usr/bin/uuidgen /tmp/usr/bin 


# build runtime image
FROM armdocker.rnd.ericsson.se/rhel7-atomic:7.4

COPY --from=rheltools-builder /tmp/usr/ /usr/

RUN set -eux; \
    mv /usr/ericsson-yum-server.repo /etc/yum.repos.d/; \
    microdnf install -y tar tcpdump curl wget --nodocs; \
    # mv tcpdump to /usr/bin to solve privilege issue
    mv /usr/sbin/tcpdump /usr/bin/tcpdump; \
    microdnf -y clean all; \
    echo "wget -O /usr/bin/busybox https://busybox.net/downloads/binaries/1.28.1-defconfig-multiarch/busybox-x86_64; chmod a+x /usr/bin/busybox" > /usr/bin/busybox_dl && chmod a+x /usr/bin/busybox_dl; \
	wget -O /usr/bin/jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 && chmod a+x /usr/bin/jq; \
    ln -s /usr/bin/busybox /usr/bin/ifconfig; \
    ln -s /usr/bin/busybox /usr/bin/netstat


