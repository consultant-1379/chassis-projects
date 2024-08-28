DEVVERSION?=latest
REGISTRY?=armdocker.rnd.ericsson.se
PROJECT?=proj-ipworks/nrfagentpmjobloader

devbuild:
	docker rmi -f pmjobloader-builder:latest || true
	docker build --no-cache=true -t pmjobloader-builder:latest -f ./Dockerfile.build ../../..
	docker rmi -f ${REGISTRY}/${PROJECT}:${DEVVERSION} || true
	docker build --no-cache=true -t ${REGISTRY}/${PROJECT}:${DEVVERSION} -f ./Dockerfile.runtime .
	docker rmi -f pmjobloader-builder:latest || true

devpush:devbuild
	docker push ${REGISTRY}/${PROJECT}:${DEVVERSION}

clean:
	docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true


