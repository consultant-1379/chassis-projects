VERSION?=1.0.0
REGISTRY?=armdocker.rnd.ericsson.se
PROJECT?=proj-ipworks/nrfagent-sidecar

relbuild:
	docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true
	docker build --no-cache=true -t ${REGISTRY}/${PROJECT}:${VERSION} -f ./Dockerfile.runtime .

relpush:relbuild
	docker push ${REGISTRY}/${PROJECT}:${VERSION}

clean:
	docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true