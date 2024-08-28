VERSION?=lsv5-0.5.6
DEVVERSION?=dev
REGISTRY?=armdocker.rnd.ericsson.se
PROJECT?=proj-ipworks/nrfagentdisc

BUILDPATH?=.

all:
	go build -ldflags "-X gerrit.ericsson.se/udm/nrfagent_discovery/cmd/nrfagent/main.Version=`date -u +%Y%m%d.%H%M%S`" -o nrfagent ../cmd/nrfagent/nrfagent.go

build: all
	docker rmi -f ${REGISTRY}/${PROJECT}:${DEVVERSION} || true
	docker build --no-cache=true -t ${REGISTRY}/${PROJECT}:${DEVVERSION} -f ./Dockerfile ../..

devbuild:
	docker rmi -f nrfagent-builder:latest || true
	docker build --no-cache=true -t nrfagent-builder:latest -f ${BUILDPATH}/Dockerfile.build ../..
	docker rmi -f ${REGISTRY}/${PROJECT}:${DEVVERSION} || true
	docker build --no-cache=true -t ${REGISTRY}/${PROJECT}:${DEVVERSION} -f ${BUILDPATH}/Dockerfile.runtime .
	docker rmi -f nrfagent-builder:latest || true

devpush: devbuild
	docker push ${REGISTRY}/${PROJECT}:${DEVVERSION}

relbuild:
	docker rmi -f nrfagent-builder:latest || true
	docker build --no-cache=true -t nrfagent-builder:latest -f ${BUILDPATH}/Dockerfile.build ../..
	docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true
	docker build --no-cache=true -t ${REGISTRY}/${PROJECT}:${VERSION} -f ${BUILDPATH}/Dockerfile.runtime .
	docker rmi -f nrfagent-builder:latest || true

relpush: relbuild
	docker push ${REGISTRY}/${PROJECT}:${VERSION}

relclean:
	docker rmi -f nrfagent-builder:latest || true
	docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true

devclean:
	docker rmi -f nrfagent-builder:latest || true
	docker rmi -f ${REGISTRY}/${PROJECT}:${DEVVERSION} || true

clean:
	rm -f nrfagent || true
	docker rmi -f ${REGISTRY}/${PROJECT}:${DEVVERSION} || true

