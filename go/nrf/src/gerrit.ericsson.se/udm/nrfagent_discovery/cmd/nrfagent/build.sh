#!/bin/bash
if [ `helm list --all nrfagent | grep -c nrfagent` -eq 1 ]; then
helm del nrfagent
helm del --purge  nrfagent
fi

docker rmi -f armdocker.rnd.ericsson.se/proj-ipworks/nrfagent:latest
go build -o nrfagent
cp nrfagent /pdu/gowork/src/gerrit.ericsson.se/udm/nrfclient/
cp /pdu/gowork/src/gerrit.ericsson.se/udm/nrfclient/deploy/helm/nrfagent/config/tls/cert.pem /pdu/gowork/src/gerrit.ericsson.se/udm/nrfclient/
cp /pdu/gowork/src/gerrit.ericsson.se/udm/nrfclient/deploy/helm/nrfagent/config/tls/key.pem /pdu/gowork/src/gerrit.ericsson.se/udm/nrfclient/
docker build -t armdocker.rnd.ericsson.se/proj-ipworks/nrfclient -f /pdu/gowork/src/gerrit.ericsson.se/udm/nrfclient/Dockerfile /pdu/gowork/src/gerrit.ericsson.se/udm/nrfclient/
docker push armdocker.rnd.ericsson.se/proj-ipworks/nrfagent:latest
docker rm `docker ps -aq -f status=exited`
helm install -n nrfagent /pdu/gowork/src/gerrit.ericsson.se/udm/nrfclient/deploy/helm/nrfagent
