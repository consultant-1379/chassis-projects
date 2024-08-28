FROM armdocker.rnd.ericsson.se/proj-ipworks/alpine-tools:v1
MAINTAINER eqianbi/ericsson
COPY nrf-istio-check.sh /bin

ENTRYPOINT ["/bin/nrf-istio-check.sh"]