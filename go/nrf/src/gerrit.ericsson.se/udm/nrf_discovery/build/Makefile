VERSION?=lsv5
DEVVERSION?=latest
REGISTRY?=armdocker.rnd.ericsson.se
PROJECT?=proj-ipworks/nrfdisc

release: clean relpush
devrelease: devpush


devbuild:
	sudo docker rmi -f nrf-builder:latest || true
	sudo docker build --no-cache=true -t nrf-builder:latest -f ./Dockerfile.build ../..
	sudo docker rmi -f ${REGISTRY}/${PROJECT}:${DEVVERSION} || true
	sudo docker build --no-cache=true -t ${REGISTRY}/${PROJECT}:${DEVVERSION} -f ./Dockerfile.runtime .
	sudo docker rmi -f nrf-builder:latest || true

devpush:devbuild
	sudo docker push ${REGISTRY}/${PROJECT}:${DEVVERSION}

relbuild:
	sudo docker rmi -f nrf-builder:latest || true
	sudo docker build --no-cache=true -t nrf-builder:latest -f ./Dockerfile.build ../..
	sudo docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true
	sudo docker build --no-cache=true -t ${REGISTRY}/${PROJECT}:${VERSION} -f ./Dockerfile.runtime .
	sudo docker rmi -f nrf-builder:latest || true

relpush:relbuild
	docker push ${REGISTRY}/${PROJECT}:${VERSION}

clean:
	sudo docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true

