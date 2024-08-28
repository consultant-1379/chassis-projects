FROM armdocker.rnd.ericsson.se/proj-ipworks/rhel-tools:v1
RUN  mkdir -p /etc/nrfagent/
COPY ./nrfagent_common/build/etc /etc/nrfagent/
COPY ./nrfagent_discovery/build/config/cache-index.json /
COPY ./nrfagent_discovery/build/nrfagent /bin/nrfagentdisc

ENTRYPOINT ["/etc/nrfagent/scripts/entrypoint.sh"]
