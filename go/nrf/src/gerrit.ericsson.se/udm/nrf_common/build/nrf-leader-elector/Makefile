VERSION?=0.1
DEVVERSION?=0.1
REGISTRY?=armdocker.rnd.ericsson.se
PROJECT?=proj-ipworks/nrf-leader-elector

release: clean relpush
devrelease: devpush


devbuild:
	docker rmi -f ${REGISTRY}/${PROJECT}:${DEVVERSION} || true
	docker build --no-cache=true -t ${REGISTRY}/${PROJECT}:${DEVVERSION} -f ./Dockerfile.runtime .

devpush:devbuild
	docker push ${REGISTRY}/${PROJECT}:${DEVVERSION}

relbuild:
	docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true
	docker build --no-cache=true -t ${REGISTRY}/${PROJECT}:${VERSION} -f ./Dockerfile.runtime .
	
relpush:relbuild
	docker push ${REGISTRY}/${PROJECT}:${VERSION}

clean:
	docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true

